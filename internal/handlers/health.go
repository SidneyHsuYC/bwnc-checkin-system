package handlers

import (
	"database/sql"
	"net/http"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/db"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
)

// HealthHandler reports service liveness and database connectivity.
type HealthHandler struct {
	DB *sql.DB
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(database *sql.DB) *HealthHandler {
	return &HealthHandler{DB: database}
}

// Check handles GET /health.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	if err := db.HealthCheck(h.DB); err != nil {
		logger.Error("Health check failed", "error", err)
		respondWithJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status":   "unhealthy",
			"database": "disconnected",
			"error":    err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"status":   "healthy",
		"database": "connected",
	})
}
