package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
)

type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type EnrollmentRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

var _ repositories.EnrollmentRepository = (*EnrollmentRepository)(nil)

func NewEnrollmentRepository(db *sql.DB, logger *zap.Logger) *EnrollmentRepository {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &EnrollmentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *EnrollmentRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return nil, err
	}

	return tx, nil
}

func (r *EnrollmentRepository) GetUserCreditInfo(
	ctx context.Context,
	tx *sql.Tx,
	userID uint64,
) (models.UserCreditInfo, error) {
	const query = `
		SELECT id, max_sks
		FROM users
		WHERE id = ?
	`

	var info models.UserCreditInfo
	err := tx.QueryRowContext(ctx, query, userID).Scan(&info.ID, &info.MaxSKS)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.UserCreditInfo{}, models.ErrUserNotFound
		}

		r.logger.Error("failed to fetch user credit info", zap.Uint64("user_id", userID), zap.Error(err))
		return models.UserCreditInfo{}, err
	}

	return info, nil
}

func (r *EnrollmentRepository) SumCurrentSKSByUserID(
	ctx context.Context,
	tx *sql.Tx,
	userID uint64,
) (int, error) {
	const query = `
		SELECT COALESCE(SUM(c.sks), 0)
		FROM enrollments e
		INNER JOIN courses c ON c.id = e.course_id
		WHERE e.user_id = ?
	`

	var total int
	if err := tx.QueryRowContext(ctx, query, userID).Scan(&total); err != nil {
		r.logger.Error("failed to sum current sks", zap.Uint64("user_id", userID), zap.Error(err))
		return 0, err
	}

	return total, nil
}

func (r *EnrollmentRepository) GetCourseByID(
	ctx context.Context,
	tx *sql.Tx,
	courseID uint64,
) (models.Course, error) {
	const query = `
		SELECT id, code, name, sks, quota
		FROM courses
		WHERE id = ?
		FOR UPDATE
	`

	var course models.Course
	err := tx.QueryRowContext(ctx, query, courseID).
		Scan(&course.ID, &course.Code, &course.Name, &course.SKS, &course.Quota)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Course{}, models.ErrCourseNotFound
		}

		r.logger.Error("failed to fetch course", zap.Uint64("course_id", courseID), zap.Error(err))
		return models.Course{}, err
	}

	return course, nil
}

func (r *EnrollmentRepository) GetCourseSchedulesByCourseID(
	ctx context.Context,
	tx *sql.Tx,
	courseID uint64,
) ([]models.ScheduleSlot, error) {
	const query = `
		SELECT
			s.course_id,
			c.code,
			s.day_of_week,
			TIME_FORMAT(s.start_time, '%H:%i:%s') AS start_time,
			TIME_FORMAT(s.end_time, '%H:%i:%s') AS end_time
		FROM schedules s
		INNER JOIN courses c ON c.id = s.course_id
		WHERE s.course_id = ?
		ORDER BY s.day_of_week, s.start_time
	`

	return r.scanSchedules(ctx, tx, query, courseID)
}

func (r *EnrollmentRepository) GetTakenSchedulesByUserID(
	ctx context.Context,
	tx *sql.Tx,
	userID uint64,
) ([]models.ScheduleSlot, error) {
	const query = `
		SELECT
			s.course_id,
			c.code,
			s.day_of_week,
			TIME_FORMAT(s.start_time, '%H:%i:%s') AS start_time,
			TIME_FORMAT(s.end_time, '%H:%i:%s') AS end_time
		FROM enrollments e
		INNER JOIN schedules s ON s.course_id = e.course_id
		INNER JOIN courses c ON c.id = e.course_id
		WHERE e.user_id = ?
		ORDER BY s.day_of_week, s.start_time
	`

	return r.scanSchedules(ctx, tx, query, userID)
}

func (r *EnrollmentRepository) CountCurrentEnrollmentsByCourseID(
	ctx context.Context,
	tx *sql.Tx,
	courseID uint64,
) (int, error) {
	const query = `
		SELECT COUNT(*)
		FROM enrollments
		WHERE course_id = ?
	`

	var total int
	if err := tx.QueryRowContext(ctx, query, courseID).Scan(&total); err != nil {
		r.logger.Error("failed to count current quota", zap.Uint64("course_id", courseID), zap.Error(err))
		return 0, err
	}

	return total, nil
}

