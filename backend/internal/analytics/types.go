package analytics

import (
	"time"

	"github.com/teamart/commerce-api/internal/events"
)

type CreatorMetrics struct {
	Revenue          float64 `json:"revenue"`
	WatchTimeSeconds int64   `json:"watch_time_seconds"`
	EngagementRate   float64 `json:"engagement_rate"`
	ConversionRate   float64 `json:"conversion_rate"`
}

type MarketplaceMetrics struct {
	GMV                 float64 `json:"gmv"`
	AOV                 float64 `json:"aov"`
	RetentionRate       float64 `json:"retention_rate"`
	CartAbandonmentRate float64 `json:"cart_abandonment_rate"`
}

type AnalyticsMetrics struct {
	Creator     CreatorMetrics     `json:"creator_metrics"`
	Marketplace MarketplaceMetrics `json:"marketplace_metrics"`
}

type AnalyticsEventInput struct {
	EventType events.EventType `json:"event_type"`
	UserID    int64            `json:"user_id,omitempty"`
	SessionID string           `json:"session_id,omitempty"`
	Data      map[string]any   `json:"data,omitempty"`
	Timestamp string           `json:"timestamp,omitempty"`
}

type EventRecord struct {
	EventType events.EventType `json:"event_type"`
	UserID    int64            `json:"user_id,omitempty"`
	SessionID string           `json:"session_id,omitempty"`
	Timestamp time.Time        `json:"timestamp,omitempty"`
	Data      map[string]any   `json:"data,omitempty"`
}
