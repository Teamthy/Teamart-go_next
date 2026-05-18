package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Pool represents the database connection pool
// This is our infrastructure ownership layer - we explicitly manage the connection lifecycle
type Pool struct {
	pool   *pgxpool.Pool
	logger *logger.Logger

	// Metadata for diagnostics
	config *pgxpool.Config
}

// NewPool creates a new database connection pool with explicit lifecycle management
// This is the entry point for all database infrastructure
func NewPool(ctx context.Context, cfg *config.DatabaseConfig, log *logger.Logger) (*Pool, error) {
	log.Debugf("initializing database connection pool: %s", cfg.URL)

	// Parse the database URL into a pgxpool config
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool settings
	// These settings are explicit and observable - we own them
	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MinConnections)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime

	// Set reasonable defaults if not provided
	if poolConfig.MaxConns == 0 {
		poolConfig.MaxConns = 25
	}
	if poolConfig.MinConns == 0 {
		poolConfig.MinConns = 5
	}
	if poolConfig.MaxConnLifetime == 0 {
		poolConfig.MaxConnLifetime = 15 * time.Minute
	}
	if poolConfig.MaxConnIdleTime == 0 {
		poolConfig.MaxConnIdleTime = 5 * time.Minute
	}

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Verify the connection is alive
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Infof("database connection pool initialized successfully")
	log.Debugf("pool config: maxConns=%d minConns=%d maxConnLifetime=%v maxConnIdleTime=%v",
		poolConfig.MaxConns, poolConfig.MinConns, poolConfig.MaxConnLifetime, poolConfig.MaxConnIdleTime)

	return &Pool{
		pool:   pool,
		logger: log,
		config: poolConfig,
	}, nil
}

// Conn returns a connection from the pool for query execution
// This ensures explicit transaction and query ownership
func (p *Pool) Conn(ctx context.Context) (*pgxpool.Conn, error) {
	return p.pool.Acquire(ctx)
}

// Exec executes a command and returns a command tag with execution info
// Explicit single-query execution - no ORM magic
func (p *Pool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

// QueryRow executes a query that returns a single row
// Explicit query - we own the SQL and performance
func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

// Query executes a query that returns multiple rows
// Explicit query - caller manages row iteration
func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

// BeginTx starts an explicit transaction
// Explicit transaction ownership - caller manages commit/rollback
func (p *Pool) BeginTx(ctx context.Context, opts TxOptions) (Tx, error) {
	tx, err := p.pool.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", err)
	}
	return tx, nil
}

// Health returns health status information about the pool
// This is explicit health validation for infrastructure ownership
func (p *Pool) Health(ctx context.Context) (*HealthStatus, error) {
	return GetHealthStatus(ctx, p.pool)
}

// Close closes the connection pool and releases all resources
// Explicit cleanup - part of lifecycle management
func (p *Pool) Close() {
	if p.pool != nil {
		p.logger.Debugf("closing database connection pool")
		p.pool.Close()
		p.logger.Infof("database connection pool closed")
	}
}

// Stats returns connection pool statistics
// Observable diagnostics for infrastructure ownership
func (p *Pool) Stats() *PoolStats {
	if p.pool == nil {
		return &PoolStats{}
	}

	stat := p.pool.Stat()
	return &PoolStats{
		AcquiredConns:     stat.AcquiredConns(),
		IdleConns:         stat.IdleConns(),
		TotalConns:        stat.TotalConns(),
		ConstructingConns: stat.ConstructingConns(),
	}
}

// PoolStats represents connection pool statistics
type PoolStats struct {
	AcquiredConns     int32
	IdleConns         int32
	TotalConns        int32
	ConstructingConns int32
}
