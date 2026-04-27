package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sistemkontrakmatkul/backend/internal/domain/models"
)

type passedCourseRepository struct {
	db *sql.DB
}

func NewPassedCourseRepository(db *sql.DB) *passedCourseRepository {
	return &passedCourseRepository{db: db}
}

func (r *passedCourseRepository) Create(ctx context.Context, pc *models.PassedCourse) error {
	query := `INSERT INTO passed_courses (user_id, course_id, grade, passed_at, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, pc.UserID, pc.CourseID, pc.Grade, pc.PassedAt, pc.CreatedAt, pc.UpdatedAt)
	return err
}

func (r *passedCourseRepository) GetByUserID(ctx context.Context, userID uint64) ([]models.PassedCourse, error) {
	query := `SELECT id, user_id, course_id, grade, passed_at, created_at, updated_at FROM passed_courses WHERE user_id = ?`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.PassedCourse
	for rows.Next() {
		var pc models.PassedCourse
		if err := rows.Scan(&pc.ID, &pc.UserID, &pc.CourseID, &pc.Grade, &pc.PassedAt, &pc.CreatedAt, &pc.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, pc)
	}
	return results, nil
}

func (r *passedCourseRepository) GetPassedCourseDetailsByUserID(ctx context.Context, userID uint64) ([]models.PassedCourseResponse, error) {
	query := `SELECT c.code, c.name, pc.grade 
              FROM passed_courses pc 
              JOIN courses c ON pc.course_id = c.id 
              WHERE pc.user_id = ?`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]models.PassedCourseResponse, 0)
	for rows.Next() {
		var res models.PassedCourseResponse
		if err := rows.Scan(&res.CourseCode, &res.CourseName, &res.Grade); err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

func (r *passedCourseRepository) HasPassedCourses(ctx context.Context, userID uint64, courseIDs []uint64) (bool, error) {
	if len(courseIDs) == 0 {
		return true, nil
	}

	// Create placeholders for the IN clause: ?, ?, ?
	placeholders := make([]string, len(courseIDs))
	args := make([]interface{}, 0, len(courseIDs)+1)
	args = append(args, userID)

	for i, id := range courseIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM passed_courses WHERE user_id = ? AND course_id IN (%s)`, 
		strings.Join(placeholders, ","))

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == len(courseIDs), nil
}