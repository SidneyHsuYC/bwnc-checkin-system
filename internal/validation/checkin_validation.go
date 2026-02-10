package validation

import (
	"fmt"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

// ValidateCheckin validates check-in data
func ValidateCheckin(checkin *models.Checkin) error {
	if checkin.StudentID <= 0 {
		return fmt.Errorf("student ID is required")
	}

	if checkin.EventID <= 0 {
		return fmt.Errorf("event ID is required")
	}

	return nil
}
