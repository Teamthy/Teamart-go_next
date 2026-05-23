package livestream

import (
	"fmt"
	"sync"
	"time"
)

type streamRecord struct {
	metadata      StreamMetadata
	state         StreamState
	viewers       map[int64]time.Time
	uniqueViewers map[int64]struct{}
	totalJoins    int
	totalLeaves   int
	startedAt     time.Time
	endedAt       time.Time
	liveDuration  time.Duration
	engagement    map[EngagementType]int
	createdAt     time.Time
	updatedAt     time.Time
}

type Service struct {
	mu      sync.RWMutex
	streams map[string]*streamRecord
}

func NewService() *Service {
	return &Service{
		streams: make(map[string]*streamRecord),
	}
}

func (s *Service) CreateStream(metadata StreamMetadata, initialState StreamState) (*StreamInfo, error) {
	if metadata.ID == "" {
		return nil, fmt.Errorf("stream ID is required")
	}
	if metadata.Title == "" {
		return nil, fmt.Errorf("stream title is required")
	}
	if initialState == "" {
		initialState = StreamStateScheduled
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.streams[metadata.ID]; ok {
		return nil, fmt.Errorf("stream %q already exists", metadata.ID)
	}

	now := time.Now()
	record := &streamRecord{
		metadata:      metadata,
		state:         initialState,
		viewers:       make(map[int64]time.Time),
		uniqueViewers: make(map[int64]struct{}),
		engagement:    make(map[EngagementType]int),
		createdAt:     now,
		updatedAt:     now,
	}
	s.streams[metadata.ID] = record

	return s.toInfo(record), nil
}

func (s *Service) GetStream(streamID string) (*StreamInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, ok := s.streams[streamID]
	if !ok {
		return nil, false
	}
	return s.toInfo(record), true
}

func (s *Service) UpdateMetadata(streamID string, metadata StreamMetadata) (*StreamInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}

	record.metadata.Title = metadata.Title
	if metadata.ThumbnailURL != "" {
		record.metadata.ThumbnailURL = metadata.ThumbnailURL
	}
	if metadata.Category != "" {
		record.metadata.Category = metadata.Category
	}
	if len(metadata.Tags) > 0 {
		record.metadata.Tags = metadata.Tags
	}
	if metadata.CreatorName != "" {
		record.metadata.CreatorName = metadata.CreatorName
	}
	if len(metadata.CoHosts) > 0 {
		record.metadata.CoHosts = metadata.CoHosts
	}
	if !metadata.ScheduledAt.IsZero() {
		record.metadata.ScheduledAt = metadata.ScheduledAt
	}
	record.updatedAt = time.Now()

	return s.toInfo(record), nil
}

func (s *Service) TransitionState(streamID string, nextState StreamState) (*StreamInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}

	if !isValidTransition(record.state, nextState) {
		return nil, fmt.Errorf("invalid transition from %q to %q", record.state, nextState)
	}

	now := time.Now()
	switch nextState {
	case StreamStateLive:
		if record.startedAt.IsZero() {
			record.startedAt = now
		}
	case StreamStateEnded:
		if !record.startedAt.IsZero() && record.endedAt.IsZero() {
			record.endedAt = now
			record.liveDuration += now.Sub(record.startedAt)
		}
	case StreamStatePaused:
		if record.state == StreamStateLive {
			record.liveDuration += now.Sub(record.startedAt)
		}
	case StreamStateArchived:
		if record.state != StreamStateEnded {
			return nil, fmt.Errorf("stream must be ended before archiving")
		}
	}

	if nextState == StreamStateLive && record.state == StreamStatePaused {
		record.startedAt = now
	}

	record.state = nextState
	record.updatedAt = now

	return s.toInfo(record), nil
}

func (s *Service) AddViewer(streamID string, userID int64) (*StreamAnalytics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}

	if _, active := record.viewers[userID]; !active {
		record.viewers[userID] = time.Now()
		record.totalJoins++
		record.uniqueViewers[userID] = struct{}{}
	}
	record.updatedAt = time.Now()

	analytics := s.analytics(record)
	return &analytics, nil
}

func (s *Service) RemoveViewer(streamID string, userID int64) (*StreamAnalytics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}

	if _, active := record.viewers[userID]; active {
		delete(record.viewers, userID)
		record.totalLeaves++
	}
	record.updatedAt = time.Now()

	analytics := s.analytics(record)
	return &analytics, nil
}

func (s *Service) TrackEngagement(streamID string, engagementType EngagementType) (*StreamAnalytics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}

	record.engagement[engagementType]++
	record.updatedAt = time.Now()

	analytics := s.analytics(record)
	return &analytics, nil
}

func (s *Service) GetAnalytics(streamID string) (*StreamAnalytics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}

	analytics := s.analytics(record)
	return &analytics, nil
}

func (s *Service) ListStreams() []StreamInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]StreamInfo, 0, len(s.streams))
	for _, record := range s.streams {
		result = append(result, *s.toInfo(record))
	}
	return result
}

func (s *Service) toInfo(record *streamRecord) *StreamInfo {
	analytics := s.analytics(record)
	return &StreamInfo{
		Metadata:  record.metadata,
		State:     record.state,
		Analytics: analytics,
		CreatedAt: record.createdAt,
		UpdatedAt: record.updatedAt,
	}
}

func (s *Service) analytics(record *streamRecord) StreamAnalytics {
	engagementCopy := make(map[EngagementType]int, len(record.engagement))
	for k, v := range record.engagement {
		engagementCopy[k] = v
	}

	liveDuration := record.liveDuration
	if record.state == StreamStateLive && !record.startedAt.IsZero() {
		liveDuration += time.Since(record.startedAt)
	}

	return StreamAnalytics{
		ViewerCount:       len(record.viewers),
		UniqueViewerCount: len(record.uniqueViewers),
		TotalJoinCount:    record.totalJoins,
		TotalLeaveCount:   record.totalLeaves,
		LiveDuration:      liveDuration,
		EngagementCounts:  engagementCopy,
		StartedAt:         record.startedAt,
		EndedAt:           record.endedAt,
	}
}

func isValidTransition(current, next StreamState) bool {
	switch current {
	case StreamStateScheduled:
		return next == StreamStatePreparing || next == StreamStateArchived
	case StreamStatePreparing:
		return next == StreamStateLive || next == StreamStatePaused || next == StreamStateEnded || next == StreamStateArchived
	case StreamStateLive:
		return next == StreamStatePaused || next == StreamStateEnded
	case StreamStatePaused:
		return next == StreamStateLive || next == StreamStateEnded
	case StreamStateEnded:
		return next == StreamStateArchived
	case StreamStateArchived:
		return false
	default:
		return false
	}
}
