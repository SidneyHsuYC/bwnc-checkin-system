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

// EventHandler handles HTTP requests for events
type EventHandler struct {
	service *service.EventService
}

// NewEventHandler creates a new event handler
func NewEventHandler(db *sql.DB) *EventHandler {
	repo := repository.NewPostgresEventRepository(db)
	svc := service.NewEventService(repo)
	return &EventHandler{service: svc}
}

// CreateEvent handles POST /api/events
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		logger.Error("Invalid request body", "error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(r.Context(), &event); err != nil {
		// Check if it's a validation error
		if strings.Contains(err.Error(), "validation error") {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		logger.Error("Failed to create event", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	respondWithJSON(w, http.StatusCreated, event)
}

// ListEvents handles GET /api/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.List(r.Context())
	if err != nil {
		logger.Error("Failed to list events", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list events")
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

// ListUpcomingEvents handles GET /api/events/upcoming
func (h *EventHandler) ListUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.ListUpcoming(r.Context())
	if err != nil {
		logger.Error("Failed to list upcoming events", "error", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to list upcoming events")
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}
