package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/teamart/commerce-api/internal/events"
)

type AnalyticsService struct {
	mu sync.RWMutex

	uniqueViewers     map[int64]struct{}
	returningViewers  map[int64]struct{}
	activeViewerJoins map[int64]time.Time
	totalWatchTime    int64
	totalViewerJoins  int64
	totalReactions    int64
	totalGifts        int64
	totalProductPins  int64
	totalOrders       int64
	totalOrderValue   float64
	totalPayments     float64
	cartStarts        int64
	cartAbandonments  int64
}

func NewService() *AnalyticsService {
	return &AnalyticsService{
		uniqueViewers:     make(map[int64]struct{}),
		returningViewers:  make(map[int64]struct{}),
		activeViewerJoins: make(map[int64]time.Time),
	}
}

func (s *AnalyticsService) HandleAuditEvent(ctx context.Context, event *events.AuditEvent) error {
	return s.IngestAuditEvent(ctx, event)
}

func (s *AnalyticsService) IngestEvent(event *EventRecord) error {
	if event == nil {
		return fmt.Errorf("event is required")
	}

	auditEvent := &events.AuditEvent{
		EventType: event.EventType,
		UserID:    event.UserID,
		SessionID: event.SessionID,
		Timestamp: event.Timestamp,
		Data:      event.Data,
		Source:    "analytics",
		Severity:  "info",
	}
	return s.IngestAuditEvent(context.Background(), auditEvent)
}

func (s *AnalyticsService) IngestAuditEvent(ctx context.Context, event *events.AuditEvent) error {
	if event == nil {
		return fmt.Errorf("event is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch event.EventType {
	case events.EventTypeViewerJoined:
		s.totalViewerJoins++
		s.recordViewerJoin(event.UserID, event.Timestamp)
	case events.EventTypeViewerLeft:
		s.recordViewerLeave(event.UserID, event.Timestamp)
	case events.EventTypeReactionSent:
		s.totalReactions++
	case events.EventTypeGiftSent:
		s.totalGifts++
		s.recordPayment(event.Data)
	case events.EventTypeProductPinned:
		s.totalProductPins++
	case events.EventTypeOrderCreated:
		s.totalOrders++
		s.recordOrderValue(event.Data)
	case events.EventTypePaymentCompleted:
		s.recordPayment(event.Data)
	case events.EventTypeCartStarted:
		s.cartStarts++
	case events.EventTypeCartAbandoned:
		s.cartAbandonments++
	default:
		// Unrecognized analytics event is ignored, but not an error.
	}

	return nil
}

func (s *AnalyticsService) recordViewerJoin(userID int64, timestamp time.Time) {
	if userID == 0 {
		return
	}

	if _, seen := s.uniqueViewers[userID]; seen {
		s.returningViewers[userID] = struct{}{}
	} else {
		s.uniqueViewers[userID] = struct{}{}
	}

	if !timestamp.IsZero() {
		s.activeViewerJoins[userID] = timestamp
	}
}

func (s *AnalyticsService) recordViewerLeave(userID int64, timestamp time.Time) {
	joinTime, ok := s.activeViewerJoins[userID]
	if !ok {
		return
	}

	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	delta := timestamp.Sub(joinTime)
	if delta > 0 {
		s.totalWatchTime += int64(delta.Seconds())
	}
	delete(s.activeViewerJoins, userID)
}

func (s *AnalyticsService) recordOrderValue(data map[string]interface{}) {
	amount := parseFloat(data, "amount")
	if amount == 0 {
		amount = parseFloat(data, "total_amount")
	}
	if amount > 0 {
		s.totalOrderValue += amount
	}
}

func (s *AnalyticsService) recordPayment(data map[string]interface{}) {
	amount := parseFloat(data, "amount")
	if amount > 0 {
		s.totalPayments += amount
	}
}

func parseFloat(data map[string]interface{}, key string) float64 {
	if data == nil {
		return 0
	}

	value, ok := data[key]
	if !ok {
		return 0
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		f, err := v.Float64()
		if err == nil {
			return f
		}
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return parsed
		}
	}
	return 0
}

func (s *AnalyticsService) GetMetrics() AnalyticsMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	viewerJoins := float64(maxInt64(1, s.totalViewerJoins))
	returningViewers := float64(len(s.returningViewers))

	creatorMetrics := CreatorMetrics{
		Revenue:          s.totalPayments,
		WatchTimeSeconds: s.totalWatchTime,
		EngagementRate:   float64(s.totalReactions+s.totalGifts+s.totalProductPins) / viewerJoins,
		ConversionRate:   float64(s.totalOrders) / viewerJoins,
	}

	orderCountFloat := float64(s.totalOrders)
	if s.totalOrders == 0 {
		orderCountFloat = 1
	}

	uniqueViewerCount := float64(len(s.uniqueViewers))
	retentionRate := 0.0
	if uniqueViewerCount > 0 {
		retentionRate = returningViewers / uniqueViewerCount
	}

	return AnalyticsMetrics{
		Creator: creatorMetrics,
		Marketplace: MarketplaceMetrics{
			GMV:                 s.totalOrderValue,
			AOV:                 s.totalOrderValue / orderCountFloat,
			RetentionRate:       retentionRate,
			CartAbandonmentRate: float64(s.cartAbandonments) / float64(maxInt64(1, s.cartStarts)),
		},
	}
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
