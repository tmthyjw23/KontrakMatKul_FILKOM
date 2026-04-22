package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"sistemkontrakmatkul/backend/internal/config"
)

func NewMySQL(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	if logger != nil {
		logger.Info(
			"mysql connection pool initialized",
			zap.Int("max_open_conns", cfg.DBMaxOpenConns),
			zap.Int("max_idle_conns", cfg.DBMaxIdleConns),
			zap.Duration("conn_max_lifetime", cfg.DBConnMaxLifetime),
		)
	}

	return db, nil
}
