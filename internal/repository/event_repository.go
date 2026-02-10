package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

// EventRepository defines the interface for event data access
type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	GetByID(ctx context.Context, id int) (*models.Event, error)
	List(ctx context.Context) ([]*models.Event, error)
	ListUpcoming(ctx context.Context) ([]*models.Event, error)
}

// PostgresEventRepository implements EventRepository for PostgreSQL
type PostgresEventRepository struct {
	db *sql.DB
}

// NewPostgresEventRepository creates a new PostgreSQL event repository
func NewPostgresEventRepository(db *sql.DB) EventRepository {
	return &PostgresEventRepository{db: db}
}

// Create inserts a new event into the database
func (r *PostgresEventRepository) Create(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (event_name, event_time, event_type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		event.EventName,
		event.EventTime,
		event.EventType,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// GetByID retrieves an event by ID
func (r *PostgresEventRepository) GetByID(ctx context.Context, id int) (*models.Event, error) {
	query := `
		SELECT id, event_name, event_time, event_type, created_at, updated_at
		FROM events
		WHERE id = $1
	`

	event := &models.Event{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.EventName,
		&event.EventTime,
		&event.EventType,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event not found with id: %d", id)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}

	return event, nil
}

// List retrieves all events ordered by event time
func (r *PostgresEventRepository) List(ctx context.Context) ([]*models.Event, error) {
	query := `
		SELECT id, event_name, event_time, event_type, created_at, updated_at
		FROM events
		ORDER BY event_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event := &models.Event{}
		err := rows.Scan(
			&event.ID,
			&event.EventName,
			&event.EventTime,
			&event.EventType,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// ListUpcoming retrieves all upcoming events (event_time >= now)
func (r *PostgresEventRepository) ListUpcoming(ctx context.Context) ([]*models.Event, error) {
	query := `
		SELECT id, event_name, event_time, event_type, created_at, updated_at
		FROM events
		WHERE event_time >= $1
		ORDER BY event_time ASC
	`

	now := time.Now()
	rows, err := r.db.QueryContext(ctx, query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to list upcoming events: %w", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		event := &models.Event{}
		err := rows.Scan(
			&event.ID,
			&event.EventName,
			&event.EventTime,
			&event.EventType,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating upcoming events: %w", err)
	}

	return events, nil
}
