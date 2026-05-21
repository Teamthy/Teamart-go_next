package reactions

import (
	"context"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/internal/creator"
	"github.com/teamart/commerce-api/internal/notifications"
	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

// Reaction represents a realtime engagement reaction.
type Reaction struct {
	ReactionID  string    `json:"reaction_id"`
	RoomID      string    `json:"room_id"`
	UserID      int64     `json:"user_id"`
	RecipientID int64     `json:"recipient_id,omitempty"`
	Username    string    `json:"username"`
	Type        string    `json:"type"`
	Emoji       string    `json:"emoji,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	GiftAmount  float64   `json:"gift_amount,omitempty"`
}

// ReactionService delivers realtime reactions to subscribers.
type ReactionService struct {
	pubsub    pubsub.PubSub
	notif     *notifications.Manager
	analytics *creator.CreatorAnalytics
}

// NewReactionService creates a new reaction service.
func NewReactionService(pubsubBroker pubsub.PubSub) *ReactionService {
	return &ReactionService{pubsub: pubsubBroker}
}

// SetNotificationManager attaches a centralized notification manager.
func (s *ReactionService) SetNotificationManager(n *notifications.Manager) {
	s.notif = n
}

// SetCreatorAnalytics attaches creator analytics.
func (s *ReactionService) SetCreatorAnalytics(a *creator.CreatorAnalytics) {
	s.analytics = a
}

// SendReaction publishes a reaction to the room.
func (s *ReactionService) SendReaction(ctx context.Context, reaction *Reaction) error {
	if reaction == nil {
		return fmt.Errorf("reaction cannot be nil")
	}

	reaction.CreatedAt = time.Now()

	sentReaction := *reaction

	if err := s.pubsub.Publish(ctx, pubsub.Topic(reaction.RoomID), &sentReaction); err != nil {
		return err
	}

	if reaction.GiftAmount > 0 {
		targetUser := reaction.RecipientID
		if targetUser == 0 {
			targetUser = reaction.UserID
		}
		if s.notif != nil && targetUser != 0 {
			go func() {
				_ = s.notif.SendUserNotification(ctx, &notifications.NotificationPayload{
					UserID: targetUser,
					Title:  "Gift received",
					Body:   fmt.Sprintf("A gift worth %.2f was sent in the live room.", reaction.GiftAmount),
					Type:   "reaction.gift",
					Data: map[string]interface{}{
						"room_id":     reaction.RoomID,
						"sender_id":   reaction.UserID,
						"gift_amount": reaction.GiftAmount,
						"reaction_id": reaction.ReactionID,
					},
				}, notifications.ChannelRealtime)
			}()
		}
	}

	if reaction.Type == "moderation" && s.notif != nil {
		go func() {
			_ = s.notif.SendUserNotification(ctx, &notifications.NotificationPayload{
				UserID: reaction.UserID,
				Title:  "Chat moderation alert",
				Body:   "One of your reactions was flagged by moderation.",
				Type:   "reaction.moderation",
				Data: map[string]interface{}{
					"room_id":     reaction.RoomID,
					"reaction_id": reaction.ReactionID,
				},
			}, notifications.ChannelRealtime)
		}()
	}

	if s.analytics != nil {
		go func() {
			_ = s.analytics.TrackEvent(ctx, reaction.UserID, "reaction.sent", map[string]interface{}{
				"room_id":       reaction.RoomID,
				"reaction_type": reaction.Type,
				"gift_amount":   reaction.GiftAmount,
				"recipient_id":  reaction.RecipientID,
			})
		}()
	}

	return nil
}
