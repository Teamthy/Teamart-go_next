package auth

import (
	"context"
	"time"
)

// SessionRepository defines the interface for session persistence
type SessionRepository interface {
	// CreateSession creates a new session in the repository
	CreateSession(ctx context.Context, session *Session) error

	// GetSession retrieves a session by ID
	GetSession(ctx context.Context, sessionID string) (*Session, error)

	// GetUserSessions retrieves all active sessions for a user
	GetUserSessions(ctx context.Context, userID int64) ([]*Session, error)

	// GetUserSession retrieves a specific session for a user
	GetUserSession(ctx context.Context, userID int64, sessionID string) (*Session, error)

	// UpdateSession updates an existing session
	UpdateSession(ctx context.Context, session *Session) error

	// TouchSession updates the last activity timestamp
	TouchSession(ctx context.Context, sessionID string) error

	// RevokeSession marks a session as revoked
	RevokeSession(ctx context.Context, sessionID string, reason string) error

	// RevokeUserSessions revokes all sessions for a user except optionally one
	RevokeUserSessions(ctx context.Context, userID int64, exceptSessionID string, reason string) error

	// CleanupExpiredSessions deletes expired sessions (older than given age)
	CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) (int64, error)

	// SessionExists checks if a session exists and is valid
	SessionExists(ctx context.Context, sessionID string) (bool, error)
}
