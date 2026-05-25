package notifications

import (
	"context"
	"fmt"

	events "github.com/teamart/commerce-api/internal/events"
	emailservice "github.com/teamart/commerce-api/internal/notifications/email"
)

// NotificationChannel represents a delivery channel for notifications.
type NotificationChannel string

const (
	ChannelEmail    NotificationChannel = "email"
	ChannelRealtime NotificationChannel = "realtime"
)

// NotificationPayload contains notification content and metadata.
type NotificationPayload struct {
	UserID int64
	Email  string
	Title  string
	Body   string
	Type   string
	Data   map[string]interface{}
}

type realtimeNotifier interface {
	CreateNotification(ctx context.Context, notification *Notification) error
}

// Notification represents a realtime notification payload.
type Notification struct {
	UserID  int64
	Title   string
	Body    string
	Type    string
	Payload map[string]interface{}
}

// Manager orchestrates centralized notification dispatch.
type Manager struct {
	emailService    *emailservice.EmailService
	realtimeService realtimeNotifier
}

// NewManager creates a centralized notification manager.
func NewManager(emailService *emailservice.EmailService, realtimeService realtimeNotifier) *Manager {
	return &Manager{
		emailService:    emailService,
		realtimeService: realtimeService,
	}
}

// SendUserNotification dispatches a notification through configured channels.
func (m *Manager) SendUserNotification(ctx context.Context, payload *NotificationPayload, channels ...NotificationChannel) error {
	if payload == nil {
		return fmt.Errorf("notification payload cannot be nil")
	}
	if payload.UserID == 0 && payload.Email == "" {
		return fmt.Errorf("at least one delivery target is required")
	}
	if payload.Title == "" || payload.Body == "" {
		return fmt.Errorf("notification title and body are required")
	}

	useRealtime := len(channels) == 0
	useEmail := len(channels) == 0
	for _, ch := range channels {
		switch ch {
		case ChannelRealtime:
			useRealtime = true
			useEmail = false
		case ChannelEmail:
			useEmail = true
			useRealtime = false
		}
	}

	if useRealtime && m.realtimeService != nil && payload.UserID != 0 {
		notification := &Notification{
			UserID:  payload.UserID,
			Title:   payload.Title,
			Body:    payload.Body,
			Type:    payload.Type,
			Payload: payload.Data,
		}
		if err := m.realtimeService.CreateNotification(ctx, notification); err != nil {
			return fmt.Errorf("realtime notification failed: %w", err)
		}
	}

	if useEmail && m.emailService != nil && payload.Email != "" {
		emailInput := &emailservice.SendEmailInput{
			To:      []string{payload.Email},
			Subject: payload.Title,
			Body:    payload.Body,
		}
		if _, err := m.emailService.SendEmail(ctx, emailInput); err != nil {
			return fmt.Errorf("email notification failed: %w", err)
		}
	}

	return nil
}

// HandleEvent sends notifications for platform events.
func (m *Manager) HandleEvent(ctx context.Context, event *events.AuditEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	payload := &NotificationPayload{
		UserID: event.UserID,
		Title:  "Platform Notification",
		Body:   "You have a new update.",
		Type:   string(event.EventType),
		Data:   event.Data,
	}

	switch event.EventType {
	case events.EventTypeOrderCreated:
		payload.Title = "Order Created"
		payload.Body = "Your order has been received and is being processed."
	case events.EventTypePaymentCompleted:
		payload.Title = "Payment Completed"
		payload.Body = "Your payment was processed successfully."
	case events.EventTypeCreatorOnboarded:
		payload.Title = "Creator Onboarded"
		payload.Body = "A new creator has been onboarded."
	case events.EventTypeStreamStarted:
		payload.Title = "Livestream Started"
		payload.Body = "A creator livestream you follow has started."
	case events.EventTypeStreamEnded:
		payload.Title = "Livestream Ended"
		payload.Body = "A creator livestream you follow has ended."
	default:
		payload.Title = fmt.Sprintf("Event: %s", event.EventType)
		payload.Body = "An event was published to your account."
	}

	return m.SendUserNotification(ctx, payload)
}
