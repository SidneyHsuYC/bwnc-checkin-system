package validation

import (
	"testing"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
)

func TestValidateStudent(t *testing.T) {
	tests := []struct {
		name    string
		student *models.Student
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid student with all fields",
			student: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				ClassInfo: "CS 101",
				Email:     "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name: "valid student without class info (optional)",
			student: &models.Student{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing first name",
			student: &models.Student{
				LastName: "Doe",
				Email:    "john.doe@example.com",
			},
			wantErr: true,
			errMsg:  "first_name is required",
		},
		{
			name: "missing last name",
			student: &models.Student{
				FirstName: "John",
				Email:     "john.doe@example.com",
			},
			wantErr: true,
			errMsg:  "last_name is required",
		},
		{
			name: "missing email",
			student: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
		{
			name: "invalid email format - no @",
			student: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalidemail.com",
			},
			wantErr: true,
			errMsg:  "email format is invalid",
		},
		{
			name: "invalid email format - no domain",
			student: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@",
			},
			wantErr: true,
			errMsg:  "email format is invalid",
		},
		{
			name: "invalid email format - no local part",
			student: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "@example.com",
			},
			wantErr: true,
			errMsg:  "email format is invalid",
		},
		{
			name: "whitespace only first name",
			student: &models.Student{
				FirstName: "   ",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			wantErr: true,
			errMsg:  "first_name is required",
		},
		{
			name: "email too long",
			student: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     string(make([]byte, 256)) + "@example.com",
			},
			wantErr: true,
			errMsg:  "email format is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStudent(tt.student)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStudent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateStudent() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@mail.example.com", true},
		{"valid email with plus", "test+tag@example.com", true},
		{"valid email with dash", "test-user@example.com", true},
		{"valid email with underscore", "test_user@example.com", true},
		{"valid email with numbers", "user123@example.com", true},
		{"invalid - no @", "testexample.com", false},
		{"invalid - no domain", "test@", false},
		{"invalid - no local part", "@example.com", false},
		{"invalid - no TLD", "test@example", false},
		{"invalid - spaces", "test user@example.com", false},
		{"invalid - empty", "", false},
		{"invalid - only @", "@", false},
		{"invalid - multiple @", "test@@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestSanitizeStudent(t *testing.T) {
	tests := []struct {
		name     string
		input    *models.Student
		expected *models.Student
	}{
		{
			name: "trim whitespace from all fields",
			input: &models.Student{
				FirstName: "  John  ",
				LastName:  "  Doe  ",
				ClassInfo: "  CS 101  ",
				Email:     "  JOHN.DOE@EXAMPLE.COM  ",
			},
			expected: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				ClassInfo: "CS 101",
				Email:     "john.doe@example.com",
			},
		},
		{
			name: "lowercase email",
			input: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "JOHN.DOE@EXAMPLE.COM",
			},
			expected: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
		},
		{
			name: "no changes needed",
			input: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			expected: &models.Student{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SanitizeStudent(tt.input)
			if tt.input.FirstName != tt.expected.FirstName {
				t.Errorf("FirstName = %q, want %q", tt.input.FirstName, tt.expected.FirstName)
			}
			if tt.input.LastName != tt.expected.LastName {
				t.Errorf("LastName = %q, want %q", tt.input.LastName, tt.expected.LastName)
			}
			if tt.input.ClassInfo != tt.expected.ClassInfo {
				t.Errorf("ClassInfo = %q, want %q", tt.input.ClassInfo, tt.expected.ClassInfo)
			}
			if tt.input.Email != tt.expected.Email {
				t.Errorf("Email = %q, want %q", tt.input.Email, tt.expected.Email)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)+1 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
