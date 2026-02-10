package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateStudent validates a student model
func ValidateStudent(student *models.Student) error {
	var errors []string

	// Validate first name
	if strings.TrimSpace(student.FirstName) == "" {
		errors = append(errors, "first_name is required")
	}

	// Validate last name
	if strings.TrimSpace(student.LastName) == "" {
		errors = append(errors, "last_name is required")
	}

	// Validate class info
	if strings.TrimSpace(student.ClassInfo) == "" {
		errors = append(errors, "class_info is required")
	}

	// Validate email
	if strings.TrimSpace(student.Email) == "" {
		errors = append(errors, "email is required")
	} else if !IsValidEmail(student.Email) {
		errors = append(errors, "email format is invalid")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

// IsValidEmail checks if an email address is valid
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) == 0 || len(email) > 255 {
		return false
	}
	return emailRegex.MatchString(email)
}

// SanitizeStudent trims whitespace from student fields
func SanitizeStudent(student *models.Student) {
	student.FirstName = strings.TrimSpace(student.FirstName)
	student.LastName = strings.TrimSpace(student.LastName)
	student.ClassInfo = strings.TrimSpace(student.ClassInfo)
	student.Email = strings.TrimSpace(strings.ToLower(student.Email))
}
