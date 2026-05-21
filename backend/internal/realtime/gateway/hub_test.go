package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

func TestHubBroadcastRoutesMessagesAcrossClients(t *testing.T) {
	broker := pubsub.NewInMemoryPubSub()
	hubA := NewHub(broker)
	hubB := NewHub(broker)

	clientA := &Client{ID: "client-a", UserID: 1, Subscriptions: make(map[pubsub.Topic]struct{}), Send: make(chan interface{}, 4)}
	clientB := &Client{ID: "client-b", UserID: 2, Subscriptions: make(map[pubsub.Topic]struct{}), Send: make(chan interface{}, 4)}

	hubA.RegisterClient(clientA)
	hubB.RegisterClient(clientB)

	if err := hubA.SubscribeClientToTopic(context.Background(), clientA.ID, "chat:123"); err != nil {
		t.Fatalf("clientA subscribe failed: %v", err)
	}
	if err := hubB.SubscribeClientToTopic(context.Background(), clientB.ID, "chat:123"); err != nil {
		t.Fatalf("clientB subscribe failed: %v", err)
	}

	if err := hubA.Broadcast(context.Background(), "chat:123", "welcome"); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}

	receivedA := waitForMessage(t, clientA.Send)
	receivedB := waitForMessage(t, clientB.Send)

	if receivedA != "welcome" || receivedB != "welcome" {
		t.Fatalf("expected both clients to receive welcome, got %v and %v", receivedA, receivedB)
	}
}

func TestHubUnsubscribeStopsDelivery(t *testing.T) {
	broker := pubsub.NewInMemoryPubSub()
	hub := NewHub(broker)
	client := &Client{ID: "client-1", UserID: 1, Subscriptions: make(map[pubsub.Topic]struct{}), Send: make(chan interface{}, 4)}

	hub.RegisterClient(client)
	if err := hub.SubscribeClientToTopic(context.Background(), client.ID, "chat:123"); err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}
	if err := hub.UnsubscribeClientFromTopic(client.ID, "chat:123"); err != nil {
		t.Fatalf("unsubscribe failed: %v", err)
	}

	if err := hub.Broadcast(context.Background(), "chat:123", "second message"); err != nil {
		t.Fatalf("broadcast failed: %v", err)
	}

	select {
	case msg := <-client.Send:
		t.Fatalf("expected no message after unsubscribe, got %v", msg)
	case <-time.After(100 * time.Millisecond):
	}
}

func waitForMessage(t *testing.T, ch chan interface{}) interface{} {
	select {
	case msg := <-ch:
		return msg
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for message")
		return nil
	}
}
