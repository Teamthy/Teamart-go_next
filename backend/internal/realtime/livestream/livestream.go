package livestream

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/teamart/commerce-api/internal/realtime/pubsub"
)

// LivestreamState represents the lifecycle state of a live stream.
type LivestreamState string

const (
	StatePreparing LivestreamState = "preparing"
	StateLive      LivestreamState = "live"
	StatePaused    LivestreamState = "paused"
	StateEnded     LivestreamState = "ended"
)

// LivestreamSession stores metadata for an active stream.
type LivestreamSession struct {
	StreamID     string            `json:"stream_id"`
	CreatorID    int64             `json:"creator_id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	ThumbnailURL string            `json:"thumbnail_url,omitempty"`
	State        LivestreamState   `json:"state"`
	StartedAt    time.Time         `json:"started_at"`
	EndedAt      *time.Time        `json:"ended_at,omitempty"`
	ViewerCount  int               `json:"viewer_count"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// LivestreamService manages live sessions and realtime state.
type LivestreamService struct {
	mu       sync.RWMutex
	sessions map[string]*LivestreamSession
	pubsub   pubsub.PubSub
}

// NewLivestreamService creates a livestream service.
func NewLivestreamService(pubsubBroker pubsub.PubSub) *LivestreamService {
	return &LivestreamService{
		sessions: make(map[string]*LivestreamSession),
		pubsub:   pubsubBroker,
	}
}

// StartStream creates a new livestream session and publishes a start event.
func (s *LivestreamService) StartStream(ctx context.Context, session *LivestreamSession) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	s.mu.Lock()
	s.sessions[session.StreamID] = session
	s.mu.Unlock()

	return s.pubsub.Publish(ctx, pubsub.Topic("livestream:"+session.StreamID), session)
}

// EndStream marks a livestream finished and notifies subscribers.
func (s *LivestreamService) EndStream(ctx context.Context, streamID string) error {
	s.mu.Lock()
	sess, ok := s.sessions[streamID]
	if !ok {
		s.mu.Unlock()
		return fmt.Errorf("stream session not found")
	}

	endAt := time.Now()
	sess.State = StateEnded
	sess.EndedAt = &endAt
	s.mu.Unlock()

	return s.pubsub.Publish(ctx, pubsub.Topic("livestream:"+streamID), sess)
}

// UpdateViewerCount updates a stream's live viewer count.
func (s *LivestreamService) UpdateViewerCount(ctx context.Context, streamID string, count int) error {
	s.mu.Lock()
	sess, ok := s.sessions[streamID]
	if !ok {
		s.mu.Unlock()
		return fmt.Errorf("stream session not found")
	}

	sess.ViewerCount = count
	s.mu.Unlock()

	return s.pubsub.Publish(ctx, pubsub.Topic("livestream.viewer_count"), map[string]interface{}{
		"stream_id":    streamID,
		"viewer_count": count,
	})
}

// GetSession returns the current session.
func (s *LivestreamService) GetSession(streamID string) (*LivestreamSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[streamID]
	return sess, ok
}