func (r *EnrollmentRepository) HasEnrollment(
	ctx context.Context,
	tx *sql.Tx,
	userID uint64,
	courseID uint64,
) (bool, error) {
	const query = `
		SELECT COUNT(*)
		FROM enrollments
		WHERE user_id = ? AND course_id = ?
	`

	var total int
	if err := tx.QueryRowContext(ctx, query, userID, courseID).Scan(&total); err != nil {
		r.logger.Error(
			"failed to check enrollment existence",
			zap.Uint64("user_id", userID),
			zap.Uint64("course_id", courseID),
			zap.Error(err),
		)
		return false, err
	}

	return total > 0, nil
}

func (r *EnrollmentRepository) CreateEnrollment(
	ctx context.Context,
	tx *sql.Tx,
	enrollment models.Enrollment,
) (*models.Enrollment, error) {
	const query = `
		INSERT INTO enrollments (user_id, course_id, enrolled_at)
		VALUES (?, ?, ?)
	`

	enrolledAt := enrollment.EnrolledAt
	if enrolledAt.IsZero() {
		enrolledAt = time.Now().UTC()
	}

	result, err := tx.ExecContext(ctx, query, enrollment.UserID, enrollment.CourseID, enrolledAt)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, models.ErrAlreadyEnrolled
		}

		r.logger.Error(
			"failed to create enrollment",
			zap.Uint64("user_id", enrollment.UserID),
			zap.Uint64("course_id", enrollment.CourseID),
			zap.Error(err),
		)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error("failed to read enrollment insert id", zap.Error(err))
		return nil, err
	}

	return &models.Enrollment{
		ID:         uint64(id),
		UserID:     enrollment.UserID,
		CourseID:   enrollment.CourseID,
		EnrolledAt: enrolledAt,
	}, nil
}

func (r *EnrollmentRepository) scanSchedules(
	ctx context.Context,
	exec sqlExecutor,
	query string,
	arg any,
) ([]models.ScheduleSlot, error) {
	rows, err := exec.QueryContext(ctx, query, arg)
	if err != nil {
		r.logger.Error("failed to query schedules", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	schedules := make([]models.ScheduleSlot, 0)
	for rows.Next() {
		var slot models.ScheduleSlot
		if err := rows.Scan(
			&slot.CourseID,
			&slot.CourseCode,
			&slot.Day,
			&slot.StartTime,
			&slot.EndTime,
		); err != nil {
			r.logger.Error("failed to scan schedule row", zap.Error(err))
			return nil, err
		}

		schedules = append(schedules, slot)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("failed while iterating schedule rows", zap.Error(err))
		return nil, err
	}

	return schedules, nil
}

func (r *EnrollmentRepository) String() string {
	return fmt.Sprintf("EnrollmentRepository{db:%t}", r.db != nil)
}

func (r *EnrollmentRepository) ListAllEnrollments(ctx context.Context) ([]models.Enrollment, error) {
	const query = `SELECT id, user_id, course_id, enrolled_at, created_at, updated_at FROM enrollments`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error("failed to list all enrollments", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	results := make([]models.Enrollment, 0)
	for rows.Next() {
		var e models.Enrollment
		if err := rows.Scan(&e.ID, &e.UserID, &e.CourseID, &e.EnrolledAt, &e.CreatedAt, &e.UpdatedAt); err != nil {
			r.logger.Error("failed to scan enrollment row", zap.Error(err))
			return nil, err
		}
		results = append(results, e)
	}

	return results, nil
}

func (r *EnrollmentRepository) UpdateStatus(ctx context.Context, enrollmentID uint64, status string) error {
	const query = `UPDATE enrollments SET status = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, enrollmentID)
	if err != nil {
		r.logger.Error("failed to update enrollment status", zap.Uint64("enrollment_id", enrollmentID), zap.Error(err))
		return err
	}
	return nil
}
