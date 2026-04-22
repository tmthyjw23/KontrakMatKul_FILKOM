package repositories

import (
	"context"

	"sistemkontrakmatkul/backend/internal/domain/models"
)

type CourseRepository interface {
	ListCourses(ctx context.Context) ([]models.CourseResponse, error)
}
