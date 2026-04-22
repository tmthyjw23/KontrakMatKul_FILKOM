package domain

import "context"

// Registration represents a student's course registration record.
type Registration struct {
	ID         int    `json:"id"          db:"id"`
	StudentNIM string `json:"student_nim" db:"student_nim"`
	CourseCode string `json:"course_code" db:"course_code"`
	Status     string `json:"status"      db:"status"` // e.g., "registered", "cancelled"
	CreatedAt  string `json:"created_at"  db:"created_at"`
}

// RegistrationRepository defines the contract for Registration data operations.
type RegistrationRepository interface {
	Create(ctx context.Context, reg *Registration) error
	GetByNIM(ctx context.Context, nim string) ([]Registration, error)
	GetAll(ctx context.Context) ([]Registration, error)
	Cancel(ctx context.Context, id int) error
}

// RegistrationUsecase defines the contract for Registration business logic.
type RegistrationUsecase interface {
	RegisterCourse(ctx context.Context, nim, courseCode string) (*Registration, error)
	GetRegistrationsByNIM(ctx context.Context, nim string) ([]Registration, error)
	GetAllRegistrations(ctx context.Context) ([]Registration, error)
	CancelRegistration(ctx context.Context, id int) error
}
