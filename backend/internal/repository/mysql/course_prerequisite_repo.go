package mysql

import (
	"context"
	"database/sql"
	"sistemkontrakmatkul/backend/internal/domain/models"
)

type coursePrerequisiteRepository struct {
	db *sql.DB
}

func NewCoursePrerequisiteRepository(db *sql.DB) *coursePrerequisiteRepository {
	return &coursePrerequisiteRepository{db: db}
}

func (r *coursePrerequisiteRepository) Create(ctx context.Context, cp *models.CoursePrerequisite) error {
	query := `INSERT INTO course_prerequisites (course_id, prerequisite_course_id) VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, cp.CourseID, cp.PrerequisiteCourseID)
	return err
}

func (r *coursePrerequisiteRepository) Delete(ctx context.Context, courseID uint64, prereqID uint64) error {
	query := `DELETE FROM course_prerequisites WHERE course_id = ? AND prerequisite_course_id = ?`
	_, err := r.db.ExecContext(ctx, query, courseID, prereqID)
	return err
}

func (r *coursePrerequisiteRepository) GetPrerequisitesByCourseID(ctx context.Context, courseID uint64) ([]models.CoursePrerequisite, error) {
	query := `SELECT id, course_id, prerequisite_course_id FROM course_prerequisites WHERE course_id = ?`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.CoursePrerequisite
	for rows.Next() {
		var cp models.CoursePrerequisite
		if err := rows.Scan(&cp.ID, &cp.CourseID, &cp.PrerequisiteCourseID); err != nil {
			return nil, err
		}
		results = append(results, cp)
	}
	return results, nil
}

func (r *coursePrerequisiteRepository) GetPrerequisitesForCourse(ctx context.Context, courseID uint64) ([]uint64, error) {
	query := `SELECT prerequisite_course_id FROM course_prerequisites WHERE course_id = ?`
	rows, err := r.db.QueryContext(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *coursePrerequisiteRepository) CheckPrerequisitesMet(ctx context.Context, userID uint64, courseID uint64) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM course_prerequisites cp
		JOIN passed_courses pc ON cp.prerequisite_course_id = pc.course_id
		WHERE cp.course_id = ? AND pc.user_id = ?`

	var passedCount int
	err := r.db.QueryRowContext(ctx, query, courseID, userID).Scan(&passedCount)
	if err != nil {
		return false, err
	}

	var totalPrereqs int
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM course_prerequisites WHERE course_id = ?`, courseID).Scan(&totalPrereqs)
	if err != nil {
		return false, err
	}

	return passedCount == totalPrereqs, nil
}