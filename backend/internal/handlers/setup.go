package handlers

import (
	"context"
	"net/http"

	"github.com/teamart/commerce-api/internal/admin"
	"github.com/teamart/commerce-api/internal/analytics"
	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/internal/feed"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/internal/merchant"
	"github.com/teamart/commerce-api/internal/middleware"
	"github.com/teamart/commerce-api/internal/moderation"
	"github.com/teamart/commerce-api/internal/orders"
	"github.com/teamart/commerce-api/internal/products"
	rec "github.com/teamart/commerce-api/internal/recommendation"
	"github.com/teamart/commerce-api/internal/staff"
	"github.com/teamart/commerce-api/internal/tenant"
	"github.com/teamart/commerce-api/internal/users"
	"github.com/teamart/commerce-api/pkg/logger"
)

// SetupHandlers initializes all service layers and HTTP handlers
// This is a convenience function to set up the complete request handling pipeline
type Router interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
	Handle(string, http.Handler)
}

func SetupHandlers(
	mux Router,
	q *queries.Queries,
	db *database.Pool,
	authConfig *auth.AuthConfig,
	moderationService *moderation.ModerationService,
	analyticsService *analytics.AnalyticsService,
	log *logger.Logger,
) {
	// ==================== Auth Services ====================
	// Create auth repositories (using PostgreSQL adapters)
	identityRepo := auth.NewIdentityRepositoryPostgres(db, log)
	sessionRepo := auth.NewSessionRepositoryPostgres(db, log)

	// Create auth services
	identityService := auth.NewIdentityService(authConfig, log, identityRepo)
	sessionService := auth.NewSessionService(authConfig, log, sessionRepo)

	// Create auth handlers
	authHandler := NewAuthHandler(identityService, sessionService, log)
	sessionHandler := NewSessionHandler(sessionService, log)

	// Register auth routes
	RegisterAuthRoutes(mux, authHandler)
	RegisterSessionRoutes(mux, sessionHandler)

	// ==================== User Services ====================
	// Create service layers
	userService := users.NewService(q, log)
	productService := products.NewService(q, log)
	orderService := orders.NewService(q, log)

	merchantRepo := merchant.NewPostgresRepository(db, log)
	merchantService := merchant.NewService(merchantRepo, log)
	staffService := staff.NewService(db, log)
	tenantService := tenant.NewService(db, log)

	// Create HTTP handlers
	userHandler := NewUserHandler(userService, log)
	productHandler := NewProductHandler(productService, log)
	orderHandler := NewOrderHandler(orderService, log)
	merchantHandler := NewMerchantHandler(merchantService, staffService, tenantService, log)
	tenantHandler := NewTenantHandler(merchantService, staffService, tenantService, log)

	// Register routes
	RegisterUserRoutes(mux, userHandler)
	RegisterProductRoutes(mux, productHandler)
	RegisterOrderRoutes(mux, orderHandler)
	RegisterMerchantRoutes(mux, merchantHandler)
	RegisterTenantRoutes(mux, tenantHandler)

	moderationHandler := moderation.NewHandler(moderationService, log)
	RegisterModerationRoutes(mux, moderationHandler)

	analyticsHandler := analytics.NewHandler(analyticsService, log)
	analytics.RegisterAnalyticsRoutes(mux, analyticsHandler)

	// ==================== Feed / Recommendation ====================
	// Create a Postgres-backed recommendation repository and a simple
	// in-memory scorer that will load candidates from the repo at startup.
	recRepo := rec.NewPostgresRepository(log)

	// Load candidates from disk/db and create a service. For now we load
	// candidates at setup and create an in-memory service seeded with them.
	candidates, err := recRepo.ListCandidates(context.Background(), db)
	if err != nil {
		log.Errorf("failed to load recommendation candidates: %v", err)
		candidates = nil
	}

	// Default weights (toy example) — tune in production.
	weights := rec.Weights{WatchTime: 1, Purchases: 3, Reactions: 0.5, Follows: 0.2, CategoryAffinity: 2}
	inmemoryRec := rec.NewInMemoryRecommendationService(candidates, weights)
	feedSvc := feed.NewService(inmemoryRec)

	feedHandler := NewFeedHandler(feedSvc, recRepo, db, log)
	RegisterFeedRoutes(mux, feedHandler)

	// ==================== Admin & Operations ====================
	// Prefer Postgres-backed admin repo/service when DB is available
	var adminSvcInterface admin.Service
	repo := admin.NewPostgresAdminRepository(log)
	if db != nil {
		adminSvcInterface = admin.NewPostgresAdminService(repo, db, log)
	} else {
		adminSvcInterface = admin.NewInMemoryService()
	}
	adminHandler := NewAdminHandler(adminSvcInterface, log)

	// Protect admin routes with permission middleware (requires `admin:access` permission)
	rpm := middleware.NewRequirePermissionMiddleware("admin:access", log)
	mux.Handle("GET /admin/dashboard", rpm.Middleware(http.HandlerFunc(adminHandler.Dashboard)))
	mux.Handle("GET /admin/disputes", rpm.Middleware(http.HandlerFunc(adminHandler.ListDisputes)))
	mux.Handle("POST /admin/disputes", rpm.Middleware(http.HandlerFunc(adminHandler.CreateDispute)))
	mux.Handle("POST /admin/payouts/approve", rpm.Middleware(http.HandlerFunc(adminHandler.ApprovePayout)))
	mux.Handle("GET /admin/fraud/alerts", rpm.Middleware(http.HandlerFunc(adminHandler.ListFraudAlerts)))
	mux.Handle("POST /admin/payouts/request", rpm.Middleware(http.HandlerFunc(adminHandler.RequestPayoutApproval)))
	mux.Handle("GET /admin/audit/logs", rpm.Middleware(http.HandlerFunc(adminHandler.ListAuditLogs)))
	mux.Handle("GET /admin/notifications", rpm.Middleware(http.HandlerFunc(adminHandler.ListNotifications)))
	mux.Handle("POST /admin/creators/verify", rpm.Middleware(http.HandlerFunc(adminHandler.VerifyCreator)))
	mux.Handle("POST /admin/support/refund", rpm.Middleware(http.HandlerFunc(adminHandler.Refund)))
	mux.Handle("POST /admin/support/suspend", rpm.Middleware(http.HandlerFunc(adminHandler.SuspendAccount)))

	log.Infof("all HTTP handlers registered successfully (auth, user, product, order, moderation, analytics)")
}

// HealthCheckHandler handles GET /health requests for health checks
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

// RegisterHealthRoutes registers health check routes
func RegisterHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", HealthCheckHandler)
}
