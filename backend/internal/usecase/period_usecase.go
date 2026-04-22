package usecase

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
)

type PeriodUsecase struct {
	settingRepo repositories.SystemSettingsRepository
	logger      *zap.Logger
}

func NewPeriodUsecase(settingRepo repositories.SystemSettingsRepository, logger *zap.Logger) *PeriodUsecase {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &PeriodUsecase{
		settingRepo: settingRepo,
		logger:      logger,
	}
}

func (u *PeriodUsecase) GetPeriodStatus(ctx context.Context) (*models.PeriodStatusResponse, error) {
	val, err := u.settingRepo.GetSetting(ctx, "is_enrollment_open")
	if err != nil {
		return nil, err
	}
	
	isOpen, _ := strconv.ParseBool(val)
	return &models.PeriodStatusResponse{IsOpen: isOpen}, nil
}

func (u *PeriodUsecase) UpdatePeriodStatus(ctx context.Context, isOpen bool) error {
	val := strconv.FormatBool(isOpen)
	return u.settingRepo.SetSetting(ctx, "is_enrollment_open", val)
}
