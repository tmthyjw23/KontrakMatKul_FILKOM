package usecase

import (
	"context"
	"fmt"

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