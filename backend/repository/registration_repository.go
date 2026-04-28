package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
// MySQL does not support RETURNING, so we use LastInsertId() to retrieve the new ID
// and set CreatedAt manually.
func (r *registrationRepository) Create(ctx context.Context, reg *domain.Registration) error {
	query := `
		INSERT INTO registrations (student_nim, course_code, status)
		VALUES (?, ?, 'pending')
	`

	result, err := r.db.ExecContext(ctx, query, reg.StudentNIM, reg.CourseCode)
	if err != nil {
		return fmt.Errorf("error creating registration: %w", err)
	}

	// MySQL equivalent of PostgreSQL's RETURNING id
	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error retrieving last insert ID: %w", err)
	}

	reg.ID = int(lastID)
	reg.Status = "pending"
	reg.CreatedAt = time.Now().Format(time.DateTime)
	return nil
}

// GetByNIM retrieves all registrations for a given student NIM.
func (r *registrationRepository) GetByNIM(ctx context.Context, nim string) ([]domain.Registration, error) {
	query := `
		SELECT r.id, r.student_nim, s.name, r.course_code, c.name, r.status, r.created_at
		FROM registrations r
		LEFT JOIN students s ON r.student_nim = s.nim
		LEFT JOIN courses c ON r.course_code = c.code
		WHERE r.student_nim = ?
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, nim)
	if err != nil {
		return nil, fmt.Errorf("error fetching registrations for NIM %s: %w", nim, err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		var studentName, courseName sql.NullString
		if err := rows.Scan(
			&reg.ID,
			&reg.StudentNIM,
			&studentName,
			&reg.CourseCode,
			&courseName,
			&reg.Status,
			&reg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning registration row: %w", err)
		}
		if studentName.Valid {
			reg.StudentName = studentName.String
		} else {
			reg.StudentName = reg.StudentNIM
		}
		if courseName.Valid {
			reg.CourseName = courseName.String
		} else {
			reg.CourseName = reg.CourseCode
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
		SELECT r.id, r.student_nim, s.name, r.course_code, c.name, r.status, r.created_at
		FROM registrations r
		LEFT JOIN students s ON r.student_nim = s.nim
		LEFT JOIN courses c ON r.course_code = c.code
		ORDER BY r.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all registrations: %w", err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		var studentName, courseName sql.NullString
		if err := rows.Scan(
			&reg.ID,
			&reg.StudentNIM,
			&studentName,
			&reg.CourseCode,
			&courseName,
			&reg.Status,
			&reg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning registration row: %w", err)
		}
		if studentName.Valid {
			reg.StudentName = studentName.String
		} else {
			reg.StudentName = reg.StudentNIM
		}
		if courseName.Valid {
			reg.CourseName = courseName.String
		} else {
			reg.CourseName = reg.CourseCode
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
	query := `UPDATE registrations SET status = 'cancelled' WHERE id = ?`

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
