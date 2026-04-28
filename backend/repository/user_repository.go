package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kontrak-matkul/domain"
)

// userRepository is the concrete implementation of domain.UserRepository.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates and returns a new userRepository.
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// GetByNIM retrieves a single student record by their NIM.
func (r *userRepository) GetByNIM(ctx context.Context, nim string) (*domain.Student, error) {
	query := `
		SELECT nim, name, faculty, study_program, cohort_year, role
		FROM students
		WHERE nim = ?
	`

	row := r.db.QueryRowContext(ctx, query, nim)

	var s domain.Student
	err := row.Scan(
		&s.NIM,
		&s.Name,
		&s.Faculty,
		&s.StudyProgram,
		&s.CohortYear,
		&s.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student with NIM %s not found", nim)
		}
		return nil, fmt.Errorf("error fetching student: %w", err)
	}

	return &s, nil
}

// GetPasswordHashByNIM retrieves only the password hash for a given student.
func (r *userRepository) GetPasswordHashByNIM(ctx context.Context, nim string) (string, error) {
	query := `SELECT password_hash FROM students WHERE nim = ?`
	
	row := r.db.QueryRowContext(ctx, query, nim)
	var hash string
	if err := row.Scan(&hash); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("student with NIM %s not found", nim)
		}
		return "", fmt.Errorf("error fetching password hash: %w", err)
	}
	
	return hash, nil
}

// GetAll retrieves all student records from the database.
func (r *userRepository) GetAll(ctx context.Context) ([]domain.Student, error) {
	query := `
		SELECT nim, name, faculty, study_program, cohort_year, role
		FROM students
		ORDER BY cohort_year DESC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all students: %w", err)
	}
	defer rows.Close()

	var students []domain.Student
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(
			&s.NIM,
			&s.Name,
			&s.Faculty,
			&s.StudyProgram,
			&s.CohortYear,
			&s.Role,
		); err != nil {
			return nil, fmt.Errorf("error scanning student row: %w", err)
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return students, nil
}

// CreateStudent inserts a new student record into the database.
func (r *userRepository) CreateStudent(ctx context.Context, student *domain.Student, passwordHash string) error {
	query := `
		INSERT INTO students (nim, name, faculty, study_program, cohort_year, password_hash, role)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		student.NIM,
		student.Name,
		student.Faculty,
		student.StudyProgram,
		student.CohortYear,
		passwordHash,
		student.Role,
	)
	if err != nil {
		return fmt.Errorf("error creating student: %w", err)
	}

	return nil
}

// UpdatePassword updates the password hash for a given student NIM.
func (r *userRepository) UpdatePassword(ctx context.Context, nim, passwordHash string) error {
	query := `UPDATE students SET password_hash = ? WHERE nim = ?`

	result, err := r.db.ExecContext(ctx, query, passwordHash, nim)
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("student with NIM %s not found", nim)
	}

	return nil
}