package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

// CheckinRepository defines the interface for check-in data access
type CheckinRepository interface {
	Create(ctx context.Context, checkin *models.Checkin) error
	GetByStudentAndEvent(ctx context.Context, studentID, eventID int) (*models.Checkin, error)
	ListByEvent(ctx context.Context, eventID int) ([]*models.CheckinDetail, error)
}

// PostgresCheckinRepository implements CheckinRepository for PostgreSQL
type PostgresCheckinRepository struct {
	db *sql.DB
}

// NewPostgresCheckinRepository creates a new PostgreSQL check-in repository
func NewPostgresCheckinRepository(db *sql.DB) CheckinRepository {
	return &PostgresCheckinRepository{db: db}
}

// Create inserts a new check-in into the database
func (r *PostgresCheckinRepository) Create(ctx context.Context, checkin *models.Checkin) error {
	query := `
		INSERT INTO check_ins (student_id, event_id, checked_in_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		checkin.StudentID,
		checkin.EventID,
		checkin.CheckedInAt,
	).Scan(&checkin.ID, &checkin.CreatedAt)

	if err != nil {
		// Check for unique constraint violation (duplicate check-in)
		if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "unique constraint") {
			return fmt.Errorf("student already checked in to this event")
		}
		// Check for foreign key constraint violation
		if strings.Contains(err.Error(), "foreign key constraint") {
			return fmt.Errorf("invalid student or event ID")
		}
		return fmt.Errorf("failed to create check-in: %w", err)
	}

	return nil
}

// GetByStudentAndEvent retrieves a check-in by student and event IDs
func (r *PostgresCheckinRepository) GetByStudentAndEvent(ctx context.Context, studentID, eventID int) (*models.Checkin, error) {
	query := `
		SELECT id, student_id, event_id, checked_in_at, created_at
		FROM check_ins
		WHERE student_id = $1 AND event_id = $2
	`

	checkin := &models.Checkin{}
	err := r.db.QueryRowContext(ctx, query, studentID, eventID).Scan(
		&checkin.ID,
		&checkin.StudentID,
		&checkin.EventID,
		&checkin.CheckedInAt,
		&checkin.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("check-in not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get check-in: %w", err)
	}

	return checkin, nil
}

// ListByEvent retrieves all check-ins for a specific event with student details
func (r *PostgresCheckinRepository) ListByEvent(ctx context.Context, eventID int) ([]*models.CheckinDetail, error) {
	query := `
		SELECT 
			c.id, c.checked_in_at,
			s.id, s.first_name, s.last_name, s.class_info, s.email, s.created_at, s.updated_at,
			e.id, e.event_name, e.event_time, e.event_type, e.created_at, e.updated_at
		FROM check_ins c
		INNER JOIN students s ON c.student_id = s.id
		INNER JOIN events e ON c.event_id = e.id
		WHERE c.event_id = $1
		ORDER BY c.checked_in_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to list check-ins by event: %w", err)
	}
	defer rows.Close()

	checkins := make([]*models.CheckinDetail, 0)
	for rows.Next() {
		detail := &models.CheckinDetail{
			Student: models.Student{},
			Event:   models.Event{},
		}
		err := rows.Scan(
			&detail.ID,
			&detail.CheckedInAt,
			&detail.Student.ID,
			&detail.Student.FirstName,
			&detail.Student.LastName,
			&detail.Student.ClassInfo,
			&detail.Student.Email,
			&detail.Student.CreatedAt,
			&detail.Student.UpdatedAt,
			&detail.Event.ID,
			&detail.Event.EventName,
			&detail.Event.EventTime,
			&detail.Event.EventType,
			&detail.Event.CreatedAt,
			&detail.Event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan check-in detail: %w", err)
		}
		checkins = append(checkins, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating check-ins: %w", err)
	}

	return checkins, nil
}
