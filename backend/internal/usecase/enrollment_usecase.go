package usecase

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
	"sistemkontrakmatkul/backend/internal/domain/services"
)

type EnrollmentUsecase struct {
	repo         repositories.EnrollmentRepository
	passedRepo   repositories.PassedCourseRepository
	prereqRepo   repositories.CoursePrerequisiteRepository
	logger       *zap.Logger
}

var _ services.EnrollmentService = (*EnrollmentUsecase)(nil)

func NewEnrollmentUsecase(
	repo repositories.EnrollmentRepository,
	passedRepo repositories.PassedCourseRepository,
	prereqRepo repositories.CoursePrerequisiteRepository,
	logger *zap.Logger,
) *EnrollmentUsecase {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &EnrollmentUsecase{
		repo:       repo,
		passedRepo: passedRepo,
		prereqRepo: prereqRepo,
		logger:     logger,
	}
}

func (u *EnrollmentUsecase) Enroll(
	ctx context.Context,
	request models.EnrollmentRequest,
) (*models.EnrollmentResult, error) {
	if request.UserID == 0 || request.CourseID == 0 {
		return nil, models.ErrInvalidEnrollmentRequest
	}

	tx, err := u.repo.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if committed {
			return
		}

		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			u.logger.Error("failed to rollback enrollment transaction", zap.Error(rollbackErr))
		}
	}()

	userCreditInfo, err := u.repo.GetUserCreditInfo(ctx, tx, request.UserID)
	if err != nil {
		return nil, err
	}

	course, err := u.repo.GetCourseByID(ctx, tx, request.CourseID)
	if err != nil {
		return nil, err
	}

	// 1. Check Prerequisites
	prereqs, err := u.prereqRepo.GetPrerequisitesForCourse(ctx, request.CourseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check prerequisites: %w", err)
	}

	if len(prereqs) > 0 {
		hasAllPassed, err := u.passedRepo.HasPassedCourses(ctx, request.UserID, prereqs)
		if err != nil {
			return nil, fmt.Errorf("failed to check passed courses: %w", err)
		}
		if !hasAllPassed {
			return nil, models.ErrPrerequisiteNotMet
		}
	}

	alreadyEnrolled, err := u.repo.HasEnrollment(ctx, tx, request.UserID, request.CourseID)
	if err != nil {
		return nil, err
	}
	if alreadyEnrolled {
		return nil, models.ErrAlreadyEnrolled
	}

	currentSKS, err := u.repo.SumCurrentSKSByUserID(ctx, tx, request.UserID)
	if err != nil {
		return nil, err
	}

	if currentSKS+course.SKS > userCreditInfo.MaxSKS {
		return nil, fmt.Errorf(
			"%w: current=%d, course=%d, max=%d",
			models.ErrCreditLimitExceeded,
			currentSKS,
			course.SKS,
			userCreditInfo.MaxSKS,
		)
	}

	currentQuota, err := u.repo.CountCurrentEnrollmentsByCourseID(ctx, tx, request.CourseID)
	if err != nil {
		return nil, err
	}
	if currentQuota >= course.Quota {
		return nil, fmt.Errorf(
			"%w: current=%d, quota=%d",
			models.ErrQuotaExceeded,
			currentQuota,
			course.Quota,
		)
	}

	takenSchedules, err := u.repo.GetTakenSchedulesByUserID(ctx, tx, request.UserID)
	if err != nil {
		return nil, err
	}

	newCourseSchedules, err := u.repo.GetCourseSchedulesByCourseID(ctx, tx, request.CourseID)
	if err != nil {
		return nil, err
	}

	if hasScheduleConflict(takenSchedules, newCourseSchedules) {
		return nil, models.ErrScheduleConflict
	}

	enrollment, err := u.repo.CreateEnrollment(ctx, tx, models.Enrollment{
		UserID:     request.UserID,
		CourseID:   request.CourseID,
		EnrolledAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		u.logger.Error(
			"failed to commit enrollment transaction",
			zap.Uint64("user_id", request.UserID),
			zap.Uint64("course_id", request.CourseID),
			zap.Error(err),
		)
		return nil, err
	}

	committed = true

	u.logger.Info(
		"enrollment committed successfully",
		zap.Uint64("user_id", request.UserID),
		zap.Uint64("course_id", request.CourseID),
	)

	return &models.EnrollmentResult{
		Enrollment: *enrollment,
	}, nil
}

func (u *EnrollmentUsecase) ListAll(ctx context.Context) ([]models.Enrollment, error) {
	return u.repo.ListAllEnrollments(ctx)
}

func (u *EnrollmentUsecase) Approve(ctx context.Context, enrollmentID uint64) error {
	return u.repo.UpdateStatus(ctx, enrollmentID, "APPROVED")
}

func (u *EnrollmentUsecase) Reject(ctx context.Context, enrollmentID uint64) error {
	return u.repo.UpdateStatus(ctx, enrollmentID, "REJECTED")
}

func hasScheduleConflict(
	currentSchedules []models.ScheduleSlot,
	newSchedules []models.ScheduleSlot,
) bool {
	for _, existing := range currentSchedules {
		existingStart, err := parseClock(existing.StartTime)
		if err != nil {
			return true
		}

		existingEnd, err := parseClock(existing.EndTime)
		if err != nil {
			return true
		}

		for _, candidate := range newSchedules {
			if !strings.EqualFold(existing.Day, candidate.Day) {
				continue
			}

			candidateStart, err := parseClock(candidate.StartTime)
			if err != nil {
				return true
			}

			candidateEnd, err := parseClock(candidate.EndTime)
			if err != nil {
				return true
			}

			if candidateStart.Before(existingEnd) && candidateEnd.After(existingStart) {
				return true
			}
		}
	}

	return false
}

func parseClock(value string) (time.Time, error) {
	layouts := []string{"15:04:05", "15:04"}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid clock value: %s", value)
}