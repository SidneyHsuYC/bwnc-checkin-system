package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/validation"
)

// CheckinHandler handles check-in HTTP requests
type CheckinHandler struct {
	checkinRepo repository.CheckinRepository
	studentRepo repository.StudentRepository
	eventRepo   repository.EventRepository
}

// NewCheckinHandler creates a new check-in handler
func NewCheckinHandler(db *sql.DB) *CheckinHandler {
	return &CheckinHandler{
		checkinRepo: repository.NewPostgresCheckinRepository(db),
		studentRepo: repository.NewPostgresStudentRepository(db),
		eventRepo:   repository.NewPostgresEventRepository(db),
	}
}

// CreateCheckin handles POST /api/checkins
func (h *CheckinHandler) CreateCheckin(w http.ResponseWriter, r *http.Request) {
	var checkin models.Checkin

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&checkin); err != nil {
		logger.Error("Failed to decode check-in request", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate check-in data
	if err := validation.ValidateCheckin(&checkin); err != nil {
		logger.Warn("Check-in validation failed", "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Verify student exists
	_, err := h.studentRepo.GetByID(r.Context(), checkin.StudentID)
	if err != nil {
		logger.Warn("Student not found", "student_id", checkin.StudentID)
		respondWithError(w, http.StatusBadRequest, "Student not found")
		return
	}

	// Verify event exists
	_, err = h.eventRepo.GetByID(r.Context(), checkin.EventID)
	if err != nil {
		logger.Warn("Event not found", "event_id", checkin.EventID)
		respondWithError(w, http.StatusBadRequest, "Event not found")
		return
	}

	// Set check-in timestamp
	checkin.CheckedInAt = time.Now()

	// Create check-in
	if err := h.checkinRepo.Create(r.Context(), &checkin); err != nil {
		if strings.Contains(err.Error(), "already checked in") {
			respondWithError(w, http.StatusConflict, "Student already checked in to this event")
			return
		}
		logger.Error("Failed to create check-in", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create check-in")
		return
	}

	logger.Info("Check-in created", "student_id", checkin.StudentID, "event_id", checkin.EventID)
	respondWithJSON(w, http.StatusCreated, checkin)
}

// ListCheckins handles GET /api/checkins?event_id={id}
func (h *CheckinHandler) ListCheckins(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("event_id")
	if eventIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "event_id query parameter is required")
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid event_id")
		return
	}

	checkins, err := h.checkinRepo.ListByEvent(r.Context(), eventID)
	if err != nil {
		logger.Error("Failed to list check-ins", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list check-ins")
		return
	}

	respondWithJSON(w, http.StatusOK, checkins)
}
