package usecase

import (
	"context"
	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
)

type CoursePrerequisiteUsecase interface {
	AddPrerequisite(ctx context.Context, req models.CoursePrerequisiteRequest) error
	RemovePrerequisite(ctx context.Context, courseID uint64, prereqID uint64) error
	GetPrerequisites(ctx context.Context, courseID uint64) ([]models.CoursePrerequisite, error)
}

type coursePrerequisiteUsecase struct {
	repo repositories.CoursePrerequisiteRepository
}

func NewCoursePrerequisiteUsecase(repo repositories.CoursePrerequisiteRepository) CoursePrerequisiteUsecase {
	return &coursePrerequisiteUsecase{repo: repo}
}

func (u *coursePrerequisiteUsecase) AddPrerequisite(ctx context.Context, req models.CoursePrerequisiteRequest) error {
	return u.repo.Create(ctx, &models.CoursePrerequisite{
		CourseID:             req.CourseID,
		PrerequisiteCourseID: req.PrerequisiteCourseID,
	})
}

func (u *coursePrerequisiteUsecase) RemovePrerequisite(ctx context.Context, courseID uint64, prereqID uint64) error {
	return u.repo.Delete(ctx, courseID, prereqID)
}

func (u *coursePrerequisiteUsecase) GetPrerequisites(ctx context.Context, courseID uint64) ([]models.CoursePrerequisite, error) {
	return u.repo.GetPrerequisitesByCourseID(ctx, courseID)
}