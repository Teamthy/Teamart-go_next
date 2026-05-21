//go:build pubsub_kafka
// +build pubsub_kafka

package pubsub

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaPubSub is a Kafka-backed PubSub adapter (build-tagged).
type KafkaPubSub struct {
	brokers       []string
	groupID       string
	mu            sync.Mutex
	subscriptions map[string]*kafkaReaderSubscription
	nextID        int64
}

type kafkaReaderSubscription struct {
	reader *kafka.Reader
	cancel context.CancelFunc
}

// NewKafkaPubSub creates a new Kafka PubSub instance.
func NewKafkaPubSub(brokers []string, groupID string) (*KafkaPubSub, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers cannot be empty")
	}

	return &KafkaPubSub{
		brokers:       brokers,
		groupID:       groupID,
		subscriptions: make(map[string]*kafkaReaderSubscription),
	}, nil
}

// Publish implements PubSub.Publish using Kafka.
func (k *KafkaPubSub) Publish(ctx context.Context, topic Topic, payload interface{}) error {
	msg := fmt.Sprintf("%v", payload)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  k.brokers,
		Topic:    string(topic),
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(time.Now().UTC().Format(time.RFC3339Nano)),
		Value: []byte(msg),
	})
}

// Subscribe begins consuming messages from a Kafka topic.
func (k *KafkaPubSub) Subscribe(topic Topic, handler SubscriptionHandler) (string, error) {
	if handler == nil {
		return "", fmt.Errorf("handler cannot be nil")
	}

	subID := fmt.Sprintf("kafka-%d", atomic.AddInt64(&k.nextID, 1))
	ctx, cancel := context.WithCancel(context.Background())
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     k.brokers,
		GroupID:     k.groupID,
		Topic:       string(topic),
		StartOffset: kafka.LastOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
	})

	k.mu.Lock()
	k.subscriptions[subID] = &kafkaReaderSubscription{reader: reader, cancel: cancel}
	k.mu.Unlock()

	go func() {
		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				return
			}
			handler(ctx, &Message{Topic: topic, Payload: string(msg.Value), CreatedAt: time.Now()})
		}
	}()

	return subID, nil
}

// Unsubscribe closes the Kafka reader for the subscription.
func (k *KafkaPubSub) Unsubscribe(subscriptionID string) error {
	k.mu.Lock()
	sub, ok := k.subscriptions[subscriptionID]
	if !ok {
		k.mu.Unlock()
		return fmt.Errorf("subscription not found")
	}
	delete(k.subscriptions, subscriptionID)
	k.mu.Unlock()

	sub.cancel()
	return sub.reader.Close()
}
