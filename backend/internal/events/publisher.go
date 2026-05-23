package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/teamart/commerce-api/pkg/logger"
)

// EventType represents the type of event
type EventType string

const (
	// Auth events
	EventTypeLoginSuccess    EventType = "auth.login.success"
	EventTypeLoginFailed     EventType = "auth.login.failed"
	EventTypeLogout          EventType = "auth.logout"
	EventTypePasswordChanged EventType = "auth.password.changed"
	EventTypeSessionRevoked  EventType = "auth.session.revoked"
	EventTypeTokenRefreshed  EventType = "auth.token.refreshed"

	// Merchant events
	EventTypeMerchantOnboarded EventType = "merchant.onboarded"
	EventTypeMerchantRejected  EventType = "merchant.rejected"
	EventTypeStoreCreated      EventType = "store.created"
	EventTypeStoreUpdated      EventType = "store.updated"

	// Creator events
	EventTypeCreatorOnboarded EventType = "creator.onboarded"
	EventTypeStreamStarted    EventType = "stream.started"
	EventTypeStreamEnded      EventType = "stream.ended"

	// Analytics events
	EventTypeViewerJoined     EventType = "analytics.viewer.joined"
	EventTypeViewerLeft       EventType = "analytics.viewer.left"
	EventTypeReactionSent     EventType = "analytics.reaction.sent"
	EventTypeGiftSent         EventType = "analytics.gift.sent"
	EventTypeProductPinned    EventType = "analytics.product.pinned"
	EventTypeOrderCreated     EventType = "analytics.order.created"
	EventTypePaymentCompleted EventType = "analytics.payment.completed"
	EventTypeCartStarted      EventType = "analytics.cart.started"
	EventTypeCartAbandoned    EventType = "analytics.cart.abandoned"

	// KYC events
	EventTypeKYCSubmitted EventType = "kyc.submitted"
	EventTypeKYCApproved  EventType = "kyc.approved"
	EventTypeKYCRejected  EventType = "kyc.rejected"

	// Security events
	EventTypeSuspiciousLogin EventType = "security.suspicious_login"
	EventTypeAccountLocked   EventType = "security.account_locked"
	EventTypeAccountUnlocked EventType = "security.account_unlocked"

	// Fraud events
	EventTypeFraudDetected EventType = "fraud.detected"
	EventTypeIPBlocked     EventType = "fraud.ip_blocked"
)

