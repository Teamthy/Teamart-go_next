package handlers

import (
	"net/http"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/internal/livestream"
	commerce "github.com/teamart/commerce-api/internal/livestream/commerce"
	"github.com/teamart/commerce-api/internal/orders"
	"github.com/teamart/commerce-api/internal/products"
	"github.com/teamart/commerce-api/internal/streaming"
	"github.com/teamart/commerce-api/internal/users"
	"github.com/teamart/commerce-api/pkg/logger"
)

// SetupHandlers initializes all service layers and HTTP handlers
// This is a convenience function to set up the complete request handling pipeline
func SetupHandlers(
	mux *http.ServeMux,
	q *queries.Queries,
	db *database.Pool,
	authConfig *auth.AuthConfig,
	livestreamService *livestream.Service,
	streamingService *streaming.Service,
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

	commerceService := commerce.NewService(
		commerce.WithLogger(log),
		commerce.WithProductFetcher(commerce.NewProductServiceAdapter(productService)),
		commerce.WithOrderCreator(commerce.NewOrderServiceAdapter(orderService)),
		commerce.WithPaymentProcessor(commerce.NewSimplePaymentProcessor(log)),
		commerce.WithWalletManager(commerce.NewInMemoryWalletManager(log)),
		commerce.WithEventDispatcher(commerce.NewLoggerEventDispatcher(log)),
	)

	// Create HTTP handlers
	userHandler := NewUserHandler(userService, log)
	productHandler := NewProductHandler(productService, log)
	orderHandler := NewOrderHandler(orderService, log)
	livestreamHandler := livestream.NewHandler(livestreamService, log)
	commerceHandler := commerce.NewHandler(commerceService, log)

	// Register routes
	RegisterUserRoutes(mux, userHandler)
	RegisterProductRoutes(mux, productHandler)
	RegisterOrderRoutes(mux, orderHandler)
	livestream.RegisterLivestreamRoutes(mux, livestreamHandler)
	commerce.RegisterRoutes(mux, commerceHandler)
	streamingHandler := streaming.NewHandler(streamingService, log)
	streaming.RegisterStreamingRoutes(mux, streamingHandler)

	log.Infof("all HTTP handlers registered successfully (auth, user, product, order, livestream, commerce, streaming)")
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
