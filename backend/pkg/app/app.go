package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/pkg/logger"
)

// App represents the application instance
type App struct {
	config *config.Config
	logger *logger.Logger
	server *http.Server

	// Lifecycle management
	mu        sync.Mutex
	isRunning bool
	stopChan  chan struct{}
	errorChan chan error

	// Resource management
	cleanup []func(context.Context) error
}

// New creates a new application instance
func New(cfg *config.Config, log *logger.Logger) *App {
	return &App{
		config:    cfg,
		logger:    log,
		stopChan:  make(chan struct{}),
		errorChan: make(chan error, 1),
		cleanup:   make([]func(context.Context) error, 0),
	}
}

// RegisterCleanup registers a cleanup function to be executed on graceful shutdown
func (a *App) RegisterCleanup(fn func(context.Context) error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cleanup = append(a.cleanup, fn)
}

// Run starts the application and blocks until shutdown
func (a *App) Run(handler http.Handler) error {
	a.mu.Lock()
	if a.isRunning {
		a.mu.Unlock()
		return fmt.Errorf("application is already running")
	}
	a.isRunning = true
	a.mu.Unlock()

	// Configure HTTP server
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
	a.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  a.config.Server.Timeout,
		WriteTimeout: a.config.Server.Timeout,
		IdleTimeout:  30 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		a.logger.Infof("starting server on %s", addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.errorChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for shutdown signal
	return a.waitForShutdown()
}

// waitForShutdown blocks until a shutdown signal is received
func (a *App) waitForShutdown() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		a.logger.Infof("received shutdown signal: %v", sig)
		return a.Shutdown(context.Background())
	case err := <-a.errorChan:
		a.logger.Errorf("application error: %v", err)
		a.Shutdown(context.Background())
		return err
	}
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	a.mu.Lock()
	if !a.isRunning {
		a.mu.Unlock()
		return fmt.Errorf("application is not running")
	}
	a.isRunning = false
	a.mu.Unlock()

	a.logger.Info("initiating graceful shutdown...")

	// Create shutdown context with timeout
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	// Shutdown HTTP server
	if a.server != nil {
		a.logger.Info("shutting down HTTP server...")
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Errorf("error shutting down server: %v", err)
		}
	}

	// Run cleanup functions in reverse order
	a.logger.Info("running cleanup handlers...")
	for i := len(a.cleanup) - 1; i >= 0; i-- {
		cleanupName := fmt.Sprintf("cleanup_%d", i)
		cleanupCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

		a.logger.Debugf("executing %s", cleanupName)
		if err := a.cleanup[i](cleanupCtx); err != nil {
			a.logger.Errorf("error in %s: %v", cleanupName, err)
		}
		cancel()
	}

	// Sync logger
	if err := a.logger.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "error syncing logger: %v\n", err)
	}

	a.logger.Info("shutdown complete")
	return nil
}

// IsRunning returns whether the application is currently running
func (a *App) IsRunning() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.isRunning
}

// GetConfig returns the application configuration
func (a *App) GetConfig() *config.Config {
	return a.config
}

// GetLogger returns the application logger
func (a *App) GetLogger() *logger.Logger {
	return a.logger
}

// NotifyError sends an error to the error channel to trigger shutdown
func (a *App) NotifyError(err error) {
	select {
	case a.errorChan <- err:
	default:
		a.logger.Errorf("error channel full, dropping error: %v", err)
	}
}
