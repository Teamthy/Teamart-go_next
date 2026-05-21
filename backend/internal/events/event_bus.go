package events

import (
	"context"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// EventBusService orchestrates event publishing and consuming
type EventBusService struct {
	publisher       *EventPublisher
	consumer        *EventConsumer
	retryQueue      *RetryQueue
	deadLetterQueue *DeadLetterQueue
	logger          *logger.Logger
}

// NewEventBusService creates a new event bus service
func NewEventBusService(
	brokers []string,
	logger *logger.Logger,
) *EventBusService {
	publisher := NewEventPublisher(brokers, "events", logger)
	retryQueue := NewRetryQueue(brokers, logger)
	dlq := NewDeadLetterQueue(brokers, logger)
	consumer := NewEventConsumer(brokers, "teamart-commerce", "events", logger, retryQueue, dlq)

	return &EventBusService{
		publisher:       publisher,
		consumer:        consumer,
		retryQueue:      retryQueue,
		deadLetterQueue: dlq,
		logger:          logger,
	}
}

// Publish publishes an event to the bus
func (ebs *EventBusService) Publish(ctx context.Context, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.ID == "" {
		event.ID = generateEventID()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	return ebs.publisher.Publish(ctx, event)
}

// Stop gracefully stops the event bus
func (ebs *EventBusService) Stop(ctx context.Context) error {
	if err := ebs.consumer.Stop(ctx); err != nil {
		ebs.logger.Errorf("error stopping consumer: %v", err)
	}

	if err := ebs.retryQueue.Stop(ctx); err != nil {
		ebs.logger.Errorf("error stopping retry queue: %v", err)
	}

	if err := ebs.deadLetterQueue.Stop(ctx); err != nil {
		ebs.logger.Errorf("error stopping dead letter queue: %v", err)
	}

	return nil
}