// AuditEvent represents an event for audit logging
type AuditEvent struct {
	EventID   string                 `json:"event_id"`
	EventType EventType              `json:"event_type"`
	UserID    int64                  `json:"user_id,omitempty"`
	Email     string                 `json:"email,omitempty"`
	SessionID string                 `json:"session_id,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Severity  string                 `json:"severity"` // info, warning, critical
	Source    string                 `json:"source"`   // auth, fraud, system, etc.
}

// EventPublisher publishes events to Kafka
type EventPublisher struct {
	writer *kafka.Writer
	logger *logger.Logger
	topic  string
}

// NewEventPublisher creates a new event publisher
func NewEventPublisher(brokers []string, topic string, logger *logger.Logger) *EventPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &EventPublisher{
		writer: writer,
		logger: logger,
		topic:  topic,
	}
}

// PublishEvent publishes an event to Kafka
func (ep *EventPublisher) PublishEvent(ctx context.Context, event *AuditEvent) error {
	// Generate unique event ID
	if event.EventID == "" {
		event.EventID = generateEventID()
	}

	// Set timestamp if not set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Marshal event to JSON
	eventData, err := json.Marshal(event)
	if err != nil {
		ep.logger.Errorf("failed to marshal event: %v", err)
		return err
	}

	// Create Kafka message
	message := kafka.Message{
		Key:   []byte(event.EventID),
		Value: eventData,
	}

	// Write to Kafka
	err = ep.writer.WriteMessages(ctx, message)
	if err != nil {
		ep.logger.Errorf("failed to publish event: %v", err)
		return err
	}

	ep.logger.Debugf("published event: %s (type: %s)", event.EventID, event.EventType)
	return nil
}

// PublishAuthEvent publishes an authentication event
func (ep *EventPublisher) PublishAuthEvent(ctx context.Context, eventType EventType, userID int64, email string, data map[string]interface{}) error {
	event := &AuditEvent{
		EventType: eventType,
		UserID:    userID,
		Email:     email,
		Timestamp: time.Now(),
		Data:      data,
		Source:    "auth",
	}

	// Determine severity
	switch eventType {
	case EventTypeLoginFailed, EventTypeSuspiciousLogin, EventTypeAccountLocked:
		event.Severity = "warning"
	case EventTypeFraudDetected, EventTypeIPBlocked:
		event.Severity = "critical"
	default:
		event.Severity = "info"
	}

	return ep.PublishEvent(ctx, event)
}

// PublishMerchantEvent publishes a merchant event
func (ep *EventPublisher) PublishMerchantEvent(ctx context.Context, eventType EventType, userID int64, data map[string]interface{}) error {
	event := &AuditEvent{
		EventType: eventType,
		UserID:    userID,
		Timestamp: time.Now(),
		Data:      data,
		Source:    "merchant",
		Severity:  "info",
	}

	return ep.PublishEvent(ctx, event)
}

// PublishKYCEvent publishes a KYC event
func (ep *EventPublisher) PublishKYCEvent(ctx context.Context, eventType EventType, userID int64, data map[string]interface{}) error {
	event := &AuditEvent{
		EventType: eventType,
		UserID:    userID,
		Timestamp: time.Now(),
		Data:      data,
		Source:    "kyc",
		Severity:  "info",
	}

	return ep.PublishEvent(ctx, event)
}

// PublishSecurityEvent publishes a security event
func (ep *EventPublisher) PublishSecurityEvent(ctx context.Context, eventType EventType, userID int64, data map[string]interface{}) error {
	event := &AuditEvent{
		EventType: eventType,
		UserID:    userID,
		Timestamp: time.Now(),
		Data:      data,
		Source:    "security",
		Severity:  "critical",
	}

	return ep.PublishEvent(ctx, event)
}

func (ep *EventPublisher) PublishAnalyticsEvent(ctx context.Context, eventType EventType, userID int64, data map[string]interface{}) error {
	event := &AuditEvent{
		EventType: eventType,
		UserID:    userID,
		Timestamp: time.Now(),
		Data:      data,
		Source:    "analytics",
		Severity:  "info",
	}
	return ep.PublishEvent(ctx, event)
}

// Close closes the publisher
func (ep *EventPublisher) Close() error {
	return ep.writer.Close()
}

// ===== Kafka Consumer for Processing Events =====

// EventConsumer consumes events from Kafka
type EventConsumer struct {
	reader   *kafka.Reader
	logger   *logger.Logger
	handlers map[EventType]EventHandler
}

// EventHandler is a function that handles an event
type EventHandler func(ctx context.Context, event *AuditEvent) error

// NewEventConsumer creates a new event consumer
func NewEventConsumer(brokers []string, topic string, groupID string, logger *logger.Logger) *EventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &EventConsumer{
		reader:   reader,
		logger:   logger,
		handlers: make(map[EventType]EventHandler),
	}
}

// RegisterHandler registers a handler for an event type
func (ec *EventConsumer) RegisterHandler(eventType EventType, handler EventHandler) {
	ec.handlers[eventType] = handler
}

// Start starts consuming events
func (ec *EventConsumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read message from Kafka
		msg, err := ec.reader.ReadMessage(ctx)
		if err != nil {
			ec.logger.Errorf("error reading message: %v", err)
			continue
		}

		// Unmarshal event
		var event AuditEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			ec.logger.Errorf("failed to unmarshal event: %v", err)
			continue
		}

		// Find handler
		handler, ok := ec.handlers[event.EventType]
		if !ok {
			ec.logger.Debugf("no handler for event type: %s", event.EventType)
			continue
		}

		// Call handler
		err = handler(ctx, &event)
		if err != nil {
			ec.logger.Errorf("error processing event: %v", err)
		}
	}
}

// Close closes the consumer
func (ec *EventConsumer) Close() error {
	return ec.reader.Close()
}

// ===== Helper Functions =====

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// HandleAnalyticsEvent processes events for analytics
func HandleAnalyticsEvent(ctx context.Context, event *AuditEvent) error {
	// Send to analytics service
	// Increment counters, track user engagement, etc.
	return nil
}

// HandleFraudEvent processes fraud events
func HandleFraudEvent(ctx context.Context, event *AuditEvent) error {
	// Send to fraud detection service
	// Trigger additional checks, block accounts, etc.
	return nil
}

// HandleNotificationEvent processes notification events
func HandleNotificationEvent(ctx context.Context, event *AuditEvent) error {
	// Send notifications to users, admins, etc.
	return nil
}

// HandleComplianceEvent processes compliance events
func HandleComplianceEvent(ctx context.Context, event *AuditEvent) error {
	// Store for compliance/audit purposes
	// Trigger automated compliance checks
	return nil
}
