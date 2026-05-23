package streaming

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/internal/media"
	"github.com/teamart/commerce-api/pkg/logger"
)

type Service struct {
	mu         sync.RWMutex
	sessions   map[string]*StreamingSession
	transcoder media.Transcoder
	cfg        config.StreamingConfig
	logger     *logger.Logger
}

func NewService(cfg config.StreamingConfig, transcoder media.Transcoder, logger *logger.Logger) *Service {
	return &Service{
		sessions:   make(map[string]*StreamingSession),
		transcoder: transcoder,
		cfg:        cfg,
		logger:     logger,
	}
}

func (s *Service) CreateSession(title, requestedID string) (*StreamingSession, error) {
	if title == "" {
		return nil, fmt.Errorf("stream title is required")
	}

	sessionID := strings.TrimSpace(requestedID)
	if sessionID == "" {
		sessionID = s.newSessionID()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[sessionID]; ok {
		return nil, fmt.Errorf("session %q already exists", sessionID)
	}

	streamKey := s.generateStreamKey()
	now := time.Now()
	session := &StreamingSession{
		ID:           sessionID,
		StreamKey:    streamKey,
		Title:        title,
		State:        SessionStateIdle,
		CDNProvider:  s.cfg.CDN.Provider,
		HLSDirectory: filepath.Join(s.cfg.HLSOutputPath, sessionID),
		Profiles:     s.profileNames(),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.sessions[sessionID] = session
	return session, nil
}

func (s *Service) GetSession(sessionID string) (*StreamingSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, false
	}
	return cloneSession(session), true
}

func (s *Service) ListSessions() []StreamingSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]StreamingSession, 0, len(s.sessions))
	for _, session := range s.sessions {
		result = append(result, *cloneSession(session))
	}
	return result
}

func (s *Service) StartIngest(sessionID string) (*StreamingSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session %q not found", sessionID)
	}
	if session.State != SessionStateIdle {
		return nil, fmt.Errorf("stream ingestion can only start from idle state")
	}

	session.IngestURL = fmt.Sprintf("%s/%s", strings.TrimRight(s.cfg.RTMPIngestURL, "/"), session.StreamKey)
	session.State = SessionStateIngesting
	session.UpdatedAt = time.Now()

	return cloneSession(session), nil
}

func (s *Service) StartTranscoding(ctx context.Context, sessionID string) (*StreamingSession, error) {
	s.mu.Lock()
	session, ok := s.sessions[sessionID]
	if !ok {
		s.mu.Unlock()
		return nil, fmt.Errorf("session %q not found", sessionID)
	}
	if session.State != SessionStateIngesting {
		if session.State == SessionStateLive {
			s.mu.Unlock()
			return cloneSession(session), nil
		}
		s.mu.Unlock()
		return nil, fmt.Errorf("transcoding can only begin after ingesting")
	}

	session.HLSDirectory = filepath.Join(s.cfg.HLSOutputPath, session.ID)
	session.PlaybackURL = s.transcoder.GetPlaybackManifestURL(session.ID)
	session.State = SessionStateLive
	session.UpdatedAt = time.Now()
	s.mu.Unlock()

	mediaSession := media.Session{
		ID:              session.ID,
		SourceURL:       session.IngestURL,
		OutputPath:      s.cfg.HLSOutputPath,
		PlaylistName:    s.cfg.MasterPlaylistName,
		SegmentDuration: s.cfg.SegmentDurationSeconds,
		Profiles:        s.cfg.Profiles,
		BaseURL:         s.cfg.HLSBaseURL,
	}

	if err := s.transcoder.StartTranscoding(ctx, mediaSession); err != nil {
		s.mu.Lock()
		session.State = SessionStateEnded
		session.LastError = err.Error()
		session.UpdatedAt = time.Now()
		s.mu.Unlock()
		return nil, err
	}

	s.mu.Lock()
	session.UpdatedAt = time.Now()
	s.mu.Unlock()

	return cloneSession(session), nil
}

func (s *Service) StopStream(sessionID string) (*StreamingSession, error) {
	if err := s.transcoder.StopTranscoding(sessionID); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session %q not found", sessionID)
	}

	session.State = SessionStateEnded
	session.UpdatedAt = time.Now()
	return cloneSession(session), nil
}

func (s *Service) profileNames() []string {
	profiles := make([]string, 0, len(s.cfg.Profiles))
	for _, profile := range s.cfg.Profiles {
		profiles = append(profiles, profile.Name)
	}
	return profiles
}

func (s *Service) newSessionID() string {
	return fmt.Sprintf("stream-%s", s.generateStreamKey())
}

func (s *Service) generateStreamKey() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func cloneSession(src *StreamingSession) *StreamingSession {
	copyProfiles := make([]string, len(src.Profiles))
	copy(copyProfiles, src.Profiles)

	return &StreamingSession{
		ID:           src.ID,
		StreamKey:    src.StreamKey,
		Title:        src.Title,
		State:        src.State,
		IngestURL:    src.IngestURL,
		PlaybackURL:  src.PlaybackURL,
		HLSDirectory: src.HLSDirectory,
		Profiles:     copyProfiles,
		CDNProvider:  src.CDNProvider,
		LastError:    src.LastError,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}
