package events

import (
	"context"
	"encoding/json"
	"time"
)

// EventType represents the type of event
type EventType string

// Event Topics - Core Platform Events
const (
	// Order events
	OrderCreated   EventType = "order.created"
	OrderUpdated   EventType = "order.updated"
	OrderCancelled EventType = "order.cancelled"
	OrderCompleted EventType = "order.completed"
	OrderDispute   EventType = "order.dispute"

	// Payment events
	PaymentInitiated EventType = "payment.initiated"
	PaymentCompleted EventType = "payment.completed"
	PaymentFailed    EventType = "payment.failed"
	RefundInitiated  EventType = "refund.initiated"
	RefundCompleted  EventType = "refund.completed"
	PayoutScheduled  EventType = "payout.scheduled"
	PayoutCompleted  EventType = "payout.completed"

	// Inventory events
	InventoryUpdated EventType = "inventory.updated"
	InventoryLow     EventType = "inventory.low"
	StockOutOfStock  EventType = "stock.out_of_stock"
	StockBackInStock EventType = "stock.back_in_stock"

	// Livestream events
	LivestreamStarted  EventType = "livestream.started"
	LivestreamEnded    EventType = "livestream.ended"
	LivestreamPaused   EventType = "livestream.paused"
	LivestreamResumed  EventType = "livestream.resumed"
	ViewerCountUpdated EventType = "viewer_count.updated"

	// Chat events
	ChatMessageCreated EventType = "chat.message_created"
	ChatMessageDeleted EventType = "chat.message_deleted"
	ChatTypingStarted  EventType = "chat.typing_started"
	ChatTypingEnded    EventType = "chat.typing_ended"

	// Product events
	ProductPinned    EventType = "product.pinned"
	ProductUnpinned  EventType = "product.unpinned"
	FlashSaleStarted EventType = "flash_sale.started"
	FlashSaleEnded   EventType = "flash_sale.ended"

	// Engagement events
	ReactionSent   EventType = "reaction.sent"
	GiftSent       EventType = "gift.sent"
	CommentCreated EventType = "comment.created"

	// Wallet events
	WalletUpdated    EventType = "wallet.updated"
	WalletWithdrawn  EventType = "wallet.withdrawn"
	CommissionEarned EventType = "commission.earned"

	// User events
	UserFollowCreator   EventType = "user.follow_creator"
	UserUnfollowCreator EventType = "user.unfollow_creator"
	UserCreatedStore    EventType = "user.created_store"

	// Notification events
	NotificationCreated EventType = "notification.created"

	// Moderation events
	MessageReported EventType = "message.reported"
	UserReported    EventType = "user.reported"
	ModeratorAction EventType = "moderator.action"
)

// Event represents a platform event
type Event struct {
	ID            string                 `json:"id"`
	Type          EventType              `json:"type"`
	Timestamp     time.Time              `json:"timestamp"`
	AggregateID   string                 `json:"aggregate_id"`   // User, Order, Store, Livestream ID
	AggregateType string                 `json:"aggregate_type"` // user, order, store, livestream
	UserID        *int64                 `json:"user_id,omitempty"`
	Payload       map[string]interface{} `json:"payload"`
	Metadata      map[string]string      `json:"metadata,omitempty"`
	Version       int64                  `json:"version"`
	CorrelationID string                 `json:"correlation_id,omitempty"` // For tracing
}

// EventBus defines the interface for event publishing and consuming
type EventBus interface {
	// Publish publishes an event to the bus
	Publish(ctx context.Context, event *Event) error

	// Subscribe subscribes to events of a specific type
	Subscribe(ctx context.Context, eventType EventType, handler EventHandler) error

	// SubscribeToAll subscribes to all events
	SubscribeToAll(ctx context.Context, handler EventHandler) error

	// Unsubscribe unsubscribes from events
	Unsubscribe(eventType EventType, handlerID string) error

	// Close closes the event bus
	Close() error

	// GetTopics returns all available topics
	GetTopics() []EventType
}

// EventHandler is a function that handles an event
type EventHandler interface {
	// Handle processes the event
	Handle(ctx context.Context, event *Event) error

	// GetID returns the unique ID of this handler
	GetID() string
}

// EventHandlerFunc is a function adapter for EventHandler
type EventHandlerFunc func(ctx context.Context, event *Event) error

func (f EventHandlerFunc) Handle(ctx context.Context, event *Event) error {
	return f(ctx, event)
}

