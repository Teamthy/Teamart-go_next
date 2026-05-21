package livestream

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/teamart/commerce-api/internal/creator"
	"github.com/teamart/commerce-api/internal/events"
	"github.com/teamart/commerce-api/internal/notifications"
	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

// Session represents a livestream session.
type Session struct {
	ID        string
	CreatorID int64
	Title     string
	StartedAt time.Time
	EndedAt   *time.Time
	Active    bool
	Viewers   int
}

// Manager manages livestream sessions and publishes lifecycle events.
type Manager struct {
	mu        sync.RWMutex
	sessions  map[string]*Session
	broker    pubsub.PubSub
	notif     *notifications.Manager
	analytics *creator.CreatorAnalytics
}

// NewManager creates a session manager.
func NewManager(broker pubsub.PubSub) *Manager {
	return &Manager{sessions: make(map[string]*Session), broker: broker}
}

// SetNotificationManager attaches a notifications.Manager to the livestream manager.
func (m *Manager) SetNotificationManager(n *notifications.Manager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notif = n
}

// SetCreatorAnalytics attaches a CreatorAnalytics instance for event tracking.
func (m *Manager) SetCreatorAnalytics(a *creator.CreatorAnalytics) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.analytics = a
}

// StartSession creates and starts a new session.
func (m *Manager) StartSession(ctx context.Context, id string, creatorID int64, title string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.sessions[id]; ok {
		return nil, fmt.Errorf("session already exists: %s", id)
	}
	s := &Session{ID: id, CreatorID: creatorID, Title: title, StartedAt: time.Now(), Active: true}
	m.sessions[id] = s
	// publish lifecycle event
	_ = m.broker.Publish(ctx, pubsub.Topic("livestream.started"), s)
	// notify users and track analytics
	ev := &events.Event{
		Type:          events.LivestreamStarted,
		AggregateID:   id,
		AggregateType: "livestream",
		Payload: map[string]interface{}{
			"session_id": id,
			"title":      title,
		},
	}
	if m.notif != nil {
		go func() {
			if err := m.notif.HandleEvent(ctx, ev); err != nil {
				fmt.Printf("notification handle error: %v\n", err)
			}
		}()
	}
	if m.analytics != nil {
		go m.analytics.TrackEvent(ctx, creatorID, "livestream.started", map[string]interface{}{"session_id": id, "title": title})
	}
	return s, nil
}

// EndSession marks a session ended and publishes lifecycle event.
func (m *Manager) EndSession(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[id]
	if !ok {
		return fmt.Errorf("session not found: %s", id)
	}
	now := time.Now()
	s.EndedAt = &now
	s.Active = false
	_ = m.broker.Publish(ctx, pubsub.Topic("livestream.ended"), s)
	ev := &events.Event{
		Type:          events.LivestreamEnded,
		AggregateID:   id,
		AggregateType: "livestream",
		Payload: map[string]interface{}{
			"session_id": id,
		},
	}
	if m.notif != nil {
		go func() {
			if err := m.notif.HandleEvent(ctx, ev); err != nil {
				fmt.Printf("notification handle error: %v\n", err)
			}
		}()
	}
	if m.analytics != nil {
		go m.analytics.TrackEvent(ctx, s.CreatorID, "livestream.ended", map[string]interface{}{"session_id": id})
	}
	return nil
}

// GetSession returns a session by ID.
func (m *Manager) GetSession(id string) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.sessions[id]
	return s, ok
}

// UpdateViewerCount updates the viewer count and publishes an update event.
func (m *Manager) UpdateViewerCount(ctx context.Context, id string, viewers int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[id]
	if !ok {
		return fmt.Errorf("session not found: %s", id)
	}
	s.Viewers = viewers
	_ = m.broker.Publish(ctx, pubsub.Topic("livestream.viewer_count"), map[string]interface{}{"session_id": id, "viewers": viewers})
	// analytics hook for viewer count updates
	if m.analytics != nil {
		go m.analytics.TrackEvent(ctx, s.CreatorID, "livestream.viewer_count", map[string]interface{}{"session_id": id, "viewers": viewers})
	}
	return nil
}
