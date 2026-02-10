package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/service"
)

// CheckinHandler handles check-in HTTP requests
type CheckinHandler struct {
	service *service.CheckinService
}

// NewCheckinHandler creates a new check-in handler
func NewCheckinHandler(db *sql.DB) *CheckinHandler {
	return &CheckinHandler{
		service: service.NewCheckinService(db),
	}
}

// CreateCheckin handles POST /api/checkins
func (h *CheckinHandler) CreateCheckin(w http.ResponseWriter, r *http.Request) {
	var checkin models.Checkin
	if err := json.NewDecoder(r.Body).Decode(&checkin); err != nil {
		logger.Error("Failed to decode check-in request", "error", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &checkin); err != nil {
		logger.Error("Failed to create check-in", "error", err)

		// Handle specific error cases
		errMsg := err.Error()
		if strings.Contains(errMsg, "student ID is required") || strings.Contains(errMsg, "event ID is required") {
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(errMsg, "student not found") || strings.Contains(errMsg, "event not found") {
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusBadRequest)
			return
		}
		if strings.Contains(errMsg, "already checked in") {
			http.Error(w, `{"error": "`+errMsg+`"}`, http.StatusConflict)
			return
		}

		http.Error(w, `{"error": "Failed to create check-in"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(checkin)
	logger.Info("Check-in created", "student_id", checkin.StudentID, "event_id", checkin.EventID)
}

// ListCheckins handles GET /api/checkins?event_id={id}
func (h *CheckinHandler) ListCheckins(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.URL.Query().Get("event_id")
	if eventIDStr == "" {
		http.Error(w, `{"error": "event_id query parameter is required"}`, http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid event_id"}`, http.StatusBadRequest)
		return
	}

	checkins, err := h.service.ListByEvent(r.Context(), eventID)
	if err != nil {
		logger.Error("Failed to list check-ins", "error", err)
		http.Error(w, `{"error": "Failed to list check-ins"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkins)
}
