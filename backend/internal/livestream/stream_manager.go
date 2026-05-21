package livestream

import (
	"context"
	"sync"
	"time"
)

// StreamManager tracks active livestream sessions and dispatches lifecycle events.
type StreamManager struct {
	mu       sync.RWMutex
	sessions map[string]*StreamSession
}

// StreamSession represents a live streaming room.
type StreamSession struct {
	StreamID    string    `json:"stream_id"`
	CreatorID   int64     `json:"creator_id"`
	Title       string    `json:"title"`
	StartedAt   time.Time `json:"started_at"`
	State       string    `json:"state"`
	ViewerCount int       `json:"viewer_count"`
}

// NewStreamManager creates a new stream manager.
func NewStreamManager() *StreamManager {
	return &StreamManager{
		sessions: make(map[string]*StreamSession),
	}
}

// RegisterSession registers a new live stream session.
func (m *StreamManager) RegisterSession(ctx context.Context, session *StreamSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.StreamID] = session
}

// GetSession fetches an active session by ID.
func (m *StreamManager) GetSession(streamID string) (*StreamSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	session, ok := m.sessions[streamID]
	return session, ok
}

// EndSession marks a stream as ended.
func (m *StreamManager) EndSession(streamID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, ok := m.sessions[streamID]; ok {
		session.State = "ended"
	}
}
