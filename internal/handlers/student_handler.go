package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/service"
)

// StudentHandler handles HTTP requests for students
type StudentHandler struct {
	service *service.StudentService
}

// NewStudentHandler creates a new student handler
func NewStudentHandler(db *sql.DB) *StudentHandler {
	repo := repository.NewPostgresStudentRepository(db)
	svc := service.NewStudentService(repo)
	return &StudentHandler{service: svc}
}

// CreateStudent handles POST /api/students
func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		logger.Error("Invalid request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(r.Context(), &student); err != nil {
		// Check if it's a validation error or duplicate email
		if strings.Contains(err.Error(), "validation error") {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		if strings.Contains(err.Error(), "email already exists") {
			respondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
		logger.Error("Failed to create student", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create student")
		return
	}

	respondWithJSON(w, http.StatusCreated, student)
}

// SearchStudents handles GET /api/students/search?q={query}
func (h *StudentHandler) SearchStudents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		respondWithJSON(w, http.StatusOK, []models.StudentSearchResult{})
		return
	}

	results, err := h.service.Search(r.Context(), query)
	if err != nil {
		logger.Error("Failed to search students", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to search students")
		return
	}

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
