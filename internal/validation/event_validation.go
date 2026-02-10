package validation

import (
	"fmt"
	"strings"
	"time"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

// ValidateEvent validates an event model
func ValidateEvent(event *models.Event) error {
	var errors []string

	// Validate event name
	if strings.TrimSpace(event.EventName) == "" {
		errors = append(errors, "event_name is required")
	}

	// Validate event type
	if strings.TrimSpace(event.EventType) == "" {
		errors = append(errors, "event_type is required")
	}

	// Validate event time
	if event.EventTime.IsZero() {
		errors = append(errors, "event_time is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

// SanitizeEvent trims whitespace from event fields
func SanitizeEvent(event *models.Event) {
	event.EventName = strings.TrimSpace(event.EventName)
	event.EventType = strings.TrimSpace(event.EventType)
}

// IsValidEventTime checks if an event time is valid (not in the past)
func IsValidEventTime(eventTime time.Time) bool {
	return !eventTime.Before(time.Now().Add(-24 * time.Hour)) // Allow events from yesterday
}
