package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/repository"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/validation"
)

// EventHandler handles HTTP requests for events
type EventHandler struct {
	repo repository.EventRepository
}

// NewEventHandler creates a new event handler
func NewEventHandler(db *sql.DB) *EventHandler {
	repo := repository.NewPostgresEventRepository(db)
	return &EventHandler{repo: repo}
}

// CreateEvent handles POST /api/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		logger.Error("Invalid request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Sanitize and validate
	validation.SanitizeEvent(&event)
	if err := validation.ValidateEvent(&event); err != nil {
		logger.Warn("Event validation failed", "error", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Create event
	if err := h.repo.Create(r.Context(), &event); err != nil {
		logger.Error("Failed to create event", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	logger.Info("Event created successfully", "id", event.ID, "name", event.EventName)
	respondWithJSON(w, http.StatusCreated, event)
}

// ListEvents handles GET /api/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.repo.List(r.Context())
	if err != nil {
		logger.Error("Failed to list events", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list events")
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

// ListUpcomingEvents handles GET /api/events/upcoming
func (h *EventHandler) ListUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.repo.ListUpcoming(r.Context())
	if err != nil {
		logger.Error("Failed to list upcoming events", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list upcoming events")
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}
