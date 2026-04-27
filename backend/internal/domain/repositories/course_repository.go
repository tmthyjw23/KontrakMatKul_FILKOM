package repositories

import (
	"context"

	"sistemkontrakmatkul/backend/internal/domain/models"
)

type CourseRepository interface {
	ListCourses(ctx context.Context) ([]models.CourseResponse, error)
	Create(ctx context.Context, course *models.Course) error
	Update(ctx context.Context, course *models.Course) error
	Delete(ctx context.Context, id uint64) error
}
