package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/pkg/app"
	"github.com/teamart/commerce-api/pkg/logger"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.Logger.Level, cfg.IsDevelopment())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Infof("starting application in %s mode", cfg.Server.Environment)
	log.Debugf("server listening on %s:%d", cfg.Server.Host, cfg.Server.Port)

	// Create application instance
	application := app.New(cfg, log)

	// Create router (placeholder for now)
	router := createRouter(log)

	// Run application
	if err := application.Run(router); err != nil {
		log.Errorf("application error: %v", err)
		os.Exit(1)
	}
}

// createRouter creates and configures the HTTP router
func createRouter(log *logger.Logger) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"teamart-api"}`)
		log.Debug("health check")
	})

	// Ready check endpoint
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready"}`)
		log.Debug("readiness check")
	})

	// API version endpoint
	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"v1","service":"teamart-commerce-api"}`)
		log.Debug("api version endpoint")
	})

	return mux
}
