package usecase

import (
	"context"
	"fmt"
	"time"

	"kontrak-matkul/domain"
)

type contractPeriodUsecase struct {
	repo domain.ContractPeriodRepository
}

func NewContractPeriodUsecase(repo domain.ContractPeriodRepository) domain.ContractPeriodUsecase {
	return &contractPeriodUsecase{repo: repo}
}

func (u *contractPeriodUsecase) GetCurrent(ctx context.Context) (*domain.ContractPeriod, error) {
	return u.repo.Get(ctx)
}

func (u *contractPeriodUsecase) Update(ctx context.Context, isOpen bool, startDate, endDate string) error {
	p, err := u.repo.Get(ctx)
	if err != nil {
		return err
	}
	if p == nil {
		return fmt.Errorf("contract period record not found in database")
	}

	p.IsOpen = isOpen
	
	if startDate != "" {
		t, err := time.Parse(time.RFC3339, startDate)
		if err == nil {
			p.StartDate = t
		}
	}
	
	if endDate != "" {
		t, err := time.Parse(time.RFC3339, endDate)
		if err == nil {
			p.EndDate = t
		}
	}

	return u.repo.Update(ctx, p)
}
