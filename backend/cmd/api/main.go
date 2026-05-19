package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/internal/handlers"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/internal/infra/migrations"
	"github.com/teamart/commerce-api/internal/infra/queries"
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

	// Initialize database connection pool (Step 1: PostgreSQL Connection Layer)
	// This demonstrates explicit infrastructure ownership
	db, err := database.NewPool(context.Background(), &cfg.Database, log)
	if err != nil {
		log.Errorf("failed to initialize database: %v", err)
		os.Exit(1)
	}

	// Register database cleanup on application shutdown
	// This ensures explicit resource management - we own the lifecycle
	application.RegisterCleanup(func(shutdownCtx context.Context) error {
		log.Info("closing database connection pool...")
		db.Close()
		log.Info("database connection pool closed")
		return nil
	})

	// Check database health on startup
	health, err := db.Health(context.Background())
	if err != nil {
		log.Warnf("database health check failed: %v", err)
	} else {
		log.Infof("database health: %s (response_time=%dms)", health.Status, health.ResponseTime)
	}

	// Run database migrations (Step 2: Migration System)
	// This is explicit schema evolution - we own the process
	log.Infof("running database migrations...")
	runner := migrations.NewRunner(db, log, migrations.Migrations)
	if err := runner.Migrate(context.Background()); err != nil {
		log.Errorf("migration failed: %v", err)
		os.Exit(1)
	}

	// Check migration status
	status, err := runner.Status(context.Background())
	if err != nil {
		log.Warnf("failed to get migration status: %v", err)
	} else {
		log.Infof("migrations status: %d applied, %d pending",
			len(status.AppliedMigrations), status.PendingMigrations)
	}

	// Initialize SQLC queries (Step 3: Type-Safe SQL Generation)
	// This provides compile-time safe access to the database
	log.Infof("initializing SQLC queries")
	q := queries.New(db)

	// Create auth configuration
	authConfig := createAuthConfig()

	// Create router with all handlers
	router := createRouter(log, db, runner, q, authConfig)

	// Run application
	if err := application.Run(router); err != nil {
		log.Errorf("application error: %v", err)
		os.Exit(1)
	}
}
}

// createRouter creates and configures the HTTP router
func createRouter(
	log *logger.Logger,
	db *database.Pool,
	migrationRunner migrations.MigrationRunner,
	q *queries.Queries,
	authConfig *auth.AuthConfig,
) http.Handler {
	mux := http.NewServeMux()

	// Register health check endpoints
	handlers.RegisterHealthRoutes(mux)

	// Register all API handlers (auth, users, products, orders)
	// This initializes service layers and sets up all routes
	handlers.SetupHandlers(mux, q, db, authConfig, log)

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

		// Check database readiness
		dbHealth, err := db.Health(r.Context())
		if err != nil || dbHealth.Status != "healthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"not_ready","database":"unavailable"}`)
			log.Warnf("readiness check failed: database unhealthy")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready","database":"healthy"}`)
		log.Debug("readiness check")
	})

	// API version endpoint
	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"v1","service":"teamart-commerce-api"}`)
		log.Debug("api version endpoint")
	})

	// Database diagnostics endpoint (development only)
	mux.HandleFunc("GET /api/v1/diagnostics/db", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := db.Stats()
		health, err := db.Health(r.Context())
		if err != nil {
			health = &database.HealthStatus{
				Status: "error",
				Error:  err.Error(),
			}
		}

		fmt.Fprintf(w, `{
			"status":"%s",
			"response_time_ms":%d,
			"pool_stats":{
				"acquired_conns":%d,
				"idle_conns":%d,
				"total_conns":%d,
				"constructing_conns":%d
			}
		}`, health.Status, health.ResponseTime, stats.AcquiredConns, stats.IdleConns, stats.TotalConns, stats.ConstructingConns)
		log.Debug("database diagnostics")
	})

	// Migration status endpoint
	mux.HandleFunc("GET /api/v1/diagnostics/migrations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		status, err := migrationRunner.Status(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error":"failed to get migration status"}`)
			log.Errorf("failed to get migration status: %v", err)
			return
		}

		currentVersion := "none"
		if status.CurrentVersion != "" {
			currentVersion = status.CurrentVersion
		}

		fmt.Fprintf(w, `{
			"current_version":"%s",
			"applied":%d,
			"pending":%d,
			"total":%d
		}`, currentVersion, len(status.AppliedMigrations), status.PendingMigrations, status.TotalMigrations)
		log.Debug("migration status")
	})

	log.Infof("HTTP router configured with all endpoints")

	return mux
}

// createAuthConfig creates the authentication configuration with sensible defaults
func createAuthConfig() *auth.AuthConfig {
	return &auth.AuthConfig{
		// JWT Configuration (from environment or defaults)
		JWTSecret:           os.Getenv("JWT_SECRET"),
		JWTAccessTokenTTL:   getDurationFromEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
		JWTRefreshTokenTTL:  getDurationFromEnv("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		JWTEmailTokenTTL:    getDurationFromEnv("JWT_EMAIL_TOKEN_TTL", 24*time.Hour),
		JWTPasswordResetTTL: getDurationFromEnv("JWT_PASSWORD_RESET_TTL", 1*time.Hour),

		// OTP Configuration
		OTPLength:      6,
		OTPTTL:         getDurationFromEnv("OTP_TTL", 5*time.Minute),
		OTPMaxAttempts: 3,

		// Session Configuration
		SessionTTL:         getDurationFromEnv("SESSION_TTL", 24*time.Hour),
		SessionIdleTimeout: getDurationFromEnv("SESSION_IDLE_TIMEOUT", 30*time.Minute),

		// Security Configuration
		MaxLoginAttempts:       5,
		LoginAttemptWindow:     getDurationFromEnv("LOGIN_ATTEMPT_WINDOW", 15*time.Minute),
		PasswordMinLength:      8,
		PasswordRequireSpecial: true,
		PasswordRequireNumbers: true,

		// Device Trust
		RequireDeviceVerification: false,
	}
}

// getDurationFromEnv gets a duration from environment variable or returns default
func getDurationFromEnv(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return d
}
