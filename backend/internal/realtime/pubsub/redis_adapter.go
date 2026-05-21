//go:build pubsub_redis
// +build pubsub_redis

package pubsub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisPubSub is a Redis-backed PubSub implementation. This file is build-tagged
// and only included when building with `-tags pubsub_redis`.
type RedisPubSub struct {
	client        *redis.Client
	mu            sync.Mutex
	subscriptions map[string]*redis.PubSub
	cancels       map[string]context.CancelFunc
	nextID        int64
}

type redisSubscription struct {
	pubsub *redis.PubSub
	cancel context.CancelFunc
}

// NewRedisPubSub creates a new Redis PubSub instance.
func NewRedisPubSub(opts *redis.Options) (*RedisPubSub, error) {
	client := redis.NewClient(opts)
	return &RedisPubSub{
		client:        client,
		subscriptions: make(map[string]*redis.PubSub),
		cancels:       make(map[string]context.CancelFunc),
	}, nil
}

// Publish implements PubSub.Publish using Redis PUBLISH.
func (r *RedisPubSub) Publish(ctx context.Context, topic Topic, payload interface{}) error {
	msg := fmt.Sprintf("%v", payload)
	return r.client.Publish(ctx, string(topic), msg).Err()
}

// Subscribe creates a Redis subscription and forwards messages to handler.
func (r *RedisPubSub) Subscribe(topic Topic, handler SubscriptionHandler) (string, error) {
	if handler == nil {
		return "", fmt.Errorf("handler cannot be nil")
	}

	r.mu.Lock()
	r.nextID++
	subscriptionID := fmt.Sprintf("redis-%d", r.nextID)
	r.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	pubsub := r.client.Subscribe(ctx, string(topic))

	if _, err := pubsub.Receive(ctx); err != nil {
		cancel()
		return "", fmt.Errorf("unable to subscribe to redis topic %q: %w", topic, err)
	}

	r.mu.Lock()
	r.subscriptions[subscriptionID] = pubsub
	r.cancels[subscriptionID] = cancel
	r.mu.Unlock()

	go func() {
		ch := pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				handler(ctx, &Message{Topic: topic, Payload: msg.Payload, CreatedAt: time.Now()})
			case <-ctx.Done():
				return
			}
		}
	}()

	return subscriptionID, nil
}

func (r *RedisPubSub) Unsubscribe(subscriptionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	pubsub, ok := r.subscriptions[subscriptionID]
	if !ok {
		return fmt.Errorf("subscription not found")
	}

	cancel := r.cancels[subscriptionID]
	cancel()
	delete(r.cancels, subscriptionID)
	delete(r.subscriptions, subscriptionID)

	return pubsub.Close()
}
