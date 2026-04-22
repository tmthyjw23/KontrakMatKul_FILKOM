package domain

import "context"

// Course represents a subject available for registration.
type Course struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Class        string `json:"class"`
	LecturerName string `json:"lecturer_name"`
	Credits      int    `json:"credits"`
	CohortTarget int    `json:"cohort_target"`
}

// Schedule represents the time and place a course is held.
type Schedule struct {
	CourseCode      string `json:"course_code"`
	Day             string `json:"day"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	Room            string `json:"room"`
	AdditionalNotes string `json:"additional_notes"`
}

// CourseRepository outlines data storage operations for Course.
type CourseRepository interface {
	FetchAll(ctx context.Context) ([]Course, error)
	GetByCode(ctx context.Context, code string) (*Course, error)
}

// CourseUsecase outlines business logic operations for Course.
type CourseUsecase interface {
	FetchAllCourses(ctx context.Context) ([]Course, error)
	GetCourseDetails(ctx context.Context, code string) (*Course, error)
}
