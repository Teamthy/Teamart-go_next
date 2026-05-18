package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthStatus represents the current health status of the database connection pool
// This is explicit health validation for infrastructure ownership
type HealthStatus struct {
	Status        string    `json:"status"` // "healthy", "degraded", "unhealthy"
	Timestamp     time.Time `json:"timestamp"`
	ResponseTime  int64     `json:"response_time_ms"` // Ping time in milliseconds
	PoolStats     PoolStats `json:"pool_stats"`
	Error         string    `json:"error,omitempty"` // Error message if unhealthy
	LastCheckTime time.Time `json:"last_check_time"`
}

// GetHealthStatus performs an explicit health check on the database pool
// This validates both connectivity and pool availability
func GetHealthStatus(ctx context.Context, pool *pgxpool.Pool) (*HealthStatus, error) {
	if pool == nil {
		return &HealthStatus{
			Status:        "unhealthy",
			Timestamp:     time.Now(),
			Error:         "pool is nil",
			LastCheckTime: time.Now(),
		}, fmt.Errorf("pool is nil")
	}

	startTime := time.Now()
	status := &HealthStatus{
		Timestamp:     startTime,
		LastCheckTime: startTime,
	}

	// Test connection ping
	err := pool.Ping(ctx)
	responseTime := time.Since(startTime)
	status.ResponseTime = responseTime.Milliseconds()

	if err != nil {
		status.Status = "unhealthy"
		status.Error = err.Error()
		return status, fmt.Errorf("health check failed: %w", err)
	}

	// Get pool statistics
	poolStat := pool.Stat()
	status.PoolStats = PoolStats{
		AcquiredConns:     poolStat.AcquiredConns(),
		IdleConns:         poolStat.IdleConns(),
		TotalConns:        poolStat.TotalConns(),
		ConstructingConns: poolStat.ConstructingConns(),
	}

	// Determine health status based on connection availability
	if poolStat.IdleConns() == 0 && poolStat.AcquiredConns() > 0 {
		// If all connections are acquired, mark as degraded if we have no idle ones
		if poolStat.AcquiredConns() >= poolStat.TotalConns() {
			status.Status = "degraded"
			status.Error = "no idle connections available"
			return status, nil
		}
	}

	status.Status = "healthy"
	return status, nil
}

// HealthChecker provides periodic health checking capabilities
// This ensures continuous infrastructure validation
type HealthChecker struct {
	pool             *pgxpool.Pool
	interval         time.Duration
	ctx              context.Context
	cancel           context.CancelFunc
	lastHealthStatus *HealthStatus
	healthStatusChan chan *HealthStatus
}

// NewHealthChecker creates a new health checker for the database pool
func NewHealthChecker(pool *pgxpool.Pool, interval time.Duration) *HealthChecker {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthChecker{
		pool:             pool,
		interval:         interval,
		ctx:              ctx,
		cancel:           cancel,
		healthStatusChan: make(chan *HealthStatus, 1),
	}
}

// Start begins periodic health checking
func (hc *HealthChecker) Start() {
	go func() {
		ticker := time.NewTicker(hc.interval)
		defer ticker.Stop()

		for {
			select {
			case <-hc.ctx.Done():
				return
			case <-ticker.C:
				// Perform health check with timeout
				checkCtx, cancel := context.WithTimeout(hc.ctx, 5*time.Second)
				status, _ := GetHealthStatus(checkCtx, hc.pool)
				cancel()

				hc.lastHealthStatus = status

				// Send latest status (non-blocking)
				select {
				case hc.healthStatusChan <- status:
				default:
					// If channel is full, skip this update
				}
			}
		}
	}()
}

// Stop stops the health checker
func (hc *HealthChecker) Stop() {
	hc.cancel()
}

// GetLastStatus returns the last known health status
func (hc *HealthChecker) GetLastStatus() *HealthStatus {
	return hc.lastHealthStatus
}

// StatusChan returns a channel that receives health status updates
func (hc *HealthChecker) StatusChan() <-chan *HealthStatus {
	return hc.healthStatusChan
}
