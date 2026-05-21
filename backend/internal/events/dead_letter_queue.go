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

// DeadLetterEvent represents an event that failed processing
type DeadLetterEvent struct {
	Event         *Event    `json:"event"`
	Error         string    `json:"error"`
	ErrorCount    int       `json:"error_count"`
	FirstError    time.Time `json:"first_error"`
	LastError     time.Time `json:"last_error"`
	FailureReason string    `json:"failure_reason"`
}

// DeadLetterQueue stores events that couldn't be processed
type DeadLetterQueue struct {
	writer *kafka.Writer
	reader *kafka.Reader
	logger *logger.Logger
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewDeadLetterQueue creates a new dead letter queue
func NewDeadLetterQueue(brokers []string, logger *logger.Logger) *DeadLetterQueue {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "events-dlq",
		Balancer: &kafka.LeastBytes{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          "events-dlq",
		GroupID:        "teamart-dlq-processor",
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
		MaxBytes:       1e6,
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &DeadLetterQueue{
		writer: writer,
		reader: reader,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Enqueue adds a failed event to the DLQ
func (dlq *DeadLetterQueue) Enqueue(ctx context.Context, event *Event, err error) error {
	dlqEvent := &DeadLetterEvent{
		Event:         event,
		Error:         err.Error(),
		FailureReason: "handler_failure",
		FirstError:    time.Now(),
		LastError:     time.Now(),
		ErrorCount:    1,
	}

	data, marshalErr := json.Marshal(dlqEvent)
	if marshalErr != nil {
		return fmt.Errorf("failed to marshal DLQ event: %w", marshalErr)
	}

	message := kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
	}

	if writeErr := dlq.writer.WriteMessages(ctx, message); writeErr != nil {
		return fmt.Errorf("failed to enqueue to DLQ: %w", writeErr)
	}

	dlq.logger.Warnf("event %s moved to DLQ: %v", event.ID, err)
	return nil
}

// Start begins monitoring the DLQ
func (dlq *DeadLetterQueue) Start() error {
	dlq.wg.Add(1)
	go func() {
		defer dlq.wg.Done()
		dlq.monitorDLQ()
	}()

	dlq.logger.Info("dead letter queue monitor started")
	return nil
}

// Stop gracefully stops the DLQ monitor
func (dlq *DeadLetterQueue) Stop(ctx context.Context) error {
	dlq.cancel()

	done := make(chan struct{})
	go func() {
		dlq.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return fmt.Errorf("DLQ stop timeout")
	}

	dlq.writer.Close()
	return dlq.reader.Close()
}

// monitorDLQ continuously monitors for DLQ events
func (dlq *DeadLetterQueue) monitorDLQ() {
	for {
		select {
		case <-dlq.ctx.Done():
			return
		default:
		}

		msg, err := dlq.reader.FetchMessage(dlq.ctx)
		if err != nil {
			if err == context.Canceled {
				return
			}
			continue
		}

		var dlqEvent DeadLetterEvent
		if err := json.Unmarshal(msg.Value, &dlqEvent); err != nil {
			dlq.logger.Errorf("failed to unmarshal DLQ event: %v", err)
			dlq.reader.CommitMessages(dlq.ctx, msg)
			continue
		}

		// Log the DLQ event
		dlq.logDLQEvent(&dlqEvent)

		// Commit message
		if err := dlq.reader.CommitMessages(dlq.ctx, msg); err != nil {
			dlq.logger.Errorf("failed to commit DLQ message: %v", err)
		}
	}
}

// logDLQEvent logs a DLQ event for monitoring and alerting
func (dlq *DeadLetterQueue) logDLQEvent(dlqEvent *DeadLetterEvent) {
	dlq.logger.Errorf(
		"DLQ Event: ID=%s Type=%s AggregateID=%s Reason=%s Error=%s",
		dlqEvent.Event.ID,
		dlqEvent.Event.Type,
		dlqEvent.Event.AggregateID,
		dlqEvent.FailureReason,
		dlqEvent.Error,
	)

	// TODO: Send alert to monitoring system
	// - Alert on high DLQ event rate
	// - Alert on specific event types failing consistently
	// - Store for analysis dashboard
}

// GetDLQEvents retrieves events from the DLQ for analysis
func (dlq *DeadLetterQueue) GetDLQEvents(ctx context.Context, limit int) ([]*DeadLetterEvent, error) {
	events := make([]*DeadLetterEvent, 0, limit)

	for i := 0; i < limit; i++ {
		msg, err := dlq.reader.FetchMessage(ctx)
		if err != nil {
			break
		}

		var dlqEvent DeadLetterEvent
		if err := json.Unmarshal(msg.Value, &dlqEvent); err != nil {
			dlq.logger.Errorf("failed to unmarshal DLQ event: %v", err)
			continue
		}

		events = append(events, &dlqEvent)

		if err := dlq.reader.CommitMessages(ctx, msg); err != nil {
			dlq.logger.Errorf("failed to commit DLQ message: %v", err)
		}
	}

	return events, nil
}

// ReplayEvent attempts to replay a DLQ event
func (dlq *DeadLetterQueue) ReplayEvent(ctx context.Context, eventID string, event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event for replay: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(eventID),
		Value: eventData,
	}

	writer := &kafka.Writer{
		Addr:     dlq.writer.Addr,
		Topic:    MainTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	if err := writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to replay event: %w", err)
	}

	dlq.logger.Infof("replayed DLQ event: %s", eventID)
	return nil
}
