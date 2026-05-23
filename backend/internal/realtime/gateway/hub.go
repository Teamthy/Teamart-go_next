package gateway

import (
	"sync"

	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

// Client represents a connected realtime client.
type Client struct {
	ID            string
	UserID        int64
	Send          chan []byte
	Subscriptions map[pubsub.Topic]struct{}
}

// topicSubscription tracks the shared broker subscription for a topic.
type topicSubscription struct {
	msgCh <-chan []byte
	unsub func()
}

// Hub manages connected clients, room subscriptions, and distributed message routing.
type Hub struct {
	mu            sync.RWMutex
	broker        pubsub.PubSub
	clients       map[string]*Client
	rooms         map[pubsub.Topic]map[string]*Client
	subscriptions map[pubsub.Topic]*topicSubscription
}

// NewHub creates a new Hub backed by the provided PubSub.
func NewHub(broker pubsub.PubSub) *Hub {
	return &Hub{
		broker:        broker,
		clients:       make(map[string]*Client),
		rooms:         make(map[pubsub.Topic]map[string]*Client),
		subscriptions: make(map[pubsub.Topic]*topicSubscription),
	}
}

// RegisterClient makes the hub aware of a new client.
func (h *Hub) RegisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client.Subscriptions == nil {
		client.Subscriptions = make(map[pubsub.Topic]struct{})
	}
	h.clients[client.ID] = client
}

// UnregisterClient removes a client and clears any topic subscriptions.
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

	delete(h.clients, clientID)
	close(client.Send)
	h.mu.Unlock()
}

// SubscribeClient subscribes an existing client to a topic.
func (h *Hub) SubscribeClient(clientID string, topic pubsub.Topic) error {
	h.mu.Lock()
	client, ok := h.clients[clientID]
	if !ok {
		h.mu.Unlock()
		return nil
	}

	if _, subscribed := client.Subscriptions[topic]; subscribed {
		h.mu.Unlock()
		return nil
	}

	client.Subscriptions[topic] = struct{}{}
	if _, ok := h.rooms[topic]; !ok {
		h.rooms[topic] = make(map[string]*Client)
	}
	h.rooms[topic][client.ID] = client

	if _, ok := h.subscriptions[topic]; !ok {
		msgCh, unsub, err := h.broker.Subscribe(topic)
		if err != nil {
			delete(client.Subscriptions, topic)
			delete(h.rooms[topic], client.ID)
			if len(h.rooms[topic]) == 0 {
				delete(h.rooms, topic)
			}
			h.mu.Unlock()
			return err
		}

		h.subscriptions[topic] = &topicSubscription{msgCh: msgCh, unsub: unsub}
		go h.forwardTopicMessages(topic, msgCh)
	}

	h.mu.Unlock()
	return nil
}

// UnsubscribeClient removes a client's subscription from a topic.
func (h *Hub) UnsubscribeClient(clientID string, topic pubsub.Topic) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.clients[clientID]
	if !ok {
		return
	}

	if _, subscribed := client.Subscriptions[topic]; !subscribed {
		return
	}

	delete(client.Subscriptions, topic)
	h.removeClientFromTopicLocked(client, topic)
}

func (h *Hub) removeClientFromTopicLocked(client *Client, topic pubsub.Topic) {
	if members, ok := h.rooms[topic]; ok {
		delete(members, client.ID)
		if len(members) == 0 {
			delete(h.rooms, topic)
			if sub, ok := h.subscriptions[topic]; ok {
				sub.unsub()
				delete(h.subscriptions, topic)
			}
		}
	}
}

func (h *Hub) forwardTopicMessages(topic pubsub.Topic, msgCh <-chan []byte) {
	for msg := range msgCh {
		h.broadcastToTopic(topic, msg)
	}
}

func (h *Hub) broadcastToTopic(topic pubsub.Topic, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.rooms[topic] {
		payload := copyBytes(msg)
		select {
		case client.Send <- payload:
		default:
		}
	}
}

func copyBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// Publish sends a message into the shared pub/sub fabric.
func (h *Hub) Publish(topic pubsub.Topic, msg []byte) error {
	return h.broker.Publish(topic, msg)
}
