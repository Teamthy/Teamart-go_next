package pubsub

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInMemoryPubSubPublishSubscribe(t *testing.T) {
	broker := NewInMemoryPubSub()
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)

	subID, err := broker.Subscribe("room:1", func(_ context.Context, message *Message) {
		defer wg.Done()
		if message.Topic != "room:1" {
			t.Errorf("expected topic room:1, got %s", message.Topic)
		}
		if payload, ok := message.Payload.(string); !ok || payload != "hello" {
			t.Errorf("expected payload hello, got %#v", message.Payload)
		}
	})
	if err != nil {
		t.Fatalf("subscribe failed: %v", err)
	}

	if err := broker.Publish(ctx, "room:1", "hello"); err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	waitCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for subscription callback")
	}

	if err := broker.Unsubscribe(subID); err != nil {
		t.Fatalf("unsubscribe failed: %v", err)
	}
}

func TestNewPubSubFromConfigReturnsMemoryByDefault(t *testing.T) {
	broker, err := NewPubSubFromConfig(context.Background(), "invalid")
	if err != nil {
		t.Fatalf("expected default memory broker, got error: %v", err)
	}

	if _, ok := broker.(*InMemoryPubSub); !ok {
		t.Fatalf("expected InMemoryPubSub, got %T", broker)
	}
}
