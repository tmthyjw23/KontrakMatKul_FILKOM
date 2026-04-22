package domain

import "context"

// User represents the student or admin user in the system.
type User struct {
	NIM          string `json:"nim"`
	Name         string `json:"name"`
	Faculty      string `json:"faculty"`
	StudyProgram string `json:"study_program"`
	CohortYear   int    `json:"cohort_year"`
	Role         string `json:"role"`
}

// UserRepository outlines data storage operations for User.
type UserRepository interface {
	GetByNIM(ctx context.Context, nim string) (*User, error)
}

// UserUsecase outlines business logic operations for User.
type UserUsecase interface {
	GetProfile(ctx context.Context, nim string) (*User, error)
}
