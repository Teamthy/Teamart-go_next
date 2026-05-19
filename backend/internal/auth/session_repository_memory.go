package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// SessionRepositoryMemory implements SessionRepository using in-memory storage
// This is useful for testing and non-production environments
type SessionRepositoryMemory struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	logger   *logger.Logger
}

// NewSessionRepositoryMemory creates a new in-memory session repository
func NewSessionRepositoryMemory(logger *logger.Logger) *SessionRepositoryMemory {
	return &SessionRepositoryMemory{
		sessions: make(map[string]*Session),
		logger:   logger,
	}
}

// CreateSession creates a new session in memory
func (r *SessionRepositoryMemory) CreateSession(ctx context.Context, session *Session) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	if session.ID == "" {
		return fmt.Errorf("session ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[session.ID]; exists {
		return fmt.Errorf("session already exists: %s", session.ID)
	}

	// Make a copy to avoid external mutations
	sessionCopy := *session
	r.sessions[session.ID] = &sessionCopy

	r.logger.Debugf("created session %s for user %d", session.ID, session.UserID)
	return nil
}

// GetSession retrieves a session by ID
func (r *SessionRepositoryMemory) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	// Return a copy to avoid external mutations
	sessionCopy := *session
	return &sessionCopy, nil
}

// GetUserSessions retrieves all active sessions for a user
func (r *SessionRepositoryMemory) GetUserSessions(ctx context.Context, userID int64) ([]*Session, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var sessions []*Session
	now := time.Now()

	for _, session := range r.sessions {
		if session.UserID == userID && !session.IsRevoked && now.Before(session.ExpiresAt) {
			sessionCopy := *session
			sessions = append(sessions, &sessionCopy)
		}
	}

	r.logger.Debugf("found %d active sessions for user %d", len(sessions), userID)
	return sessions, nil
}

// GetUserSession retrieves a specific session for a user
func (r *SessionRepositoryMemory) GetUserSession(ctx context.Context, userID int64, sessionID string) (*Session, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	if session.UserID != userID {
		return nil, fmt.Errorf("session does not belong to user %d", userID)
	}

	sessionCopy := *session
	return &sessionCopy, nil
}

// UpdateSession updates an existing session
func (r *SessionRepositoryMemory) UpdateSession(ctx context.Context, session *Session) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	if session.ID == "" {
		return fmt.Errorf("session ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.sessions[session.ID]; !exists {
		return fmt.Errorf("session not found: %s", session.ID)
	}

	// Update timestamp
	session.UpdatedAt = time.Now()

	// Make a copy to avoid external mutations
	sessionCopy := *session
	r.sessions[session.ID] = &sessionCopy

	r.logger.Debugf("updated session %s", session.ID)
	return nil
}

// TouchSession updates the last activity timestamp
func (r *SessionRepositoryMemory) TouchSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.LastActivityAt = time.Now()
	session.UpdatedAt = time.Now()

	r.logger.Debugf("touched session %s", sessionID)
	return nil
}

// RevokeSession marks a session as revoked
func (r *SessionRepositoryMemory) RevokeSession(ctx context.Context, sessionID string, reason string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	now := time.Now()
	session.IsRevoked = true
	session.RevokedAt = &now
	session.RevokeReason = reason
	session.UpdatedAt = now

	r.logger.Infof("revoked session %s: %s", sessionID, reason)
	return nil
}

// RevokeUserSessions revokes all sessions for a user except optionally one
func (r *SessionRepositoryMemory) RevokeUserSessions(ctx context.Context, userID int64, exceptSessionID string, reason string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	revokedCount := 0
	now := time.Now()

	for _, session := range r.sessions {
		if session.UserID == userID && session.ID != exceptSessionID && !session.IsRevoked {
			session.IsRevoked = true
			session.RevokedAt = &now
			session.RevokeReason = reason
			session.UpdatedAt = now
			revokedCount++
		}
	}

	r.logger.Infof("revoked %d sessions for user %d: %s", revokedCount, userID, reason)
	return nil
}

// CleanupExpiredSessions deletes expired sessions
func (r *SessionRepositoryMemory) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-maxAge)
	deletedCount := int64(0)

	for sessionID, session := range r.sessions {
		if session.IsRevoked && session.RevokedAt != nil && session.RevokedAt.Before(cutoff) {
			delete(r.sessions, sessionID)
			deletedCount++
		} else if !session.IsRevoked && session.ExpiresAt.Before(cutoff) {
			delete(r.sessions, sessionID)
			deletedCount++
		}
	}

	r.logger.Infof("cleaned up %d expired sessions", deletedCount)
	return deletedCount, nil
}

// SessionExists checks if a session exists and is valid
func (r *SessionRepositoryMemory) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	if sessionID == "" {
		return false, fmt.Errorf("session ID is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return false, nil
	}

	// Check if valid (not revoked and not expired)
	if session.IsRevoked || time.Now().After(session.ExpiresAt) {
		return false, nil
	}

	return true, nil
}
