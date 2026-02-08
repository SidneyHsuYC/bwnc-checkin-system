package db

import (
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgres() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logger.Error("DATABASE_URL not set in environment")
		os.Exit(1)
	}

	logger.Info("Connecting to database", "dsn", maskPassword(dsn))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open database connection", "error", err)
		os.Exit(1)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection with retry
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			logger.Info("Database ping successful")
			return db
		}
		logger.Warn("Database ping failed", "attempt", i+1, "max_retries", maxRetries, "error", err)
		time.Sleep(2 * time.Second)
	}

	logger.Error("Failed to connect to database", "attempts", maxRetries, "error", err)
	os.Exit(1)
	return nil
}

// maskPassword hides the password in connection string for logging
func maskPassword(dsn string) string {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		// If parsing fails, return a generic masked string
		return "postgres://***:***@***/***/..."
	}

	// Mask the password
	if parsedURL.User != nil {
		username := parsedURL.User.Username()
		parsedURL.User = url.UserPassword(username, "***")
	}

	// Get the database name from path
	dbName := "default"
	if parsedURL.Path != "" && parsedURL.Path != "/" {
		dbName = parsedURL.Path[1:] // Remove leading slash
	}

	return fmt.Sprintf("postgres://%s@%s/%s",
		parsedURL.User.String(),
		parsedURL.Host,
		dbName)
}
