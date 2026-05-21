package gateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

// Client represents a connected realtime client.
type Client struct {
	ID            string
	UserID        int64
	Subscriptions map[pubsub.Topic]struct{}
	Send          chan interface{}
}

// Hub manages active clients, room subscriptions, and distributed message routing.
type Hub struct {
	mu                 sync.RWMutex
	clients            map[string]*Client
	rooms              map[pubsub.Topic]map[string]*Client
	consumerTrackers   map[pubsub.Topic]*topicSubscription
	pubsub             pubsub.PubSub
	nextSubscriptionID int64
}

type topicSubscription struct {
	subscriptionID string
	refCount      int
}

// NewHub creates a new Hub backed by the provided pubsub broker.
func NewHub(broker pubsub.PubSub) *Hub {
	return &Hub{
		clients:          make(map[string]*Client),
		rooms:            make(map[pubsub.Topic]map[string]*Client),
		consumerTrackers: make(map[pubsub.Topic]*topicSubscription),
		pubsub:           broker,
	}
}

// RegisterClient adds a client to the hub.
func (h *Hub) RegisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if c.Subscriptions == nil {
		c.Subscriptions = make(map[pubsub.Topic]struct{})
	}

	h.clients[c.ID] = c
}

// UnregisterClient removes a client and cleans up room subscriptions.
func (h *Hub) UnregisterClient(clientID string) {
	h.mu.Lock()
	client, ok := h.clients[clientID]
	if !ok {
		h.mu.Unlock()
		return
	}

	for topic := range client.Subscriptions {
		h.removeClientFromTopicLocked(client, topic)
	}

	close(client.Send)
	delete(h.clients, clientID)
	h.mu.Unlock()
}

// SubscribeClientToTopic subscribes a client to a topic and ensures the hub listens on that topic.
func (h *Hub) SubscribeClientToTopic(ctx context.Context, clientID string, topic pubsub.Topic) error {
	h.mu.Lock()
	client, ok := h.clients[clientID]
	if !ok {
		h.mu.Unlock()
		return fmt.Errorf("client %q not registered", clientID)
	}

	if _, subscribed := client.Subscriptions[topic]; subscribed {
		h.mu.Unlock()
		return nil
	}

	client.Subscriptions[topic] = struct{}{}

	if _, ok := h.rooms[topic]; !ok {
		h.rooms[topic] = make(map[string]*Client)
	}
		h.rooms[topic][clientID] = client

	tracker, ok := h.consumerTrackers[topic]
	if ok {
		tracker.refCount++
		h.consumerTrackers[topic] = tracker
		h.mu.Unlock()
		return nil
	}

	subscriptionID, err := h.pubsub.Subscribe(topic, h.handleIncomingMessage)
	if err != nil {
		delete(client.Subscriptions, topic)
		delete(h.rooms[topic], clientID)
		if len(h.rooms[topic]) == 0 {
			delete(h.rooms, topic)
		}
		h.mu.Unlock()
		return fmt.Errorf("failed to subscribe to topic %q: %w", topic, err)
	}

	tracker = &topicSubscription{subscriptionID: subscriptionID, refCount: 1}
	h.consumerTrackers[topic] = tracker
	h.mu.Unlock()

	return nil
}

// UnsubscribeClientFromTopic removes a client's subscription to a topic.
func (h *Hub) UnsubscribeClientFromTopic(clientID string, topic pubsub.Topic) error {
	h.mu.Lock()
	client, ok := h.clients[clientID]
	if !ok {
		h.mu.Unlock()
		return fmt.Errorf("client %q not registered", clientID)
	}

	if _, subscribed := client.Subscriptions[topic]; !subscribed {
		h.mu.Unlock()
		return nil
	}

	delete(client.Subscriptions, topic)
	h.removeClientFromTopicLocked(client, topic)
	h.mu.Unlock()
	return nil
}

func (h *Hub) removeClientFromTopicLocked(client *Client, topic pubsub.Topic) {
	if members, ok := h.rooms[topic]; ok {
		delete(members, client.ID)
		if len(members) == 0 {
			delete(h.rooms, topic)
			if tracker, ok := h.consumerTrackers[topic]; ok {
				tracker.refCount--
				if tracker.refCount <= 0 {
					h.pubsub.Unsubscribe(tracker.subscriptionID)
					delete(h.consumerTrackers, topic)
				}
			}
		}
	}
}

func (h *Hub) handleIncomingMessage(_ context.Context, message *pubsub.Message) {
	h.broadcastLocal(message.Topic, message.Payload)
}

func (h *Hub) broadcastLocal(topic pubsub.Topic, payload interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.rooms[topic] {
		select {
		case client.Send <- payload:
		default:
		}
	}
}

// Broadcast publishes a message to the shared pub/sub fabric and routes it to local clients.
func (h *Hub) Broadcast(ctx context.Context, topic pubsub.Topic, payload interface{}) error {
	if err := h.pubsub.Publish(ctx, topic, payload); err != nil {
		return fmt.Errorf("publish failed: %w", err)
	}
	h.broadcastLocal(topic, payload)
	return nil
}
