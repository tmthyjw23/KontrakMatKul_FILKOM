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
