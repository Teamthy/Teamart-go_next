package chat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/teamart/commerce-api/internal/creator"
	"github.com/teamart/commerce-api/internal/notifications"
	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

// ChatMessage describes a live chat message.
type ChatMessage struct {
	MessageID  string    `json:"message_id"`
	RoomID     string    `json:"room_id"`
	UserID     int64     `json:"user_id"`
	Username   string    `json:"username"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
	GiftAmount float64   `json:"gift_amount,omitempty"`
	Moderated  bool      `json:"moderated,omitempty"`
}

// ChatService manages chat delivery and persistence for live rooms.
type ChatService struct {
	mu        sync.RWMutex
	pubsub    pubsub.PubSub
	history   map[string][]*ChatMessage
	maxStored int
	notif     *notifications.Manager
	analytics *creator.CreatorAnalytics
}

// NewChatService creates a new chat service.
func NewChatService(pubsubBroker pubsub.PubSub) *ChatService {
	return &ChatService{
		pubsub:    pubsubBroker,
		history:   make(map[string][]*ChatMessage),
		maxStored: 100,
	}
}

// SetNotificationManager attaches a notifications manager.
func (s *ChatService) SetNotificationManager(n *notifications.Manager) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notif = n
}

// SetCreatorAnalytics attaches a creator analytics instance.
func (s *ChatService) SetCreatorAnalytics(a *creator.CreatorAnalytics) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.analytics = a
}

// SendMessage publishes a chat message to the room.
func (s *ChatService) SendMessage(ctx context.Context, message *ChatMessage) error {
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	message.CreatedAt = time.Now()

	s.mu.Lock()
	s.history[message.RoomID] = append(s.history[message.RoomID], message)
	if len(s.history[message.RoomID]) > s.maxStored {
		s.history[message.RoomID] = s.history[message.RoomID][len(s.history[message.RoomID])-s.maxStored:]
	}
	notifMgr := s.notif
	analytics := s.analytics
	s.mu.Unlock()

	if err := s.pubsub.Publish(ctx, pubsub.Topic(message.RoomID), message); err != nil {
		return err
	}

	if message.Moderated && notifMgr != nil {
		go func() {
			_ = notifMgr.SendUserNotification(ctx, &notifications.NotificationPayload{
				UserID: message.UserID,
				Title:  "Chat moderation notice",
				Body:   "Your chat message was moderated.",
				Type:   "chat.moderation",
				Data: map[string]interface{}{
					"room_id":   message.RoomID,
					"message":   message.Message,
					"moderated": true,
				},
			}, notifications.ChannelRealtime)
		}()
	}

	if message.GiftAmount > 0 && notifMgr != nil {
		go func() {
			_ = notifMgr.SendUserNotification(ctx, &notifications.NotificationPayload{
				UserID: message.UserID,
				Title:  "Gift received",
				Body:   fmt.Sprintf("You sent a gift of %.2f during the live chat.", message.GiftAmount),
				Type:   "chat.gift",
				Data: map[string]interface{}{
					"room_id":     message.RoomID,
					"gift_amount": message.GiftAmount,
				},
			}, notifications.ChannelRealtime)
		}()
	}

	if analytics != nil {
		go analytics.TrackEvent(ctx, message.UserID, "chat.message", map[string]interface{}{
			"room_id":     message.RoomID,
			"moderated":   message.Moderated,
			"gift_amount": message.GiftAmount,
		})
	}

	return nil
}

func (s *ChatService) GetRecentMessages(roomID string) []*ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]*ChatMessage(nil), s.history[roomID]...)
}
