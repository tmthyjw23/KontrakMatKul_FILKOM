package usecase

import (
	"context"
	"fmt"

	"kontrak-matkul/domain"
)

// courseUsecase is the concrete implementation of domain.CourseUsecase.
type courseUsecase struct {
	courseRepo domain.CourseRepository
}

// NewCourseUsecase creates and returns a new courseUsecase.
func NewCourseUsecase(cr domain.CourseRepository) domain.CourseUsecase {
	return &courseUsecase{courseRepo: cr}
}

// FetchAllCourses retrieves all available courses.
func (u *courseUsecase) FetchAllCourses(ctx context.Context) ([]domain.Course, error) {
	courses, err := u.courseRepo.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("FetchAllCourses: %w", err)
	}

	return courses, nil
}

// GetCourseDetails retrieves a single course by code.
func (u *courseUsecase) GetCourseDetails(ctx context.Context, code string) (*domain.Course, error) {
	if code == "" {
		return nil, fmt.Errorf("course code cannot be empty")
	}

	course, err := u.courseRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("GetCourseDetails: %w", err)
	}

	return course, nil
}

// CreateCourse validates and persists a new course.
func (u *courseUsecase) CreateCourse(ctx context.Context, course *domain.Course) error {
	if course.Code == "" {
		return fmt.Errorf("course code is required")
	}
	if course.Name == "" {
		return fmt.Errorf("course name is required")
	}
	if course.Credits <= 0 {
		return fmt.Errorf("credits must be a positive number")
	}

	if err := u.courseRepo.Create(ctx, course); err != nil {
		return fmt.Errorf("CreateCourse: %w", err)
	}

	return nil
}

// DeleteCourse removes a course by code after validating it exists.
func (u *courseUsecase) DeleteCourse(ctx context.Context, code string) error {
	if code == "" {
		return fmt.Errorf("course code cannot be empty")
	}

	// Confirm course exists before attempting deletion
	if _, err := u.courseRepo.GetByCode(ctx, code); err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}

	if err := u.courseRepo.Delete(ctx, code); err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}

	return nil
}