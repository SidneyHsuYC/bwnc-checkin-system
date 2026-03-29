package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/validation"
)

// StudentHandler handles HTTP requests for students
type StudentHandler struct {
	repo repository.StudentRepository
}

// NewStudentHandler creates a new student handler
func NewStudentHandler(db *sql.DB) *StudentHandler {
	repo := repository.NewPostgresStudentRepository(db)
	return &StudentHandler{repo: repo}
}

// CreateStudent handles POST /api/students
func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		logger.Error("Invalid request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize and validate
	validation.SanitizeStudent(&student)
	if err := validation.ValidateStudent(&student); err != nil {
		logger.Warn("Student validation failed", "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check email uniqueness
	existing, err := h.repo.GetByEmail(r.Context(), student.Email)
	if err == nil && existing != nil {
		logger.Warn("Duplicate email attempt", "email", student.Email)
		respondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	// Create student
	if err := h.repo.Create(r.Context(), &student); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			respondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
		logger.Error("Failed to create student", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create student")
		return
	}

	logger.Info("Student created successfully", "id", student.ID, "email", student.Email)
	respondWithJSON(w, http.StatusCreated, student)
}

// SearchStudents handles GET /api/students/search?q={query}
func (h *StudentHandler) SearchStudents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		respondWithJSON(w, http.StatusOK, []models.StudentSearchResult{})
		return
	}

	results, err := h.repo.Search(r.Context(), query)
	if err != nil {
		logger.Error("Failed to search students", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to search students")
		return
	}

	logger.Info("Student search completed", "query", query, "results", len(results))
	respondWithJSON(w, http.StatusOK, results)
}

// Helper functions for consistent JSON responses
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal JSON response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
