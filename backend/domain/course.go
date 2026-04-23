package domain

import "context"

// Course represents a subject (mata kuliah) available for registration.
type Course struct {
	Code         string `json:"code"          db:"code"`
	Name         string `json:"name"          db:"name"`
	Class        string `json:"class"         db:"class"`
	LecturerName string `json:"lecturer_name" db:"lecturer_name"`
	Credits      int    `json:"credits"       db:"credits"`
	CohortTarget int    `json:"cohort_target" db:"cohort_target"`
}

// Schedule represents the time and room allocation for a given course.
type Schedule struct {
	CourseCode      string `json:"course_code"      db:"course_code"`
	DayOfWeek       string `json:"day_of_week"      db:"day_of_week"`
	StartTime       string `json:"start_time"       db:"start_time"`
	EndTime         string `json:"end_time"         db:"end_time"`
	Room            string `json:"room"             db:"room"`
	AdditionalNotes string `json:"additional_notes" db:"additional_notes"`
}

// CourseRepository defines the contract for Course data storage operations.
// Concrete implementations will reside in the repository layer.
type CourseRepository interface {
	FetchAll(ctx context.Context) ([]Course, error)
	GetByCode(ctx context.Context, code string) (*Course, error)
	GetSchedulesByCourseCode(ctx context.Context, code string) ([]Schedule, error)
	Create(ctx context.Context, course *Course) error
	Delete(ctx context.Context, code string) error
}

// CourseUsecase defines the contract for Course business logic operations.
// Concrete implementations will reside in the usecase layer.
type CourseUsecase interface {
	FetchAllCourses(ctx context.Context) ([]Course, error)
	GetCourseDetails(ctx context.Context, code string) (*Course, error)
	CreateCourse(ctx context.Context, course *Course) error
	DeleteCourse(ctx context.Context, code string) error
}
