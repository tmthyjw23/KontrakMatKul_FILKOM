package repositories

import (
	"context"
	"sistemkontrakmatkul/backend/internal/domain/models"
)

type PassedCourseRepository interface {
	Create(ctx context.Context, passedCourse *models.PassedCourse) error
	GetByUserID(ctx context.Context, userID uint64) ([]models.PassedCourse, error)
	GetPassedCourseDetailsByUserID(ctx context.Context, userID uint64) ([]models.PassedCourseResponse, error)
}