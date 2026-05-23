package pubsub

import (
    "testing"
    "time"
)

func TestInMemoryPublishSubscribe(t *testing.T) {
    b := NewInMemoryBroker()
    topic := Topic("room:1")

    ch, unsub, err := b.Subscribe(topic)
    if err != nil {
        t.Fatalf("subscribe error: %v", err)
    }
    defer unsub()

    payload := []byte("hello")
    if err := b.Publish(topic, payload); err != nil {
        t.Fatalf("publish error: %v", err)
    }

    select {
    case got := <-ch:
        if string(got) != string(payload) {
            t.Fatalf("message mismatch: want=%s got=%s", payload, got)
        }
    case <-time.After(1 * time.Second):
        t.Fatalf("timed out waiting for message")
    }
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
    b := NewInMemoryBroker()
    topic := Topic("room:2")

    ch, unsub, _ := b.Subscribe(topic)
    payload := []byte("one")
    unsub()

    // publish after unsubscribe — channel should be closed or not receive
    if err := b.Publish(topic, payload); err != nil {
        t.Fatalf("publish error: %v", err)
    }

    select {
    case _, ok := <-ch:
        if ok {
            t.Fatalf("expected channel closed or no value after unsubscribe")
        }
    case <-time.After(200 * time.Millisecond):
        // if channel wasn't closed immediately, that's acceptable as long as no message
    }
}
