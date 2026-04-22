package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kontrak-matkul/domain"
)

// registrationRepository is the concrete implementation of domain.RegistrationRepository.
type registrationRepository struct {
	db *sql.DB
}

// NewRegistrationRepository creates and returns a new registrationRepository.
func NewRegistrationRepository(db *sql.DB) domain.RegistrationRepository {
	return &registrationRepository{db: db}
}

// Create inserts a new registration record into the database.
func (r *registrationRepository) Create(ctx context.Context, reg *domain.Registration) error {
	query := `
		INSERT INTO registrations (student_nim, course_code, status)
		VALUES ($1, $2, 'registered')
		RETURNING id, created_at
	`

	row := r.db.QueryRowContext(ctx, query, reg.StudentNIM, reg.CourseCode)
	if err := row.Scan(&reg.ID, &reg.CreatedAt); err != nil {
		return fmt.Errorf("error creating registration: %w", err)
	}

	reg.Status = "registered"
	return nil
}

// GetByNIM retrieves all registrations for a given student NIM.
func (r *registrationRepository) GetByNIM(ctx context.Context, nim string) ([]domain.Registration, error) {
	query := `
		SELECT id, student_nim, course_code, status, created_at
		FROM registrations
		WHERE student_nim = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, nim)
	if err != nil {
		return nil, fmt.Errorf("error fetching registrations for NIM %s: %w", nim, err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		if err := rows.Scan(
			&reg.ID,
			&reg.StudentNIM,
			&reg.CourseCode,
			&reg.Status,
			&reg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning registration row: %w", err)
		}
		registrations = append(registrations, reg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return registrations, nil
}

// GetAll retrieves every registration record (admin use).
func (r *registrationRepository) GetAll(ctx context.Context) ([]domain.Registration, error) {
	query := `
		SELECT id, student_nim, course_code, status, created_at
		FROM registrations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all registrations: %w", err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		if err := rows.Scan(
			&reg.ID,
			&reg.StudentNIM,
			&reg.CourseCode,
			&reg.Status,
			&reg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning registration row: %w", err)
		}
		registrations = append(registrations, reg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return registrations, nil
}

// Cancel updates a registration's status to "cancelled".
func (r *registrationRepository) Cancel(ctx context.Context, id int) error {
	query := `UPDATE registrations SET status = 'cancelled' WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error cancelling registration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("registration with id %d not found", id)
	}

	return nil
}
