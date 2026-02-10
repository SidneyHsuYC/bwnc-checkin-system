package service

import (
	"context"
	"fmt"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/validation"
)

// StudentService handles business logic for students
type StudentService struct {
	repo repository.StudentRepository
}

// NewStudentService creates a new student service
func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// Create creates a new student after validation
func (s *StudentService) Create(ctx context.Context, student *models.Student) error {
	// Sanitize input
	validation.SanitizeStudent(student)

	// Validate student data
	if err := validation.ValidateStudent(student); err != nil {
		logger.Warn("Student validation failed", "error", err)
		return fmt.Errorf("validation error: %w", err)
	}

	// Check email uniqueness
	existing, err := s.repo.GetByEmail(ctx, student.Email)
	if err == nil && existing != nil {
		logger.Warn("Duplicate email attempt", "email", student.Email)
		return fmt.Errorf("email already exists: %s", student.Email)
	}

	// Create student
	if err := s.repo.Create(ctx, student); err != nil {
		logger.Error("Failed to create student", "error", err)
		return fmt.Errorf("failed to create student: %w", err)
	}

	logger.Info("Student created successfully", "id", student.ID, "email", student.Email)
	return nil
}

// GetByEmail retrieves a student by email
func (s *StudentService) GetByEmail(ctx context.Context, email string) (*models.Student, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	student, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return student, nil
}

// GetByID retrieves a student by ID
func (s *StudentService) GetByID(ctx context.Context, id int) (*models.Student, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid student ID: %d", id)
	}

	student, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return student, nil
}

// Search searches for students by name
func (s *StudentService) Search(ctx context.Context, query string) ([]*models.StudentSearchResult, error) {
	if query == "" {
		return []*models.StudentSearchResult{}, nil
	}

	results, err := s.repo.Search(ctx, query)
	if err != nil {
		logger.Error("Failed to search students", "query", query, "error", err)
		return nil, fmt.Errorf("failed to search students: %w", err)
	}

	logger.Info("Student search completed", "query", query, "results", len(results))
	return results, nil
}
