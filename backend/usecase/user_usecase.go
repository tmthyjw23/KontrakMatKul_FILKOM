package usecase

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"kontrak-matkul/domain"
)

// userUsecase is the concrete implementation of domain.UserUsecase.
type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates and returns a new userUsecase.
func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: ur}
}

// GetProfile fetches a student's profile by NIM.
func (u *userUsecase) GetProfile(ctx context.Context, nim string) (*domain.Student, error) {
	if nim == "" {
		return nil, fmt.Errorf("NIM cannot be empty")
	}

	student, err := u.userRepo.GetByNIM(ctx, nim)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: %w", err)
	}

	return student, nil
}

// GetAllStudents fetches all registered students.
func (u *userUsecase) GetAllStudents(ctx context.Context) ([]domain.Student, error) {
	students, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllStudents: %w", err)
	}

	return students, nil
}

// CreateStudent hashes the password and saves the student.
func (u *userUsecase) CreateStudent(ctx context.Context, student *domain.Student, rawPassword string) error {
	if student.NIM == "" || student.Name == "" {
		return fmt.Errorf("NIM and Name cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	student.Role = "Student" // force role

	if err := u.userRepo.CreateStudent(ctx, student, string(hash)); err != nil {
		return fmt.Errorf("CreateStudent: %w", err)
	}

	return nil
}

// ResetPassword hashes the new password and updates it.
func (u *userUsecase) ResetPassword(ctx context.Context, nim, newRawPassword string) error {
	if nim == "" {
		return fmt.Errorf("NIM cannot be empty")
	}
	if newRawPassword == "" {
		return fmt.Errorf("new password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newRawPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	if err := u.userRepo.UpdatePassword(ctx, nim, string(hash)); err != nil {
		return fmt.Errorf("ResetPassword: %w", err)
	}

	return nil
}

// DeleteStudent removes a student by NIM.
func (u *userUsecase) DeleteStudent(ctx context.Context, nim string) error {
	if nim == "" {
		return fmt.Errorf("NIM cannot be empty")
	}

	if err := u.userRepo.Delete(ctx, nim); err != nil {
		return fmt.Errorf("DeleteStudent: %w", err)
	}

	return nil
}