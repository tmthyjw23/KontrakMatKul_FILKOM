package mysql

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
)

type CourseRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

var _ repositories.CourseRepository = (*CourseRepository)(nil)

func NewCourseRepository(db *sql.DB, logger *zap.Logger) *CourseRepository {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &CourseRepository{
		db:     db,
		logger: logger,
	}
}

func (r *CourseRepository) ListCourses(ctx context.Context) ([]models.CourseResponse, error) {
	const query = `
		SELECT
			c.id,
			c.code,
			c.name,
			c.sks,
			c.quota,
			c.lecturer,
			s.day_of_week,
			TIME_FORMAT(s.start_time, '%H:%i:%s') AS start_time,
			TIME_FORMAT(s.end_time, '%H:%i:%s') AS end_time,
			s.room
		FROM courses c
		LEFT JOIN schedules s ON s.course_id = c.id
		ORDER BY c.code, s.day_of_week, s.start_time
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error("failed to query courses", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	courseMap := make(map[uint64]*models.CourseResponse)
	courseOrder := make([]uint64, 0)

	for rows.Next() {
		var (
			id        uint64
			code      string
			name      string
			sks       int
			quota     int
			lecturer  string
			dayOfWeek sql.NullString
			startTime sql.NullString
			endTime   sql.NullString
			room      sql.NullString
		)

		if err := rows.Scan(
			&id,
			&code,
			&name,
			&sks,
			&quota,
			&lecturer,
			&dayOfWeek,
			&startTime,
			&endTime,
			&room,
		); err != nil {
			r.logger.Error("failed to scan course row", zap.Error(err))
			return nil, err
		}

		course, exists := courseMap[id]
		if !exists {
			course = &models.CourseResponse{
				ID:        id,
				Code:      code,
				Name:      name,
				SKS:       sks,
				Quota:     quota,
				Lecturer:  lecturer,
				Schedules: make([]models.CourseScheduleResponse, 0),
			}
			courseMap[id] = course
			courseOrder = append(courseOrder, id)
		}

		if dayOfWeek.Valid && startTime.Valid && endTime.Valid {
			course.Schedules = append(course.Schedules, models.CourseScheduleResponse{
				DayOfWeek: dayOfWeek.String,
				StartTime: startTime.String,
				EndTime:   endTime.String,
				Room:      room.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed while iterating course rows", zap.Error(err))
		return nil, err
	}

	courses := make([]models.CourseResponse, 0, len(courseOrder))
	for _, id := range courseOrder {
		courses = append(courses, *courseMap[id])
	}

	return courses, nil
}

func (r *CourseRepository) Create(ctx context.Context, course *models.Course) error {
	const query = `
		INSERT INTO courses (code, name, sks, quota, lecturer)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, course.Code, course.Name, course.SKS, course.Quota, course.Lecturer)
	if err != nil {
		r.logger.Error("failed to create course", zap.Error(err))
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error("failed to retrieve insert id", zap.Error(err))
		return err
	}
	course.ID = uint64(id)

	return nil
}

func (r *CourseRepository) Update(ctx context.Context, course *models.Course) error {
	const query = `
		UPDATE courses 
		SET code = ?, name = ?, sks = ?, quota = ?, lecturer = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, course.Code, course.Name, course.SKS, course.Quota, course.Lecturer, course.ID)
	if err != nil {
		r.logger.Error("failed to update course", zap.Uint64("course_id", course.ID), zap.Error(err))
		return err
	}

	return nil
}

func (r *CourseRepository) Delete(ctx context.Context, id uint64) error {
	const query = `DELETE FROM courses WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete course", zap.Uint64("course_id", id), zap.Error(err))
		return err
	}

	return nil
}
