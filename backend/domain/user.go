package domain

import "context"

// Student represents a registered student user in the system.
type Student struct {
	NIM          string `json:"nim"           db:"nim"`
	Name         string `json:"name"          db:"name"`
	Faculty      string `json:"faculty"       db:"faculty"`
	StudyProgram string `json:"study_program" db:"study_program"`
	CohortYear   int    `json:"cohort_year"   db:"cohort_year"`
	Role         string `json:"role"          db:"role"`
}

// UserRepository defines the contract for User data storage operations.
// Concrete implementations will reside in the repository layer.
type UserRepository interface {
	GetByNIM(ctx context.Context, nim string) (*Student, error)
	GetAll(ctx context.Context) ([]Student, error)
	// GetPasswordHashByNIM returns only the bcrypt password hash for a given NIM.
	// The hash is intentionally NOT a field on the Student struct to prevent
	// it from being accidentally serialized into JSON API responses.
	GetPasswordHashByNIM(ctx context.Context, nim string) (string, error)
}

// UserUsecase defines the contract for User business logic operations.
// Concrete implementations will reside in the usecase layer.
type UserUsecase interface {
	GetProfile(ctx context.Context, nim string) (*Student, error)
	GetAllStudents(ctx context.Context) ([]Student, error)
}
