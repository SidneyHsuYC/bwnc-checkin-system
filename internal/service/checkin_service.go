package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/validation"
)

// CheckinService handles check-in business logic
type CheckinService struct {
	checkinRepo repository.CheckinRepository
	studentRepo repository.StudentRepository
	eventRepo   repository.EventRepository
}

// NewCheckinService creates a new check-in service
func NewCheckinService(db *sql.DB) *CheckinService {
	return &CheckinService{
		checkinRepo: repository.NewPostgresCheckinRepository(db),
		studentRepo: repository.NewPostgresStudentRepository(db),
		eventRepo:   repository.NewPostgresEventRepository(db),
	}
}

// Create creates a new check-in with validation
func (s *CheckinService) Create(ctx context.Context, checkin *models.Checkin) error {
	// Validate check-in data
	if err := validation.ValidateCheckin(checkin); err != nil {
		return err
	}

	// Verify student exists
	_, err := s.studentRepo.GetByID(ctx, checkin.StudentID)
	if err != nil {
		return fmt.Errorf("student not found")
	}

	// Verify event exists
	_, err = s.eventRepo.GetByID(ctx, checkin.EventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// Set check-in timestamp
	checkin.CheckedInAt = time.Now()

	// Create check-in
	if err := s.checkinRepo.Create(ctx, checkin); err != nil {
		return err
	}

	return nil
}

// ListByEvent retrieves all check-ins for a specific event
func (s *CheckinService) ListByEvent(ctx context.Context, eventID int) ([]*models.CheckinDetail, error) {
	if eventID <= 0 {
		return nil, fmt.Errorf("event ID is required")
	}

	return s.checkinRepo.ListByEvent(ctx, eventID)
}
