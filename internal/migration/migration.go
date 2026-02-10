package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// RunMigrations applies all pending database migrations in order
// This is a simple forward-only migration runner
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Create schema_migrations table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("could not create schema_migrations table: %w", err)
	}

	// Get list of migration files
	files, err := filepath.Glob(filepath.Join(migrationsPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("could not read migration files: %w", err)
	}

	// Sort files to ensure they run in order
	sort.Strings(files)

	// Apply each migration
	for _, file := range files {
		filename := filepath.Base(file)

		// Check if migration has already been applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", filename).Scan(&count)
		if err != nil {
			return fmt.Errorf("could not check migration status for %s: %w", filename, err)
		}

		if count > 0 {
			// Migration already applied, skip
			continue
		}

		// Read migration file
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("could not read migration file %s: %w", filename, err)
		}

		// Execute migration
		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("could not execute migration %s: %w", filename, err)
		}

		// Record migration as applied
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", filename)
		if err != nil {
			return fmt.Errorf("could not record migration %s: %w", filename, err)
		}
	}

	return nil
}

// GetAppliedMigrations returns a list of applied migration versions
func GetAppliedMigrations(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("could not query applied migrations: %w", err)
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("could not scan migration version: %w", err)
		}
		migrations = append(migrations, version)
	}

	return migrations, nil
}
