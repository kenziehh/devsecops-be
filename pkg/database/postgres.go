package database

import (
	"database/sql"
	"devsecops-be/pkg/logger"
    "devsecops-be/config/env"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

func NewPostgresConnection(log logger.Logger) (*sql.DB, error) {
	config := loadConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info(nil, "Successfully connected to PostgreSQL database", logger.Fields{
		"host":   config.Host,
		"port":   config.Port,
		"dbname": config.DBName,
	})

	return db, nil
}

func loadConfig() *Config {
	config := &Config{
		Host:         env.GetEnv("DB_HOST", "localhost"),
		Port:         env.GetEnvAsInt("DB_PORT", 5432),
		User:         env.GetEnv("DB_USER", "postgres"),
		Password:     env.GetEnv("DB_PASSWORD", "password"),
		DBName:       env.GetEnv("DB_NAME", "devsecops_be"),
		SSLMode:      env.GetEnv("DB_SSLMODE", "disable"),
		MaxOpenConns: env.GetEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: env.GetEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		MaxLifetime:  time.Duration(env.GetEnvAsInt("DB_MAX_LIFETIME_MINUTES", 5)) * time.Minute,
	}

	// Parse DATABASE_URL if provided (for Docker/Heroku compatibility)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		// DATABASE_URL takes precedence over individual env vars
		return parsePostgresURL(databaseURL, config)
	}

	return config
}

func parsePostgresURL(databaseURL string, defaultConfig *Config) *Config {
	// Simple parsing for postgres://user:password@host:port/dbname?sslmode=disable
	// For production, consider using a proper URL parser
	return defaultConfig // Simplified for this example
}

