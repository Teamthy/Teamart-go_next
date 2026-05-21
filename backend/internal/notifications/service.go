package notifications

import (
	"context"
	"fmt"

	events "github.com/teamart/commerce-api/internal/events"
	emailservice "github.com/teamart/commerce-api/internal/notifications/email"
	rtNotifications "github.com/teamart/commerce-api/internal/realtime/notifications"
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

// Manager orchestrates centralized notification dispatch.
type Manager struct {
	emailService    *emailservice.EmailService
	realtimeService *rtNotifications.NotificationService
}

// NewManager creates a centralized notification manager.
func NewManager(emailService *emailservice.EmailService, realtimeService *rtNotifications.NotificationService) *Manager {
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
		notification := &rtNotifications.Notification{
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
func (m *Manager) HandleEvent(ctx context.Context, event *events.Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	userID := int64(0)
	if event.UserID != nil {
		userID = *event.UserID
	}

	payload := &NotificationPayload{
		UserID: userID,
		Title:  "Platform Notification",
		Body:   "You have a new update.",
		Type:   string(event.Type),
		Data:   event.Payload,
	}

	switch event.Type {
	case events.OrderCreated:
		payload.Title = "Order Created"
		payload.Body = "Your order has been received and is being processed."
	case events.PaymentCompleted:
		payload.Title = "Payment Completed"
		payload.Body = "Your payment was processed successfully."
	case events.CreatorCommissionPaid:
		payload.Title = "Commission Paid"
		payload.Body = "Your creator commission has been credited."
	case events.LivestreamStarted:
		payload.Title = "Livestream Started"
		payload.Body = "A creator livestream you follow has started."
	case events.LivestreamEnded:
		payload.Title = "Livestream Ended"
		payload.Body = "A creator livestream you follow has ended."
	default:
		payload.Title = fmt.Sprintf("Event: %s", event.Type)
		payload.Body = "An event was published to your account."
	}

	return m.SendUserNotification(ctx, payload)
}
