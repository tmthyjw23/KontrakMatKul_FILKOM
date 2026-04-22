package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv              string
	AppPort             string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	JWTSecret           string
	DBMaxOpenConns      int
	DBMaxIdleConns      int
	DBConnMaxLifetime   time.Duration
	ServerReadTimeout   time.Duration
	ServerWriteTimeout  time.Duration
	ServerIdleTimeout   time.Duration
	ShutdownGracePeriod time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:              getEnv("APP_ENV", "development"),
		AppPort:             getEnv("APP_PORT", "8080"),
		DBHost:              getEnv("DB_HOST", "127.0.0.1"),
		DBPort:              getEnv("DB_PORT", "3306"),
		DBUser:              getEnv("DB_USER", "root"),
		DBPassword:          getEnv("DB_PASSWORD", ""),
		DBName:              getEnv("DB_NAME", "sistem_kontrak"),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		DBMaxOpenConns:      getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:      getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime:   time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
		ServerReadTimeout:   time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT_SECONDS", 15)) * time.Second,
		ServerWriteTimeout:  time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT_SECONDS", 15)) * time.Second,
		ServerIdleTimeout:   time.Duration(getEnvAsInt("SERVER_IDLE_TIMEOUT_SECONDS", 60)) * time.Second,
		ShutdownGracePeriod: time.Duration(getEnvAsInt("SHUTDOWN_TIMEOUT_SECONDS", 10)) * time.Second,
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) Address() string {
	return ":" + c.AppPort
}

func (c *Config) MySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=false&charset=utf8mb4&collation=utf8mb4_unicode_ci&loc=UTC",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
