package events

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/teamart/commerce-api/pkg/logger"
)

// EventConsumer consumes events from Kafka and dispatches handlers.
type EventConsumer struct {
	reader           *kafka.Reader
	logger           *logger.Logger
	handlers         map[EventType][]EventHandler
	handlersMutex    sync.RWMutex
	retryQueue       *RetryQueue
	deadLetterQueue  *DeadLetterQueue
	maxRetries       int
	processingCtx    context.Context
	processingCancel context.CancelFunc
	wg               sync.WaitGroup
}

// NewEventConsumer creates a new Kafka-backed event consumer.
func NewEventConsumer(
	brokers []string,
	groupID string,
	topic string,
	logger *logger.Logger,
	retryQueue *RetryQueue,
	dlq *DeadLetterQueue,
) *EventConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
		MaxBytes:       1 << 20,
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &EventConsumer{
		reader:           reader,
		logger:           logger,
		handlers:         make(map[EventType][]EventHandler),
		retryQueue:       retryQueue,
		deadLetterQueue:  dlq,
		maxRetries:       3,
		processingCtx:    ctx,
		processingCancel: cancel,
	}
}

// Subscribe registers a handler for a specific event type.
func (ec *EventConsumer) Subscribe(eventType EventType, handler EventHandler) {
	ec.handlersMutex.Lock()
	defer ec.handlersMutex.Unlock()

	ec.handlers[eventType] = append(ec.handlers[eventType], handler)
	ec.logger.Infof("subscribed to event type: %s", eventType)
}

// Start begins polling the Kafka topic for new events.
func (ec *EventConsumer) Start() error {
	ec.wg.Add(1)
	go func() {
		defer ec.wg.Done()
		ec.consumeMessages()
	}()

	ec.logger.Info("event consumer started")
	return nil
}

// Stop gracefully stops the consumer.
func (ec *EventConsumer) Stop(ctx context.Context) error {
	ec.processingCancel()

	done := make(chan struct{})
	go func() {
		ec.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return fmt.Errorf("consumer stop timeout")
	}

	return ec.reader.Close()
}

func (ec *EventConsumer) consumeMessages() {
	for {
		select {
		case <-ec.processingCtx.Done():
			return
		default:
		}

		msg, err := ec.reader.FetchMessage(ec.processingCtx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			ec.logger.Errorf("failed to fetch message: %v", err)
			continue
		}

		var event Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			ec.logger.Errorf("failed to unmarshal event: %v", err)
			continue
		}

		ec.processEvent(&event)

		if err := ec.reader.CommitMessages(ec.processingCtx, msg); err != nil {
			ec.logger.Errorf("failed to commit message: %v", err)
		}
	}
}

func (ec *EventConsumer) processEvent(event *Event) {
	ec.handlersMutex.RLock()
	handlers, ok := ec.handlers[event.Type]
	ec.handlersMutex.RUnlock()

	if !ok || len(handlers) == 0 {
		ec.logger.Debugf("no handlers for event type: %s", event.Type)
		return
	}

	for _, handler := range handlers {
		ctx, cancel := context.WithTimeout(ec.processingCtx, 30*time.Second)
		err := handler(ctx, event)
		cancel()

		if err != nil {
			ec.logger.Errorf("handler error for event %s: %v", event.ID, err)

			if event.Metadata == nil {
				event.Metadata = make(map[string]string)
			}

			retryCount := 0
			if v, ok := event.Metadata["retry_count"]; ok {
				fmt.Sscanf(v, "%d", &retryCount)
			}

			if retryCount < ec.maxRetries {
				event.Metadata["retry_count"] = fmt.Sprintf("%d", retryCount+1)
				if err := ec.retryQueue.Enqueue(ec.processingCtx, event); err != nil {
					ec.logger.Errorf("failed to enqueue retry: %v", err)
					if err := ec.deadLetterQueue.Enqueue(ec.processingCtx, event, err); err != nil {
						ec.logger.Errorf("failed to enqueue DLQ: %v", err)
					}
				}
			} else {
				if err := ec.deadLetterQueue.Enqueue(ec.processingCtx, event, fmt.Errorf("max retries exceeded")); err != nil {
					ec.logger.Errorf("failed to enqueue DLQ: %v", err)
				}
			}
		}
	}
}
