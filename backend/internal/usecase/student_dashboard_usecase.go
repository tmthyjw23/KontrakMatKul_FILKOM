package usecase

import (
	"context"
	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
)

type StudentDashboardUsecase interface {
	GetMySchedule(ctx context.Context, userID uint64) ([]models.StudentScheduleResponse, error)
	GetMyHistory(ctx context.Context, userID uint64) ([]models.PassedCourseResponse, error)
}

type studentDashboardUsecase struct {
	enrollRepo  repositories.EnrollmentRepository
	courseRepo  repositories.CourseRepository
	passedRepo  repositories.PassedCourseRepository
}

func NewStudentDashboardUsecase(er repositories.EnrollmentRepository, cr repositories.CourseRepository, pr repositories.PassedCourseRepository) StudentDashboardUsecase {
	return &studentDashboardUsecase{
		enrollRepo: er,
		courseRepo: cr,
		passedRepo: pr,
	}
}

func (u *studentDashboardUsecase) GetMySchedule(ctx context.Context, userID uint64) ([]models.StudentScheduleResponse, error) {
	return u.enrollRepo.GetStudentSchedule(ctx, userID)
}

func (u *studentDashboardUsecase) GetMyHistory(ctx context.Context, userID uint64) ([]models.PassedCourseResponse, error) {
	return u.passedRepo.GetPassedCourseDetailsByUserID(ctx, userID)
}