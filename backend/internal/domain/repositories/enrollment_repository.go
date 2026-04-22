package repositories

import (
	"context"
	"database/sql"

	"sistemkontrakmatkul/backend/internal/domain/models"
)

type EnrollmentRepository interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	GetUserCreditInfo(ctx context.Context, tx *sql.Tx, userID uint64) (models.UserCreditInfo, error)
	SumCurrentSKSByUserID(ctx context.Context, tx *sql.Tx, userID uint64) (int, error)
	GetCourseByID(ctx context.Context, tx *sql.Tx, courseID uint64) (models.Course, error)
	GetCourseSchedulesByCourseID(ctx context.Context, tx *sql.Tx, courseID uint64) ([]models.ScheduleSlot, error)
	GetTakenSchedulesByUserID(ctx context.Context, tx *sql.Tx, userID uint64) ([]models.ScheduleSlot, error)
	CountCurrentEnrollmentsByCourseID(ctx context.Context, tx *sql.Tx, courseID uint64) (int, error)
	HasEnrollment(ctx context.Context, tx *sql.Tx, userID uint64, courseID uint64) (bool, error)
	CreateEnrollment(ctx context.Context, tx *sql.Tx, enrollment models.Enrollment) (*models.Enrollment, error)
}
