package usecase

import (
	"context"

	"go.uber.org/zap"

	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
	"sistemkontrakmatkul/backend/internal/domain/services"
)

type CourseUsecase struct {
	repo   repositories.CourseRepository
	logger *zap.Logger
}

var _ services.CourseService = (*CourseUsecase)(nil)

func NewCourseUsecase(repo repositories.CourseRepository, logger *zap.Logger) *CourseUsecase {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &CourseUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *CourseUsecase) ListCourses(ctx context.Context) ([]models.CourseResponse, error) {
	courses, err := u.repo.ListCourses(ctx)
	if err != nil {
		u.logger.Error("failed to list courses", zap.Error(err))
		return nil, models.ErrFailedToFetchCourses
	}

	return courses, nil
}

func (u *CourseUsecase) Create(ctx context.Context, course *models.Course) error {
	// Validation: Ensure basic fields are present
	if course.Code == "" || course.Name == "" {
		return models.ErrInvalidEnrollmentRequest // Reusing basic validation error
	}

	if course.SKS < 1 || course.SKS > 6 {
		return models.ErrInvalidEnrollmentRequest 
	}

	err := u.repo.Create(ctx, course)
	if err != nil {
		u.logger.Error("failed to create course in usecase", zap.Error(err))
		return err
	}

	return nil
}

func (u *CourseUsecase) Update(ctx context.Context, course *models.Course) error {
	if course.ID == 0 {
		return models.ErrInvalidEnrollmentRequest
	}

	err := u.repo.Update(ctx, course)
	if err != nil {
		u.logger.Error("failed to update course in usecase", zap.Uint64("course_id", course.ID), zap.Error(err))
		return err
	}

	return nil
}

func (u *CourseUsecase) Delete(ctx context.Context, id uint64) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		u.logger.Error("failed to delete course in usecase", zap.Uint64("course_id", id), zap.Error(err))
		return err
	}

	return nil
}
