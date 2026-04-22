package services

import (
	"context"

	"sistemkontrakmatkul/backend/internal/domain/models"
)

type EnrollmentService interface {
	Enroll(ctx context.Context, request models.EnrollmentRequest) (*models.EnrollmentResult, error)
}

type CourseService interface {
	ListCourses(ctx context.Context) ([]models.CourseResponse, error)
}
