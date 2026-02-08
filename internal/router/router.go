package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/handlers"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
)

// Custom logging middleware
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for browser metadata requests
		if r.URL.Path == "/.well-known/appspecific/com.chrome.devtools.json" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Create a custom response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Log favicon requests at DEBUG level
		if r.URL.Path == "/favicon.ico" {
			logger.Request("DEBUG", r.Method, r.URL.Path, ww.Status(), duration.String())
		} else {
			logger.Request("INFO", r.Method, r.URL.Path, ww.Status(), duration.String())
		}
	})
}

func NewRouter(userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(requestLogger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8090"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Favicon route
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/favicon.ico")
	})

	// Health check endpoint
	r.Get("/health", userHandler.HealthCheck)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/user", userHandler.CreateUser)
		r.Get("/users", userHandler.GetUsers)
		r.Get("/user/{id}", userHandler.GetUserByID)
	})

	// Static files
	fs := http.FileServer(http.Dir("./web/static"))
	r.Handle("/*", fs)

	logger.Info("Router initialized with all endpoints")
	return r
}
