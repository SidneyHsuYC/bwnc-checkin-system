package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/go-chi/chi/v5"
)

// ClassHandler handles HTTP requests for classes
type ClassHandler struct {
	repo repository.ClassRepository
}

// NewClassHandler creates a new class handler
func NewClassHandler(db *sql.DB) *ClassHandler {
	repo := repository.NewPostgresClassRepository(db)
	return &ClassHandler{repo: repo}
}

// CreateClass handles POST /api/classes
func (h *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var class models.Class

	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		logger.Error("Invalid request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.repo.Create(r.Context(), &class); err != nil {
		// Check for unique constraint violation (duplicate class name)
		if strings.Contains(err.Error(), "unique_class_name") || strings.Contains(err.Error(), "duplicate key") {
			logger.Warn("Duplicate class name", "class_name", class.ClassName)
			respondWithError(w, http.StatusConflict, "A class with this name already exists")
			return
		}
		logger.Error("Failed to create class", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create class")
		return
	}

	logger.Info("Class created successfully", "id", class.ID, "name", class.ClassName)
	respondWithJSON(w, http.StatusCreated, class)
}

// ListClasses handles GET /api/classes
func (h *ClassHandler) ListClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.repo.List(r.Context())
	if err != nil {
		logger.Error("Failed to list classes", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list classes")
		return
	}

	respondWithJSON(w, http.StatusOK, classes)
}

// SearchClasses handles GET /api/classes/search?q={query}
func (h *ClassHandler) SearchClasses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		respondWithJSON(w, http.StatusOK, []models.ClassWithLeader{})
		return
	}

	results, err := h.repo.Search(r.Context(), query)
	if err != nil {
		logger.Error("Failed to search classes", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to search classes")
		return
	}

	logger.Info("Class search completed", "query", query, "results", len(results))
	respondWithJSON(w, http.StatusOK, results)
}

// GetClass handles GET /api/classes/{id}
func (h *ClassHandler) GetClass(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid class ID")
		return
	}

	class, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Class not found")
			return
		}
		logger.Error("Failed to get class", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get class")
		return
	}

	respondWithJSON(w, http.StatusOK, class)
}

// UpdateClass handles PUT /api/classes/{id}
func (h *ClassHandler) UpdateClass(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid class ID")
		return
	}

	var class models.Class
	if err := json.NewDecoder(r.Body).Decode(&class); err != nil {
		logger.Error("Invalid request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	class.ID = id

	if err := h.repo.Update(r.Context(), &class); err != nil {
		logger.Error("Failed to update class", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to update class")
		return
	}

	logger.Info("Class updated successfully", "id", class.ID)
	respondWithJSON(w, http.StatusOK, class)
}

// DeleteClass handles DELETE /api/classes/{id}
func (h *ClassHandler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid class ID")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		logger.Error("Failed to delete class", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete class")
		return
	}

	logger.Info("Class deleted successfully", "id", id)
	w.WriteHeader(http.StatusNoContent)
}
