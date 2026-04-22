package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kontrak-matkul/domain"
)

// courseRepository is the concrete implementation of domain.CourseRepository.
type courseRepository struct {
	db *sql.DB
}

// NewCourseRepository creates and returns a new courseRepository.
func NewCourseRepository(db *sql.DB) domain.CourseRepository {
	return &courseRepository{db: db}
}

// FetchAll retrieves all courses from the database.
func (r *courseRepository) FetchAll(ctx context.Context) ([]domain.Course, error) {
	query := `
		SELECT code, name, class, lecturer_name, credits, cohort_target
		FROM courses
		ORDER BY code ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all courses: %w", err)
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(
			&c.Code,
			&c.Name,
			&c.Class,
			&c.LecturerName,
			&c.Credits,
			&c.CohortTarget,
		); err != nil {
			return nil, fmt.Errorf("error scanning course row: %w", err)
		}
		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return courses, nil
}

// GetByCode retrieves a single course by its code.
func (r *courseRepository) GetByCode(ctx context.Context, code string) (*domain.Course, error) {
	query := `
		SELECT code, name, class, lecturer_name, credits, cohort_target
		FROM courses
		WHERE code = $1
	`

	row := r.db.QueryRowContext(ctx, query, code)

	var c domain.Course
	err := row.Scan(
		&c.Code,
		&c.Name,
		&c.Class,
		&c.LecturerName,
		&c.Credits,
		&c.CohortTarget,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("course with code %s not found", code)
		}
		return nil, fmt.Errorf("error fetching course: %w", err)
	}

	return &c, nil
}

// Create inserts a new course record into the database.
func (r *courseRepository) Create(ctx context.Context, course *domain.Course) error {
	query := `
		INSERT INTO courses (code, name, class, lecturer_name, credits, cohort_target)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		course.Code,
		course.Name,
		course.Class,
		course.LecturerName,
		course.Credits,
		course.CohortTarget,
	)
	if err != nil {
		return fmt.Errorf("error creating course: %w", err)
	}

	return nil
}

// Delete removes a course record from the database by its code.
func (r *courseRepository) Delete(ctx context.Context, code string) error {
	query := `DELETE FROM courses WHERE code = $1`

	result, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		return fmt.Errorf("error deleting course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("course with code %s not found", code)
	}

	return nil
}