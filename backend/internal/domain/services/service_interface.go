package services

import (
	"context"

	"sistemkontrakmatkul/backend/internal/domain/models"
)

type EnrollmentService interface {
	Enroll(ctx context.Context, request models.EnrollmentRequest) (*models.EnrollmentResult, error)
	ListAll(ctx context.Context) ([]models.Enrollment, error)
	Approve(ctx context.Context, enrollmentID uint64) error
	Reject(ctx context.Context, enrollmentID uint64) error
}

type CourseService interface {
	ListCourses(ctx context.Context) ([]models.CourseResponse, error)
	Create(ctx context.Context, course *models.Course) error
	Update(ctx context.Context, course *models.Course) error
	Delete(ctx context.Context, id uint64) error
}
