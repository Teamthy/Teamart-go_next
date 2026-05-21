package events

import (
	"context"
	"time"
)

// EventType represents a named event topic used across the platform.
type EventType string

const (
	OrderCreated          EventType = "order.created"
	PaymentCompleted      EventType = "payment.completed"
	InventoryUpdated      EventType = "inventory.updated"
	LivestreamStarted     EventType = "livestream.started"
	LivestreamEnded       EventType = "livestream.ended"
	ChatMessageCreated    EventType = "chat.message"
	ReactionSent          EventType = "reaction.sent"
	NotificationCreated   EventType = "notification.created"
	WalletUpdated         EventType = "wallet.updated"
	CreatorCommissionPaid EventType = "creator.commission_paid"
)

// Event represents a generic platform event that can flow through the event bus.
type Event struct {
	ID            string                 `json:"id"`
	Type          EventType              `json:"type"`
	Timestamp     time.Time              `json:"timestamp"`
	AggregateID   string                 `json:"aggregate_id,omitempty"`
	AggregateType string                 `json:"aggregate_type,omitempty"`
	UserID        *int64                 `json:"user_id,omitempty"`
	Payload       map[string]interface{} `json:"payload,omitempty"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
	Version       int64                  `json:"version,omitempty"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
}

// EventHandler handles events delivered by the event bus.
type EventHandler func(ctx context.Context, event *Event) error

// EventBus is the public interface for publishing and subscribing to events.
type EventBus interface {
	Publish(ctx context.Context, event *Event) error
	Subscribe(eventType EventType, handler EventHandler)
	Start() error
	Stop(ctx context.Context) error
	Close() error
}

// EventHandlerFunc adapts a function to the EventHandler signature.
type EventHandlerFunc func(ctx context.Context, event *Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event *Event) error {
	return f(ctx, event)
}
