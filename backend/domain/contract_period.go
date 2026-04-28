package domain

import (
	"context"
	"time"
)

// ContractPeriod defines the time range when students can register for courses.
type ContractPeriod struct {
	ID        int       `json:"id"         db:"id"`
	IsOpen    bool      `json:"is_open"    db:"is_open"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date"   db:"end_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ContractPeriodRepository defines the storage operations for contract periods.
type ContractPeriodRepository interface {
	Get(ctx context.Context) (*ContractPeriod, error)
	Update(ctx context.Context, period *ContractPeriod) error
}

// ContractPeriodUsecase defines the business logic for contract periods.
type ContractPeriodUsecase interface {
	GetCurrent(ctx context.Context) (*ContractPeriod, error)
	Update(ctx context.Context, isOpen bool, startDate, endDate string) error
}
