package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/models"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error("Invalid request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger.Info("Received request", "user", user)

	if user.LastName == "" || user.FirstName == "" || user.Phone == "" || user.Email == "" {
		logger.Warn("Missing required fields")
		http.Error(w, "Missing required fields: first_name, last_name, phone, email", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO users (first_name, last_name, phone, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := h.DB.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		logger.Error("Database error", "error", err)
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("User created successfully", "id", user.ID, "email", user.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	logger.Info("Fetching all users")

	query := `
		SELECT id, first_name, last_name, phone, email, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		logger.Error("Database error", "error", err)
		http.Error(w, "Failed to fetch users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Phone, &user.Email, &user.CreatedAt)
		if err != nil {
			logger.Error("Row scan error", "error", err)
			http.Error(w, "Failed to parse user data", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Rows iteration error", "error", err)
		http.Error(w, "Error reading users", http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully fetched users", "count", len(users))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	logger.Info("Fetching user", "id", userID)

	query := `
		SELECT id, first_name, last_name, phone, email, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := h.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		logger.Warn("User not found", "id", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		logger.Error("Database error", "error", err)
		http.Error(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Successfully fetched user", "id", user.ID, "email", user.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger.Info("Checking system health")

	// Check database connection
	err := h.DB.Ping()
	if err != nil {
		logger.Error("Database ping failed", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]any{
			"status":   "unhealthy",
			"database": "disconnected",
			"error":    err.Error(),
		})
		return
	}

	// Get user count
	var userCount int
	err = h.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		logger.Warn("Failed to count users", "error", err)
		userCount = -1
	}

	logger.Info("System healthy", "user_count", userCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"status":     "healthy",
		"database":   "connected",
		"user_count": userCount,
	})
}
