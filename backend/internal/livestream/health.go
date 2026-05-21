package livestream

import (
	"context"
	"fmt"
	"time"
)

// HealthStatus captures basic stream health metrics.
type HealthStatus struct {
	SessionID   string
	Uptime      time.Duration
	BitrateKbps int
	Viewers     int
	Healthy     bool
}

// HealthMonitor provides hooks for checking stream health.
type HealthMonitor struct{}

func (m *HealthMonitor) Check(ctx context.Context, session *Session) *HealthStatus {
	// Simple heuristic: active session with viewers is healthy.
	viewers := session.Viewers
	healthy := session != nil && session.Active && viewers >= 0
	return &HealthStatus{
		SessionID:   session.ID,
		Uptime:      time.Since(session.StartedAt),
		BitrateKbps: 0,
		Viewers:     viewers,
		Healthy:     healthy,
	}
}

// AlertIfUnhealthy is a placeholder for health-based alerts.
func (m *HealthMonitor) AlertIfUnhealthy(ctx context.Context, status *HealthStatus) {
	if status == nil {
		return
	}
	if !status.Healthy {
		fmt.Printf("[health] unhealthy stream %s: viewers=%d\n", status.SessionID, status.Viewers)
	}
}
