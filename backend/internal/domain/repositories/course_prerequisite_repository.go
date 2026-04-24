package repositories

import (
	"context"
	"sistemkontrakmatkul/backend/internal/domain/models"
)

type CoursePrerequisiteRepository interface {
	Create(ctx context.Context, prereq *models.CoursePrerequisite) error
	Delete(ctx context.Context, courseID uint64, prereqID uint64) error
	GetPrerequisitesByCourseID(ctx context.Context, courseID uint64) ([]models.CoursePrerequisite, error)
	GetPrerequisitesForCourse(ctx context.Context, courseID uint64) ([]uint64, error)
	CheckPrerequisitesMet(ctx context.Context, userID uint64, courseID uint64) (bool, error)
}