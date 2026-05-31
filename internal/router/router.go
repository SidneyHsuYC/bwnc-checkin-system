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

func NewRouter(healthHandler *handlers.HealthHandler, studentHandler *handlers.StudentHandler, eventHandler *handlers.EventHandler, checkinHandler *handlers.CheckinHandler, classHandler *handlers.ClassHandler) http.Handler {
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
	r.Get("/health", healthHandler.Check)

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Student routes
		r.Post("/students", studentHandler.CreateStudent)
		r.Get("/students/search", studentHandler.SearchStudents)

		// Class routes
		r.Post("/classes", classHandler.CreateClass)
		r.Get("/classes", classHandler.ListClasses)
		r.Get("/classes/search", classHandler.SearchClasses)
		r.Get("/classes/{id}", classHandler.GetClass)
		r.Put("/classes/{id}", classHandler.UpdateClass)
		r.Delete("/classes/{id}", classHandler.DeleteClass)

		// Event routes
		r.Post("/events", eventHandler.CreateEvent)
		r.Get("/events", eventHandler.ListEvents)
		r.Get("/events/recent", eventHandler.ListRecentEvents)

		// Check-in routes
		r.Post("/checkins", checkinHandler.CreateCheckin)
		r.Get("/checkins", checkinHandler.ListCheckins)
	})

	// Static files
	fs := http.FileServer(http.Dir("./web/static"))
	r.Handle("/*", fs)

	logger.Info("Router initialized with all endpoints")
	return r
}
