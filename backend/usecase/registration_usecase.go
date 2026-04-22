package usecase

import (
	"context"
	"fmt"

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

// RegisterCourse validates and creates a new registration entry.
func (u *registrationUsecase) RegisterCourse(ctx context.Context, nim, courseCode string) (*domain.Registration, error) {
	if nim == "" {
		return nil, fmt.Errorf("student NIM cannot be empty")
	}
	if courseCode == "" {
		return nil, fmt.Errorf("course code cannot be empty")
	}

	// Business Rule: Ensure the course exists before registering
	if _, err := u.courseRepo.GetByCode(ctx, courseCode); err != nil {
		return nil, fmt.Errorf("RegisterCourse: course not found: %w", err)
	}

	// Business Rule: Prevent duplicate registrations
	existing, err := u.regRepo.GetByNIM(ctx, nim)
	if err != nil {
		return nil, fmt.Errorf("RegisterCourse: error checking existing registrations: %w", err)
	}
	for _, r := range existing {
		if r.CourseCode == courseCode && r.Status == "registered" {
			return nil, fmt.Errorf("RegisterCourse: student %s is already registered for course %s", nim, courseCode)
		}
	}

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
