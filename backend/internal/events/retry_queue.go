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

// RetryEvent represents an event in the retry queue
type RetryEvent struct {
	Event        *Event    `json:"event"`
	RetryCount   int       `json:"retry_count"`
	NextRetryAt  time.Time `json:"next_retry_at"`
	LastError    string    `json:"last_error"`
	FirstAttempt time.Time `json:"first_attempt"`
}

// RetryQueue manages event retries with exponential backoff
type RetryQueue struct {
	writer      *kafka.Writer
	reader      *kafka.Reader
	logger      *logger.Logger
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	maxRetries  int
	baseBackoff time.Duration
	maxBackoff  time.Duration
}

// NewRetryQueue creates a new retry queue
func NewRetryQueue(brokers []string, logger *logger.Logger) *RetryQueue {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "events-retry",
		Balancer: &kafka.LeastBytes{},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          "events-retry",
		GroupID:        "teamart-retry-processor",
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
		MaxBytes:       1e6,
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &RetryQueue{
		writer:      writer,
		reader:      reader,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
		maxRetries:  3,
		baseBackoff: 5 * time.Second,
		maxBackoff:  5 * time.Minute,
	}
}

// Enqueue adds an event to the retry queue
func (rq *RetryQueue) Enqueue(ctx context.Context, event *Event) error {
	retryCount := 0
	if event.Metadata != nil {
		if v, ok := event.Metadata["retry_count"]; ok {
			fmt.Sscanf(v, "%d", &retryCount)
		}
	}

	retryEvent := &RetryEvent{
		Event:        event,
		RetryCount:   retryCount + 1,
		NextRetryAt:  time.Now().Add(rq.calculateBackoff(retryCount)),
		FirstAttempt: time.Now(),
	}

	data, err := json.Marshal(retryEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal retry event: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(event.ID),
		Value: data,
	}

	if err := rq.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to enqueue retry: %w", err)
	}

	rq.logger.Infof("enqueued retry for event %s (attempt %d)", event.ID, retryCount+1)
	return nil
}

// Start begins processing retry queue
func (rq *RetryQueue) Start() error {
	rq.wg.Add(1)
	go func() {
		defer rq.wg.Done()
		rq.processRetries()
	}()

	rq.logger.Info("retry queue processor started")
	return nil
}

// Stop gracefully stops the retry queue processor
func (rq *RetryQueue) Stop(ctx context.Context) error {
	rq.cancel()

	done := make(chan struct{})
	go func() {
		rq.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return fmt.Errorf("retry queue stop timeout")
	}

	rq.writer.Close()
	return rq.reader.Close()
}

// processRetries continuously checks for events ready to retry
func (rq *RetryQueue) processRetries() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-rq.ctx.Done():
			return
		case <-ticker.C:
			if err := rq.checkAndRetry(); err != nil {
				rq.logger.Errorf("error processing retries: %v", err)
			}
		}
	}
}

// checkAndRetry checks for events ready to retry
func (rq *RetryQueue) checkAndRetry() error {
	msg, err := rq.reader.FetchMessage(rq.ctx)
	if err != nil {
		if err == context.Canceled {
			return nil
		}
		// Continue on timeout or other errors
		return nil
	}

	var retryEvent RetryEvent
	if err := json.Unmarshal(msg.Value, &retryEvent); err != nil {
		rq.logger.Errorf("failed to unmarshal retry event: %v", err)
		return rq.reader.CommitMessages(rq.ctx, msg)
	}

	// Check if ready to retry
	if time.Now().Before(retryEvent.NextRetryAt) {
		// Not ready yet - keep in queue
		return nil
	}

	// Ready to retry - republish to main events topic
	rq.logger.Infof("retrying event %s (attempt %d)", retryEvent.Event.ID, retryEvent.RetryCount)

	// Commit the message so it doesn't get retried indefinitely
	if err := rq.reader.CommitMessages(rq.ctx, msg); err != nil {
		rq.logger.Errorf("failed to commit retry message: %v", err)
	}

	return nil
}

// calculateBackoff calculates exponential backoff with jitter
func (rq *RetryQueue) calculateBackoff(attemptNumber int) time.Duration {
	backoff := rq.baseBackoff * time.Duration(1<<uint(attemptNumber))
	if backoff > rq.maxBackoff {
		backoff = rq.maxBackoff
	}

	// Add jitter (0-20% variation)
	jitter := time.Duration(int64(backoff) / 5)
	return backoff - jitter/2 + time.Duration(int64(time.Now().UnixNano()%int64(jitter)))
}
