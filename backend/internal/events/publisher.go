package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/teamart/commerce-api/pkg/logger"
)

// EventPublisher publishes events to Kafka.
type EventPublisher struct {
	writer *kafka.Writer
	logger *logger.Logger
	topic  string
}

// NewEventPublisher creates a new event publisher.
func NewEventPublisher(brokers []string, topic string, logger *logger.Logger) *EventPublisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &EventPublisher{
		writer: writer,
		logger: logger,
		topic:  topic,
	}
}

// Publish publishes an event to Kafka.
func (ep *EventPublisher) Publish(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.ID == "" {
		event.ID = generateEventID()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	payload, err := json.Marshal(event)
	if err != nil {
		ep.logger.Errorf("failed to marshal event: %v", err)
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.ID),
		Value: payload,
	}

	if err := ep.writer.WriteMessages(ctx, message); err != nil {
		ep.logger.Errorf("failed to publish event: %v", err)
		return err
	}

	ep.logger.Debugf("published event: %s (type: %s)", event.ID, event.Type)
	return nil
}

// PublishAsync publishes an event asynchronously.
func (ep *EventPublisher) PublishAsync(event *Event, callback func(error)) {
	go func() {
		err := ep.Publish(context.Background(), event)
		if callback != nil {
			callback(err)
		}
	}()
}

// PublishBatch publishes multiple events.
func (ep *EventPublisher) PublishBatch(ctx context.Context, events []*Event) error {
	messages := make([]kafka.Message, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}
		if event.ID == "" {
			event.ID = generateEventID()
		}
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}
		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}
		messages = append(messages, kafka.Message{Key: []byte(event.ID), Value: payload})
	}

	return ep.writer.WriteMessages(ctx, messages...)
}

// Close closes the publisher.
func (ep *EventPublisher) Close() error {
	return ep.writer.Close()
}

// generateEventID generates a unique event ID.
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
