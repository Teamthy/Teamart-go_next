package gateway

import (
	"testing"
	"time"

	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

func TestHubPublishBroadcastsToSubscribedClient(t *testing.T) {
	b := pubsub.NewInMemoryBroker()
	h := NewHub(b)

	client := &Client{
		ID:            "client-1",
		UserID:        1,
		Send:          make(chan []byte, 2),
		Subscriptions: make(map[pubsub.Topic]struct{}),
	}
	h.RegisterClient(client)
	defer h.UnregisterClient(client.ID)

	topic := pubsub.Topic("room:test")
	if err := h.SubscribeClient(client.ID, topic); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}

	payload := []byte(`{"message":"hello"}`)
	if err := h.Publish(topic, payload); err != nil {
		t.Fatalf("publish error: %v", err)
	}

	select {
	case got := <-client.Send:
		if string(got) != string(payload) {
			t.Fatalf("unexpected payload: want=%s got=%s", payload, got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("timed out waiting for message")
	}
}

func TestHubUnsubscribeStopsDeliveringMessages(t *testing.T) {
	b := pubsub.NewInMemoryBroker()
	h := NewHub(b)

	client := &Client{
		ID:            "client-2",
		UserID:        2,
		Send:          make(chan []byte, 2),
		Subscriptions: make(map[pubsub.Topic]struct{}),
	}
	h.RegisterClient(client)
	defer h.UnregisterClient(client.ID)

	topic := pubsub.Topic("room:leave")
	if err := h.SubscribeClient(client.ID, topic); err != nil {
		t.Fatalf("subscribe error: %v", err)
	}
	h.UnsubscribeClient(client.ID, topic)

	payload := []byte(`{"message":"bye"}`)
	if err := h.Publish(topic, payload); err != nil {
		t.Fatalf("publish error: %v", err)
	}

	select {
	case msg := <-client.Send:
		t.Fatalf("expected no message after unsubscribe, got %s", string(msg))
	case <-time.After(200 * time.Millisecond):
	}
}
