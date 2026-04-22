package repositories

import (
	"context"
	"sistemkontrakmatkul/backend/internal/domain/models"
)

type UserRepository interface {
	GetByStudentNumber(ctx context.Context, studentNumber string) (*models.User, error)
	GetByID(ctx context.Context, id uint64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint64) error
}