func (f EventHandlerFunc) GetID() string {
	return "anonymous"
}

// RetryPolicy defines retry behavior for failed events
type RetryPolicy struct {
	MaxRetries        int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
}

// DeadLetterEntry represents an event that couldn't be processed
type DeadLetterEntry struct {
	ID         string    `json:"id"`
	Event      *Event    `json:"event"`
	Error      string    `json:"error"`
	CreatedAt  time.Time `json:"created_at"`
	RetryCount int       `json:"retry_count"`
}

// EventPublisher defines methods for publishing events
type EventPublisher interface {
	// Publish publishes an event
	Publish(ctx context.Context, event *Event) error

	// PublishAsync publishes an event asynchronously
	PublishAsync(event *Event, callback func(error))

	// PublishBatch publishes multiple events
	PublishBatch(ctx context.Context, events []*Event) error
}

// EventConsumer defines methods for consuming events
type EventConsumer interface {
	// Subscribe subscribes to events
	Subscribe(ctx context.Context, eventType EventType, handler EventHandler) (string, error)

	// SubscribeWithFilter subscribes with filtering
	SubscribeWithFilter(ctx context.Context, eventType EventType, filter EventFilter, handler EventHandler) (string, error)

	// Unsubscribe unsubscribes from events
	Unsubscribe(subscriptionID string) error

	// GetSubscriptions returns all active subscriptions
	GetSubscriptions() []string
}

// EventFilter allows filtering events
type EventFilter interface {
	// Matches returns true if the event matches the filter
	Matches(event *Event) bool
}

// SimpleFilter is a simple key-value filter
type SimpleFilter struct {
	Key   string
	Value interface{}
}

// Matches checks if event payload contains the filter key-value
func (sf *SimpleFilter) Matches(event *Event) bool {
	if val, ok := event.Payload[sf.Key]; ok {
		return val == sf.Value
	}
	return false
}

// EventPayloads define the structure of specific event payloads

// OrderCreatedPayload
type OrderCreatedPayload struct {
	OrderID   int64     `json:"order_id"`
	BuyerID   int64     `json:"buyer_id"`
	SellerID  int64     `json:"seller_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	ItemCount int       `json:"item_count"`
	CreatedAt time.Time `json:"created_at"`
}

// PaymentCompletedPayload
type PaymentCompletedPayload struct {
	PaymentID   int64     `json:"payment_id"`
	OrderID     int64     `json:"order_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Provider    string    `json:"provider"`
	CompletedAt time.Time `json:"completed_at"`
}

// LivestreamStartedPayload
type LivestreamStartedPayload struct {
	StreamID     string    `json:"stream_id"`
	CreatorID    int64     `json:"creator_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StreamURL    string    `json:"stream_url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	StartedAt    time.Time `json:"started_at"`
}

// ChatMessageCreatedPayload
type ChatMessageCreatedPayload struct {
	MessageID string    `json:"message_id"`
	RoomID    string    `json:"room_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// ReactionSentPayload
type ReactionSentPayload struct {
	ReactionID string    `json:"reaction_id"`
	UserID     int64     `json:"user_id"`
	Username   string    `json:"username"`
	RoomID     string    `json:"room_id"`
	Type       string    `json:"type"` // heart, emoji, gift, etc
	Emoji      string    `json:"emoji,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProductPinnedPayload
type ProductPinnedPayload struct {
	PinID     string    `json:"pin_id"`
	ProductID int64     `json:"product_id"`
	StreamID  string    `json:"stream_id"`
	CreatorID int64     `json:"creator_id"`
	Position  int       `json:"position"` // Order on livestream
	PinnedAt  time.Time `json:"pinned_at"`
}

// InventoryUpdatedPayload
type InventoryUpdatedPayload struct {
	ProductID         int64     `json:"product_id"`
	StoreID           int64     `json:"store_id"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// WalletUpdatedPayload
type WalletUpdatedPayload struct {
	WalletID        int64     `json:"wallet_id"`
	UserID          int64     `json:"user_id"`
	PreviousBalance float64   `json:"previous_balance"`
	NewBalance      float64   `json:"new_balance"`
	TransactionType string    `json:"transaction_type"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Helper function to convert payload to JSON
func PayloadToJSON(payload interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &result)
	return result, err
}

// Helper function to extract payload from event
func ExtractPayload(event *Event, target interface{}) error {
	jsonData, err := json.Marshal(event.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}
