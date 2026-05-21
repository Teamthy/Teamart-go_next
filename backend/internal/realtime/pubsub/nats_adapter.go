//go:build pubsub_nats
// +build pubsub_nats

package pubsub

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
)

// NatsPubSub is a NATS-backed PubSub adapter (build-tagged).
type NatsPubSub struct {
	conn          *nats.Conn
	mu            sync.Mutex
	subscriptions map[string]*nats.Subscription
	nextID        int64
}

func NewNatsPubSub(url string) (*NatsPubSub, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsPubSub{
		conn:          nc,
		subscriptions: make(map[string]*nats.Subscription),
	}, nil
}

func (n *NatsPubSub) Publish(ctx context.Context, topic Topic, payload interface{}) error {
	msg := fmt.Sprintf("%v", payload)
	return n.conn.Publish(string(topic), []byte(msg))
}

func (n *NatsPubSub) Subscribe(topic Topic, handler SubscriptionHandler) (string, error) {
	if handler == nil {
		return "", fmt.Errorf("handler cannot be nil")
	}

	subID := fmt.Sprintf("nats-%d", atomic.AddInt64(&n.nextID, 1))

	subscription, err := n.conn.Subscribe(string(topic), func(msg *nats.Msg) {
		handler(context.Background(), &Message{Topic: topic, Payload: string(msg.Data), CreatedAt: time.Now()})
	})
	if err != nil {
		return "", err
	}

	n.mu.Lock()
	n.subscriptions[subID] = subscription
	n.mu.Unlock()

	return subID, nil
}

func (n *NatsPubSub) Unsubscribe(subscriptionID string) error {
	n.mu.Lock()
	subscription, ok := n.subscriptions[subscriptionID]
	if !ok {
		n.mu.Unlock()
		return fmt.Errorf("subscription not found")
	}
	delete(n.subscriptions, subscriptionID)
	n.mu.Unlock()

	return subscription.Unsubscribe()
}
