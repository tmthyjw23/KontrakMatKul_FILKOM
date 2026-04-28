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
}

// NewRegistrationUsecase creates and returns a new registrationUsecase.
// It takes both a registration and course repository to enforce business rules
// (e.g., a student can only register for an existing course).
func NewRegistrationUsecase(rr domain.RegistrationRepository, cr domain.CourseRepository) domain.RegistrationUsecase {
	return &registrationUsecase{
		regRepo:    rr,
		courseRepo: cr,
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
		if r.Status == "registered" {
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
