package mysql

import (
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type SystemSettingsRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewSystemSettingsRepository(db *sql.DB, logger *zap.Logger) *SystemSettingsRepository {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &SystemSettingsRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SystemSettingsRepository) GetSetting(ctx context.Context, key string) (string, error) {
	const query = `SELECT setting_value FROM system_settings WHERE setting_key = ?`
	var value string
	err := r.db.QueryRowContext(ctx, query, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		r.logger.Error("failed to get system setting", zap.String("key", key), zap.Error(err))
		return "", err
	}
	return value, nil
}

func (r *SystemSettingsRepository) SetSetting(ctx context.Context, key string, value string) error {
	const query = `INSERT INTO system_settings (setting_key, setting_value) VALUES (?, ?) 
	               ON DUPLICATE KEY UPDATE setting_value = VALUES(setting_value)`
	_, err := r.db.ExecContext(ctx, query, key, value)
	if err != nil {
		r.logger.Error("failed to set system setting", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}
