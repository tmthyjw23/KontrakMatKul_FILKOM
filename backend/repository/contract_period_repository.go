package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"kontrak-matkul/domain"
)

type contractPeriodRepository struct {
	db *sql.DB
}

func NewContractPeriodRepository(db *sql.DB) domain.ContractPeriodRepository {
	return &contractPeriodRepository{db: db}
}

func (r *contractPeriodRepository) Get(ctx context.Context) (*domain.ContractPeriod, error) {
	query := `SELECT id, is_open, start_date, end_date, created_at, updated_at FROM contract_periods LIMIT 1`
	
	row := r.db.QueryRowContext(ctx, query)
	
	var p domain.ContractPeriod
	var startDate, endDate, createdAt, updatedAt string
	
	err := row.Scan(&p.ID, &p.IsOpen, &startDate, &endDate, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching contract period: %w", err)
	}
	
	p.StartDate, _ = parseTime(startDate)
	p.EndDate, _ = parseTime(endDate)
	p.CreatedAt, _ = parseTime(createdAt)
	p.UpdatedAt, _ = parseTime(updatedAt)
	
	return &p, nil
}

func (r *contractPeriodRepository) Update(ctx context.Context, p *domain.ContractPeriod) error {
	query := `UPDATE contract_periods SET is_open = ?, start_date = ?, end_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, p.IsOpen, p.StartDate, p.EndDate, p.ID)
	if err != nil {
		return fmt.Errorf("error updating contract period: %w", err)
	}
	
	return nil
}

// parseTime helper (reusing logic if available or just basic parsing)
func parseTime(s string) (time.Time, error) {
    // MySQL DATETIME format
    layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "2006-01-02T15:04:05+07:00"}
    var t time.Time
    var err error
    for _, l := range layouts {
        t, err = time.Parse(l, s)
        if err == nil {
            return t, nil
        }
    }
    return time.Time{}, err
}
