package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

// StudentRepository defines the interface for student data access
type StudentRepository interface {
	Create(ctx context.Context, student *models.Student) error
	GetByEmail(ctx context.Context, email string) (*models.Student, error)
	GetByID(ctx context.Context, id int) (*models.Student, error)
	Search(ctx context.Context, query string) ([]*models.StudentSearchResult, error)
}

// PostgresStudentRepository implements StudentRepository for PostgreSQL
type PostgresStudentRepository struct {
	db *sql.DB
}

// NewPostgresStudentRepository creates a new PostgreSQL student repository
func NewPostgresStudentRepository(db *sql.DB) StudentRepository {
	return &PostgresStudentRepository{db: db}
}

// Create inserts a new student into the database
func (r *PostgresStudentRepository) Create(ctx context.Context, student *models.Student) error {
	query := `
		INSERT INTO students (first_name, last_name, class_info, class_id, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		student.FirstName,
		student.LastName,
		student.ClassInfo,
		student.ClassID,
		student.Email,
	).Scan(&student.ID, &student.CreatedAt, &student.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation (duplicate email)
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "unique constraint") {
			return fmt.Errorf("email already exists: %w", err)
		}
		return fmt.Errorf("failed to create student: %w", err)
	}

	return nil
}

// GetByEmail retrieves a student by email address
func (r *PostgresStudentRepository) GetByEmail(ctx context.Context, email string) (*models.Student, error) {
	query := `
		SELECT id, first_name, last_name, class_info, class_id, email, created_at, updated_at
		FROM students
		WHERE email = $1
	`

	student := &models.Student{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.ClassInfo,
		&student.ClassID,
		&student.Email,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("student not found with email: %s", email)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get student by email: %w", err)
	}

	return student, nil
}

// GetByID retrieves a student by ID
func (r *PostgresStudentRepository) GetByID(ctx context.Context, id int) (*models.Student, error) {
	query := `
		SELECT id, first_name, last_name, class_info, class_id, email, created_at, updated_at
		FROM students
		WHERE id = $1
	`

	student := &models.Student{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.ClassInfo,
		&student.ClassID,
		&student.Email,
		&student.CreatedAt,
		&student.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("student not found with id: %d", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get student by id: %w", err)
	}

	return student, nil
}

// Search performs a case-insensitive search on first_name and last_name
func (r *PostgresStudentRepository) Search(ctx context.Context, query string) ([]*models.StudentSearchResult, error) {
	sqlQuery := `
		SELECT id, first_name, last_name, class_info, email
		FROM students
		WHERE LOWER(first_name) LIKE LOWER($1) OR LOWER(last_name) LIKE LOWER($1)
		ORDER BY first_name, last_name
		LIMIT 20
	`

	// Add wildcards for partial matching
	searchPattern := "%" + query + "%"

	rows, err := r.db.QueryContext(ctx, sqlQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search students: %w", err)
	}
	defer rows.Close()

	var results []*models.StudentSearchResult
	for rows.Next() {
		result := &models.StudentSearchResult{}
		err := rows.Scan(
			&result.ID,
			&result.FirstName,
			&result.LastName,
			&result.ClassInfo,
			&result.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student search result: %w", err)
		}

		// Construct full name
		result.FullName = fmt.Sprintf("%s %s", result.FirstName, result.LastName)
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating student search results: %w", err)
	}

	return results, nil
}
