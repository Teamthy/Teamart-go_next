package pubsub

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Topic identifies a pub/sub channel inside the realtime layer.
type Topic string

// Message is the payload structure passed through pub/sub.
type Message struct {
	Topic     Topic       `json:"topic"`
	Payload   interface{} `json:"payload"`
	CreatedAt time.Time   `json:"created_at"`
}

// SubscriptionHandler handles published messages.
type SubscriptionHandler func(ctx context.Context, message *Message)

// PubSub defines a simple publish/subscribe interface.
type PubSub interface {
	Publish(ctx context.Context, topic Topic, payload interface{}) error
	Subscribe(topic Topic, handler SubscriptionHandler) (string, error)
	Unsubscribe(subscriptionID string) error
}

// InMemoryPubSub is an in-process event bus for realtime services.
type InMemoryPubSub struct {
	mu            sync.RWMutex
	subscriptions map[string]*subscription
	nextID        int64
}

type subscription struct {
	topic   Topic
	handler SubscriptionHandler
}

// NewInMemoryPubSub creates a new in-memory pub/sub broker.
func NewInMemoryPubSub() *InMemoryPubSub {
	return &InMemoryPubSub{
		subscriptions: make(map[string]*subscription),
	}
}

// Publish publishes a message to a topic.
func (ps *InMemoryPubSub) Publish(ctx context.Context, topic Topic, payload interface{}) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	message := &Message{
		Topic:     topic,
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	for _, sub := range ps.subscriptions {
		if sub.topic == topic {
			go sub.handler(ctx, message)
		}
	}

	return nil
}

// Subscribe registers a handler for a topic.
func (ps *InMemoryPubSub) Subscribe(topic Topic, handler SubscriptionHandler) (string, error) {
	if handler == nil {
		return "", fmt.Errorf("handler cannot be nil")
	}

	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.nextID++
	subscriptionID := fmt.Sprintf("sub_%d", ps.nextID)
	ps.subscriptions[subscriptionID] = &subscription{topic: topic, handler: handler}
	return subscriptionID, nil
}

// Unsubscribe removes a subscription.
func (ps *InMemoryPubSub) Unsubscribe(subscriptionID string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.subscriptions[subscriptionID]; !ok {
		return fmt.Errorf("subscription not found")
	}

	delete(ps.subscriptions, subscriptionID)
	return nil
}
