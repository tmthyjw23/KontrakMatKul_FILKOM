package usecase

import (
	"context"
	"fmt"
	"time"

	"kontrak-matkul/domain"
)

// registrationUsecase is the concrete implementation of domain.RegistrationUsecase.
type registrationUsecase struct {
	regRepo    domain.RegistrationRepository
	courseRepo domain.CourseRepository
	cpRepo     domain.ContractPeriodRepository
}

// NewRegistrationUsecase creates and returns a new registrationUsecase.
// It takes both a registration and course repository to enforce business rules
// (e.g., a student can only register for an existing course).
func NewRegistrationUsecase(rr domain.RegistrationRepository, cr domain.CourseRepository, cpr domain.ContractPeriodRepository) domain.RegistrationUsecase {
	return &registrationUsecase{
		regRepo:    rr,
		courseRepo: cr,
		cpRepo:     cpr,
	}
}

// isTimeOverlap checks if two time ranges (HH:MM) overlap.
func isTimeOverlap(startA, endA, startB, endB string) (bool, error) {
	layout := "15:04" // Standard format for HH:MM in Go
	tSA, err := time.Parse(layout, startA)
	if err != nil {
		return false, fmt.Errorf("invalid time format startA: %w", err)
	}
	tEA, err := time.Parse(layout, endA)
	if err != nil {
		return false, fmt.Errorf("invalid time format endA: %w", err)
	}
	tSB, err := time.Parse(layout, startB)
	if err != nil {
		return false, fmt.Errorf("invalid time format startB: %w", err)
	}
	tEB, err := time.Parse(layout, endB)
	if err != nil {
		return false, fmt.Errorf("invalid time format endB: %w", err)
	}

	// Overlap condition: StartA < EndB && EndA > StartB
	return tSA.Before(tEB) && tEA.After(tSB), nil
}

// RegisterCourse validates and creates a new registration entry.
func (u *registrationUsecase) RegisterCourse(ctx context.Context, nim, courseCode string) (*domain.Registration, error) {
	if nim == "" {
		return nil, fmt.Errorf("student NIM cannot be empty")
	}
	if courseCode == "" {
		return nil, fmt.Errorf("course code cannot be empty")
	}

	// 0. Ensure Contract Period is OPEN
	period, err := u.cpRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("RegisterCourse: error checking contract period: %w", err)
	}
	if period != nil && !period.IsOpen {
		return nil, domain.ErrContractPeriodClosed
	}

	// 1. Ensure the new course exists
	newCourse, err := u.courseRepo.GetByCode(ctx, courseCode)
	if err != nil {
		return nil, fmt.Errorf("RegisterCourse: course not found: %w", err)
	}

	// 2. Fetch existing active registrations
	existingRegs, err := u.regRepo.GetByNIM(ctx, nim)
	if err != nil {
		return nil, fmt.Errorf("RegisterCourse: error checking existing registrations: %w", err)
	}

	// NEW RULE: If student has ANY 'pending' registrations, they cannot register for more
	for _, r := range existingRegs {
		if r.Status == "pending" {
			return nil, domain.ErrPendingRegistration
		}
	}

	totalCredits := newCourse.Credits
	var activeCourseCodes []string

	for _, r := range existingRegs {
		if r.Status == "approved" || r.Status == "registered" {
			if r.CourseCode == courseCode {
				return nil, fmt.Errorf("RegisterCourse: student %s is already registered for course %s", nim, courseCode)
			}

			// Accumulate credits for existing active courses
			c, err := u.courseRepo.GetByCode(ctx, r.CourseCode)
			if err != nil {
				return nil, fmt.Errorf("RegisterCourse: error fetching existing course details for calculating SKS: %w", err)
			}
			totalCredits += c.Credits
			activeCourseCodes = append(activeCourseCodes, r.CourseCode)
		}
	}

	// 3. Business Rule: Max SKS Limit
	if totalCredits > 24 {
		return nil, domain.ErrMaxCreditsExceeded
	}

	// 4. Business Rule: Schedule Clash Detection
	newSchedules, err := u.courseRepo.GetSchedulesByCourseCode(ctx, courseCode)
	if err != nil {
		return nil, fmt.Errorf("RegisterCourse: error fetching new course schedules: %w", err)
	}

	for _, activeCode := range activeCourseCodes {
		existSchedules, err := u.courseRepo.GetSchedulesByCourseCode(ctx, activeCode)
		if err != nil {
			return nil, fmt.Errorf("RegisterCourse: error fetching active course schedules: %w", err)
		}

		for _, nSched := range newSchedules {
			for _, eSched := range existSchedules {
				// We only care if they are on the same day
				if nSched.DayOfWeek == eSched.DayOfWeek {
					overlap, err := isTimeOverlap(nSched.StartTime, nSched.EndTime, eSched.StartTime, eSched.EndTime)
					if err != nil {
						return nil, fmt.Errorf("RegisterCourse: error parsing schedule times: %w", err)
					}
					if overlap {
						return nil, fmt.Errorf("%w: conflict with course %s on %s (%s-%s)", 
							domain.ErrScheduleConflict, activeCode, nSched.DayOfWeek, eSched.StartTime, eSched.EndTime)
					}
				}
			}
		}
	}

	// Validation passed, create registration
	reg := &domain.Registration{
		StudentNIM: nim,
		CourseCode: courseCode,
	}

	if err := u.regRepo.Create(ctx, reg); err != nil {
		return nil, fmt.Errorf("RegisterCourse: %w", err)
	}

	return reg, nil
}

