package models

import (
	"errors"
	"time"
)

var (
	ErrInvalidEnrollmentRequest = errors.New("invalid enrollment request")
	ErrUserNotFound             = errors.New("user not found")
	ErrCourseNotFound           = errors.New("course not found")
	ErrAlreadyEnrolled          = errors.New("already enrolled")
	ErrCreditLimitExceeded      = errors.New("credit limit exceeded")
	ErrQuotaExceeded            = errors.New("quota exceeded")
	ErrScheduleConflict         = errors.New("Schedule Conflict")
)

type EnrollmentRequest struct {
	UserID   uint64 `json:"-"`
	CourseID uint64 `json:"course_id" binding:"required"`
}

type Enrollment struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	CourseID   uint64    `json:"course_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type EnrollmentResult struct {
	Enrollment Enrollment `json:"enrollment"`
}

type UserCreditInfo struct {
	ID     uint64
	MaxSKS int
}

type Course struct {
	ID    uint64 `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	SKS   int    `json:"sks"`
	Quota int    `json:"quota"`
}

type ScheduleSlot struct {
	CourseID   uint64 `json:"course_id"`
	CourseCode string `json:"course_code"`
	Day        string `json:"day"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
}
