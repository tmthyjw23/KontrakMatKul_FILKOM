package models

import (
	"time"
)

type PassedCourse struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	CourseID  uint64    `json:"course_id"`
	Grade     string    `json:"grade"`
	PassedAt  *time.Time `json:"passed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PassedCourseResponse struct {
	CourseCode string `json:"course_code"`
	CourseName string `json:"course_name"`
	Grade      string `json:"grade"`
}

type StudentScheduleResponse struct {
	CourseCode string    `json:"course_code"`
	CourseName string    `json:"course_name"`
	Day        string    `json:"day"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Room       string    `json:"room"`
	Lecturer   string    `json:"lecturer"`
}