// BulkRegisterCourse handles multiple course registrations in one go.
func (u *registrationUsecase) BulkRegisterCourse(ctx context.Context, nim string, courseCodes []string) ([]domain.Registration, error) {
	if nim == "" {
		return nil, fmt.Errorf("student NIM cannot be empty")
	}
	if len(courseCodes) == 0 {
		return nil, fmt.Errorf("no course codes provided")
	}

	// 0. Ensure Contract Period is OPEN
	period, err := u.cpRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("BulkRegisterCourse: error checking contract period: %w", err)
	}
	if period != nil && !period.IsOpen {
		return nil, domain.ErrContractPeriodClosed
	}

	// 1. Fetch existing registrations
	existingRegs, err := u.regRepo.GetByNIM(ctx, nim)
	if err != nil {
		return nil, fmt.Errorf("BulkRegisterCourse: %w", err)
	}

	// 2. NEW RULE: Block if there are already pending registrations
	for _, r := range existingRegs {
		if r.Status == "pending" {
			return nil, domain.ErrPendingRegistration
		}
	}

	// 3. Collect all "active" courses (approved or registered)
	var activeCourses []domain.Course
	var activeCourseCodes []string
	totalSks := 0

	for _, r := range existingRegs {
		if r.Status == "approved" || r.Status == "registered" {
			c, err := u.courseRepo.GetByCode(ctx, r.CourseCode)
			if err != nil {
				return nil, fmt.Errorf("BulkRegisterCourse: error fetching active course: %w", err)
			}
			activeCourses = append(activeCourses, *c)
			activeCourseCodes = append(activeCourseCodes, r.CourseCode)
			totalSks += c.Credits
		}
	}

	// 4. Validate each new course
	var newRegistrations []domain.Registration
	var newCourses []domain.Course

	for _, code := range courseCodes {
		// Check for duplicate in request
		for _, alreadyAdded := range newRegistrations {
			if alreadyAdded.CourseCode == code {
				return nil, fmt.Errorf("duplicate course code in request: %s", code)
			}
		}
		// Check if already registered
		for _, activeCode := range activeCourseCodes {
			if activeCode == code {
				return nil, fmt.Errorf("student is already registered for course %s", code)
			}
		}

		course, err := u.courseRepo.GetByCode(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("course %s not found", code)
		}

		totalSks += course.Credits
		newCourses = append(newCourses, *course)
		newRegistrations = append(newRegistrations, domain.Registration{
			StudentNIM: nim,
			CourseCode: code,
		})
	}

	// 5. Check SKS limit
	if totalSks > 24 {
		return nil, domain.ErrMaxCreditsExceeded
	}

	// 6. Check conflicts (between new courses and between new & active)
	allCoursesToCompare := append([]domain.Course{}, activeCourses...)
	for _, newC := range newCourses {
		newSchedules, err := u.courseRepo.GetSchedulesByCourseCode(ctx, newC.Code)
		if err != nil {
			return nil, fmt.Errorf("error fetching schedules for %s: %w", newC.Code, err)
		}

		// Compare newC with all existing/active courses
		for _, otherC := range allCoursesToCompare {
			otherSchedules, err := u.courseRepo.GetSchedulesByCourseCode(ctx, otherC.Code)
			if err != nil {
				return nil, fmt.Errorf("error fetching schedules for %s: %w", otherC.Code, err)
			}

			for _, ns := range newSchedules {
				for _, os := range otherSchedules {
					if ns.DayOfWeek == os.DayOfWeek {
						overlap, err := isTimeOverlap(ns.StartTime, ns.EndTime, os.StartTime, os.EndTime)
						if err != nil {
							return nil, err
						}
						if overlap {
							return nil, fmt.Errorf("%w: conflict between %s and %s on %s", 
								domain.ErrScheduleConflict, newC.Code, otherC.Code, ns.DayOfWeek)
						}
					}
				}
			}
		}
		// Add newC to comparison pool for subsequent new courses in this batch
		allCoursesToCompare = append(allCoursesToCompare, newC)
	}

	// 7. Everything passed, save them
	var saved []domain.Registration
	for _, reg := range newRegistrations {
		r := reg
		if err := u.regRepo.Create(ctx, &r); err != nil {
			return nil, fmt.Errorf("BulkRegisterCourse: error saving %s: %w", reg.CourseCode, err)
		}
		saved = append(saved, r)
	}

	return saved, nil
}

// GetRegistrationsByNIM fetches all registrations for a given student.
func (u *registrationUsecase) GetRegistrationsByNIM(ctx context.Context, nim string) ([]domain.Registration, error) {
	if nim == "" {
		return nil, fmt.Errorf("NIM cannot be empty")
	}

	return u.regRepo.GetByNIM(ctx, nim)
}

// GetAllRegistrations fetches every registration (admin only).
func (u *registrationUsecase) GetAllRegistrations(ctx context.Context) ([]domain.Registration, error) {
	return u.regRepo.GetAll(ctx)
}

// CancelRegistration cancels a registration by ID.
func (u *registrationUsecase) CancelRegistration(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid registration ID")
	}

	if err := u.regRepo.Cancel(ctx, id); err != nil {
		return fmt.Errorf("CancelRegistration: %w", err)
	}

	return nil
}

// ApproveRegistration sets status to 'approved' (or 'registered').
func (u *registrationUsecase) ApproveRegistration(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid registration ID")
	}
	return u.regRepo.UpdateStatus(ctx, id, "approved")
}

// RejectRegistration sets status to 'rejected'.
func (u *registrationUsecase) RejectRegistration(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid registration ID")
	}
	return u.regRepo.UpdateStatus(ctx, id, "rejected")
}
