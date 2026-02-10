package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/db"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/handlers"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/logger"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/migration"
	"github.com/SidneyHsuYC/bwnc-checkin-system/internal/router"
)

func main() {
	// Initialize logger with lumberjack rotation
	if err := logger.Init(); err != nil {
		log.Fatalf("❌ Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	logger.Info("Starting server...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using system environment variables")
	} else {
		logger.Info("Loaded .env file")
	}

	// Connect to database
	logger.Info("Connecting to database...")
	database := db.NewPostgres()
	defer database.Close()
	logger.Info("Database connection established")

	// Run migrations
	logger.Info("Running database migrations...")
	if err := migration.RunMigrations(database, "migrations"); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Migrations completed successfully")

	// Initialize handlers and router
	userHandler := &handlers.UserHandler{DB: database}
	studentHandler := handlers.NewStudentHandler(database)
	eventHandler := handlers.NewEventHandler(database)
	checkinHandler := handlers.NewCheckinHandler(database)
	r := router.NewRouter(userHandler, studentHandler, eventHandler, checkinHandler)

	// Start server
	port := ":8090"
	logger.Info("Server running", "url", "http://localhost"+port)
	logger.Info("API Endpoints:")
	logger.Info("   GET    /health       - Health check")
	logger.Info("   POST   /api/user     - Create a new user")
	logger.Info("   GET    /api/users    - Get all users")
	logger.Info("   GET    /api/user/:id - Get user by ID")
	logger.Info("   GET    /             - Static files")
	logger.Info("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Setup graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		logger.Info("\nShutting down server...")
		logger.Close()
		os.Exit(0)
	}()

	if err := http.ListenAndServe(port, r); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
