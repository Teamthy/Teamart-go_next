package notifications

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

// Notification represents an in-app realtime notification.
type Notification struct {
	ID        string                 `json:"id"`
	UserID    int64                  `json:"user_id"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	Read      bool                   `json:"read"`
}

// NotificationService publishes in-app notifications.
type NotificationService struct {
	mu            sync.RWMutex
	pubsub        pubsub.PubSub
	notifications map[int64][]*Notification
}

// NewNotificationService creates a new in-app notification service.
func NewNotificationService(pubsubBroker pubsub.PubSub) *NotificationService {
	return &NotificationService{
		pubsub:        pubsubBroker,
		notifications: make(map[int64][]*Notification),
	}
}

// CreateNotification stores and publishes a notification.
func (s *NotificationService) CreateNotification(ctx context.Context, notification *Notification) error {
	if notification == nil {
		return fmt.Errorf("notification cannot be nil")
	}

	notification.CreatedAt = time.Now()
	notification.ID = fmt.Sprintf("notif_%d", time.Now().UnixNano())

	s.mu.Lock()
	s.notifications[notification.UserID] = append(s.notifications[notification.UserID], notification)
	s.mu.Unlock()

	return s.pubsub.Publish(ctx, pubsub.Topic(fmt.Sprintf("notifications:%d", notification.UserID)), notification)
}

// GetNotifications returns the user's current notification feed.
func (s *NotificationService) GetNotifications(userID int64) []*Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]*Notification(nil), s.notifications[userID]...)
}
