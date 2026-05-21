package presence

import (
	"sync"
	"time"
)

// PresenceService tracks online/offline status, typing state, and viewer counts.
type PresenceService struct {
	mu           sync.RWMutex
	onlineUsers  map[int64]time.Time
	typingStatus map[string]map[int64]bool
	viewerCounts map[string]int
}

// NewPresenceService creates a new presence tracker.
func NewPresenceService() *PresenceService {
	return &PresenceService{
		onlineUsers:  make(map[int64]time.Time),
		typingStatus: make(map[string]map[int64]bool),
		viewerCounts: make(map[string]int),
	}
}

// SetOnline marks a user as online.
func (s *PresenceService) SetOnline(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.onlineUsers[userID] = time.Now()
}

// SetOffline marks a user as offline.
func (s *PresenceService) SetOffline(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.onlineUsers, userID)
}

// IsOnline returns whether the user is currently online.
func (s *PresenceService) IsOnline(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.onlineUsers[userID]
	return ok
}

// SetTyping updates typing state for a user in a room.
func (s *PresenceService) SetTyping(roomID string, userID int64, typing bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.typingStatus[roomID]; !ok {
		s.typingStatus[roomID] = make(map[int64]bool)
	}

	if typing {
		s.typingStatus[roomID][userID] = true
	} else {
		delete(s.typingStatus[roomID], userID)
	}
}

// GetTypingUsers returns the list of users currently typing in a room.
func (s *PresenceService) GetTypingUsers(roomID string) []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]int64, 0)
	for userID := range s.typingStatus[roomID] {
		users = append(users, userID)
	}
	return users
}

// SetViewerCount updates the active viewer count for a livestream room.
func (s *PresenceService) SetViewerCount(roomID string, count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.viewerCounts[roomID] = count
}

// GetViewerCount returns the active viewer count for a room.
func (s *PresenceService) GetViewerCount(roomID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.viewerCounts[roomID]
}

// GetOnlineUsers returns all currently online user IDs.
func (s *PresenceService) GetOnlineUsers() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]int64, 0, len(s.onlineUsers))
	for userID := range s.onlineUsers {
		users = append(users, userID)
	}
	return users
}
