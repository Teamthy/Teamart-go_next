package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/internal/analytics"
	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/internal/events"
	"github.com/teamart/commerce-api/internal/handlers"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/internal/infra/migrations"
	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/internal/moderation"
	realtimewebsocket "github.com/teamart/commerce-api/internal/realtime/websocket"
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

	// Initialize database connection pool
	db, err := database.NewPool(context.Background(), &cfg.Database, log)
	if err != nil {
		log.Errorf("failed to initialize database: %v", err)
		os.Exit(1)
	}

	application.RegisterCleanup(func(shutdownCtx context.Context) error {
		log.Info("closing database connection pool...")
		db.Close()
		log.Info("database connection pool closed")
		return nil
	})

	health, err := db.Health(context.Background())
	if err != nil {
		log.Warnf("database health check failed: %v", err)
	} else {
		log.Infof("database health: %s (response_time=%dms)", health.Status, health.ResponseTime)
	}

	log.Infof("running database migrations...")
	runner := migrations.NewRunner(db, log, migrations.Migrations)
	if err := runner.Migrate(context.Background()); err != nil {
		log.Errorf("migration failed: %v", err)
		os.Exit(1)
	}

	status, err := runner.Status(context.Background())
	if err != nil {
		log.Warnf("failed to get migration status: %v", err)
	} else {
		log.Infof("migrations status: %d applied, %d pending", len(status.AppliedMigrations), status.PendingMigrations)
	}

	log.Infof("initializing SQLC queries")
	q := queries.New(db)

	authConfig := createAuthConfig()

	moderationService := moderation.NewService(nil)
	analyticsService := analytics.NewService()

	eventConsumer := events.NewEventConsumer(cfg.Kafka.Brokers, "teamart-events", cfg.Kafka.GroupID, log)
	eventConsumer.RegisterHandler(events.EventTypeViewerJoined, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypeViewerLeft, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypeReactionSent, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypeGiftSent, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypeProductPinned, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypeOrderCreated, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypePaymentCompleted, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypeCartStarted, analyticsService.HandleAuditEvent)
	eventConsumer.RegisterHandler(events.EventTypeCartAbandoned, analyticsService.HandleAuditEvent)

	go func() {
		if err := eventConsumer.Start(context.Background()); err != nil && err != context.Canceled {
			log.Errorf("analytics event consumer stopped unexpectedly: %v", err)
		}
	}()

	application.RegisterCleanup(func(shutdownCtx context.Context) error {
		log.Info("shutting down analytics event consumer...")
		return eventConsumer.Close()
	})

	router := createRouter(log, db, runner, q, authConfig, moderationService, analyticsService)

	if err := application.Run(router); err != nil {
		log.Errorf("application error: %v", err)
		os.Exit(1)
	}
}

type Router struct {
	router *mux.Router
}

func NewRouter() *Router {
	return &Router{router: mux.NewRouter()}
}

func (r *Router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	method, path := parseRoutePattern(pattern)
	route := r.router.HandleFunc(path, handler)
	if method != "" {
		route.Methods(method)
	}
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	method, path := parseRoutePattern(pattern)
	route := r.router.Handle(path, handler)
	if method != "" {
		route.Methods(method)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func parseRoutePattern(pattern string) (method, path string) {
	trimmed := strings.TrimSpace(pattern)
	parts := strings.SplitN(trimmed, " ", 2)
	if len(parts) == 2 {
		if isHTTPMethod(parts[0]) {
			return parts[0], strings.TrimSpace(parts[1])
		}
	}
	return "", trimmed
}

func isHTTPMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions, http.MethodHead:
		return true
	default:
		return false
	}
}

func createRouter(
	log *logger.Logger,
	db *database.Pool,
	migrationRunner migrations.MigrationRunner,
	q *queries.Queries,
	authConfig *auth.AuthConfig,
	moderationService *moderation.ModerationService,
	analyticsService *analytics.AnalyticsService,
) http.Handler {
	mux := NewRouter()

	handlers.SetupHandlers(mux, q, db, authConfig, moderationService, analyticsService, log)

	tokenService := auth.NewTokenService(authConfig, log)
	realtimeHub := realtimewebsocket.NewHub()
	realtimeServer := realtimewebsocket.NewServer(realtimeHub, realtimewebsocket.NewTokenAuthenticator(tokenService), moderationService, log)
	mux.Handle("/ws", realtimeServer)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"teamart-api"}`)
		log.Debug("health check")
	})

	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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

	mux.HandleFunc("GET /api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"v1","service":"teamart-commerce-api"}`)
		log.Debug("api version endpoint")
	})

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
