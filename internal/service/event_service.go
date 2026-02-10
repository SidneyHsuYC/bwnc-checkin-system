package service

import (
	"context"
	"fmt"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/validation"
)

// EventService handles business logic for events
type EventService struct {
	repo repository.EventRepository
}

// NewEventService creates a new event service
func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

// Create creates a new event after validation
func (s *EventService) Create(ctx context.Context, event *models.Event) error {
	// Sanitize input
	validation.SanitizeEvent(event)

	// Validate event data
	if err := validation.ValidateEvent(event); err != nil {
		logger.Warn("Event validation failed", "error", err)
		return fmt.Errorf("validation error: %w", err)
	}

	// Create event
	if err := s.repo.Create(ctx, event); err != nil {
		logger.Error("Failed to create event", "error", err)
		return fmt.Errorf("failed to create event: %w", err)
	}

	logger.Info("Event created successfully", "id", event.ID, "name", event.EventName)
	return nil
}

// GetByID retrieves an event by ID
func (s *EventService) GetByID(ctx context.Context, id int) (*models.Event, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid event ID: %d", id)
	}

	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// List retrieves all events
func (s *EventService) List(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.List(ctx)
	if err != nil {
		logger.Error("Failed to list events", "error", err)
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	logger.Info("Events listed successfully", "count", len(events))
	return events, nil
}

// ListUpcoming retrieves all upcoming events
func (s *EventService) ListUpcoming(ctx context.Context) ([]*models.Event, error) {
	events, err := s.repo.ListUpcoming(ctx)
	if err != nil {
		logger.Error("Failed to list upcoming events", "error", err)
		return nil, fmt.Errorf("failed to list upcoming events: %w", err)
	}

	logger.Info("Upcoming events listed successfully", "count", len(events))
	return events, nil
}
