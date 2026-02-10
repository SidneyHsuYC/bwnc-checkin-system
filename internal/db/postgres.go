package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"

	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := &Config{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}

	// Validate required fields
	if config.User == "" || config.Password == "" || config.DBName == "" {
		logger.Error("Missing required database configuration (DB_USER, DB_PASSWORD, DB_NAME)")
		os.Exit(1)
	}

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// NewPostgres creates a new PostgreSQL database connection with connection pooling
func NewPostgres() *sql.DB {
	config := LoadConfigFromEnv()
	return NewPostgresWithConfig(config)
}

// NewPostgresWithConfig creates a new PostgreSQL connection with the provided config
func NewPostgresWithConfig(config *Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	logger.Info("Connecting to database",
		"host", config.Host,
		"port", config.Port,
		"dbname", config.DBName,
		"user", config.User)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open database connection", "error", err)
		os.Exit(1)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection
	db.SetConnMaxIdleTime(2 * time.Minute) // Maximum idle time of a connection

	// Test connection with exponential backoff retry
	if err := connectWithRetry(db); err != nil {
		logger.Error("Failed to connect to database after retries", "error", err)
		os.Exit(1)
	}

	logger.Info("Database connection established successfully")
	return db
}

// connectWithRetry attempts to connect to the database with exponential backoff
func connectWithRetry(db *sql.DB) error {
	maxRetries := 5
	baseDelay := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		err := db.Ping()
		if err == nil {
			logger.Info("Database ping successful")
			return nil
		}

		// Calculate exponential backoff delay: 1s, 2s, 4s, 8s, 16s
		delay := baseDelay * time.Duration(1<<uint(i))
		logger.Warn("Database ping failed, retrying...",
			"attempt", i+1,
			"max_retries", maxRetries,
			"next_retry_in", delay,
			"error", err)

		if i < maxRetries-1 {
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

// HealthCheck performs a health check on the database connection
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}
