package auth

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// SessionService manages user sessions and device trust
type SessionService struct {
	config      *AuthConfig
	logger      *logger.Logger
	repo        SessionRepository
	maxSessions int // Maximum concurrent sessions per user
}

// NewSessionService creates a new session service
func NewSessionService(config *AuthConfig, logger *logger.Logger, repo SessionRepository) *SessionService {
	return &SessionService{
		config:      config,
		logger:      logger,
		repo:        repo,
		maxSessions: 5, // Default: allow 5 concurrent sessions per user
	}
}

// SessionRepository returns the underlying session repository
func (ss *SessionService) SessionRepository() SessionRepository {
	return ss.repo
}

// CreateSessionInput represents input for session creation
type CreateSessionInput struct {
	UserID     int64
	DeviceID   string
	DeviceName string
	DeviceType string
	IPAddress  string
	UserAgent  string
}

// CreateSessionOutput represents the result of session creation
type CreateSessionOutput struct {
	Session *Session
}

// CreateSession creates a new session with proper validation and device fingerprinting
func (ss *SessionService) CreateSession(ctx context.Context, input *CreateSessionInput) (*CreateSessionOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.IPAddress == "" {
		return nil, fmt.Errorf("IP address is required")
	}
	if input.DeviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	fingerprint := ss.generateDeviceFingerprint(input.UserAgent, input.IPAddress)
	sessionID := ss.generateSessionID()

	session := &Session{
		ID:                sessionID,
		UserID:            input.UserID,
		DeviceID:          input.DeviceID,
		DeviceFingerprint: fingerprint,
		IPAddress:         input.IPAddress,
		UserAgent:         input.UserAgent,
		TrustLevel:        TrustLevelUntrusted,
		ExpiresAt:         time.Now().Add(ss.config.SessionTTL),
		LastActivityAt:    time.Now(),
		IsRevoked:         false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := ss.repo.CreateSession(ctx, session); err != nil {
		ss.logger.Errorf("failed to create session: %v", err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := ss.enforceSessionLimit(ctx, input.UserID); err != nil {
		ss.logger.Warnf("failed to enforce session limit: %v", err)
	}

	ss.logger.Infof("session created for user %d on device %s from IP %s",
		input.UserID, input.DeviceID, input.IPAddress)

	return &CreateSessionOutput{Session: session}, nil
}

// ValidateSessionInput represents input for session validation
type ValidateSessionInput struct {
	SessionID string
	UserID    int64
	IPAddress string
	UserAgent string
}

// ValidateSessionOutput represents the result of session validation
type ValidateSessionOutput struct {
	IsValid     bool
	Session     *Session
	IsTrusted   bool
	RequiresMFA bool
}

// ValidateSession validates a session with comprehensive checks
func (ss *SessionService) ValidateSession(ctx context.Context, input *ValidateSessionInput) (*ValidateSessionOutput, error) {
	if input.SessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	session, err := ss.repo.GetUserSession(ctx, input.UserID, input.SessionID)
	if err != nil {
		ss.logger.Warnf("session validation failed: %v", err)
		return &ValidateSessionOutput{IsValid: false}, nil
	}

	if session.IsRevoked {
		ss.logger.Warnf("session %s is revoked", session.ID)
		return &ValidateSessionOutput{IsValid: false, Session: session}, nil
	}

	if session.IsExpired() {
		ss.logger.Debugf("session %s is expired", session.ID)
		return &ValidateSessionOutput{IsValid: false, Session: session}, nil
	}

	if session.IsInactive(ss.config.SessionIdleTimeout) {
		ss.logger.Infof("session %s exceeded idle timeout", session.ID)
		if err := ss.repo.RevokeSession(ctx, session.ID, "idle_timeout"); err != nil {
			ss.logger.Errorf("failed to revoke idle session: %v", err)
		}
		return &ValidateSessionOutput{IsValid: false, Session: session}, nil
	}

	expectedFingerprint := ss.generateDeviceFingerprint(input.UserAgent, input.IPAddress)
	fingerprintMatches := expectedFingerprint == session.DeviceFingerprint
	ipChanged := input.IPAddress != session.IPAddress

	requiresMFA := false
	if session.TrustLevel == TrustLevelUntrusted {
		requiresMFA = true
	}
	if !fingerprintMatches || ipChanged {
		requiresMFA = true
		session.TrustLevel = TrustLevelPartial
	}

	if err := ss.repo.TouchSession(ctx, session.ID); err != nil {
		ss.logger.Warnf("failed to touch session: %v", err)
	}

	ss.logger.Debugf("session %s validated", session.ID)

	return &ValidateSessionOutput{
		IsValid:     true,
		Session:     session,
		IsTrusted:   session.TrustLevel == TrustLevelTrusted,
		RequiresMFA: requiresMFA,
	}, nil
}

// UpdateSessionActivityInput represents input for updating session activity
type UpdateSessionActivityInput struct {
	SessionID string
	UserID    int64
	IPAddress string
	UserAgent string
}

// UpdateSessionActivity updates the last activity time for a session
func (ss *SessionService) UpdateSessionActivity(ctx context.Context, input *UpdateSessionActivityInput) error {
	if input.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	if err := ss.repo.TouchSession(ctx, input.SessionID); err != nil {
		ss.logger.Errorf("failed to update session activity: %v", err)
		return fmt.Errorf("failed to update session activity: %w", err)
	}

	ss.logger.Debugf("updated activity for session %s", input.SessionID)
	return nil
}

// RevokeSessionInput represents input for session revocation
type RevokeSessionInput struct {
	SessionID string
	UserID    int64
	Reason    string
}

// RevokeSession revokes a session
func (ss *SessionService) RevokeSession(ctx context.Context, input *RevokeSessionInput) error {
	if input.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	session, err := ss.repo.GetUserSession(ctx, input.UserID, input.SessionID)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}

	if err := ss.repo.RevokeSession(ctx, session.ID, input.Reason); err != nil {
		ss.logger.Errorf("failed to revoke session: %v", err)
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	ss.logger.Infof("revoked session %s for user %d", input.SessionID, input.UserID)
	return nil
}

// RevokeAllUserSessionsInput represents input for revoking all user sessions
type RevokeAllUserSessionsInput struct {
	UserID          int64
	ExceptSessionID string
	Reason          string
}

// RevokeAllUserSessions revokes all sessions for a user
func (ss *SessionService) RevokeAllUserSessions(ctx context.Context, input *RevokeAllUserSessionsInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	if err := ss.repo.RevokeUserSessions(ctx, input.UserID, input.ExceptSessionID, input.Reason); err != nil {
		ss.logger.Errorf("failed to revoke all user sessions: %v", err)
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	ss.logger.Infof("revoked all sessions for user %d", input.UserID)
	return nil
}

// ListUserSessionsInput represents input for listing user sessions
type ListUserSessionsInput struct {
	UserID int64
	Limit  int32
	Offset int32
}

// ListUserSessionsOutput represents the result of listing sessions
type ListUserSessionsOutput struct {
	Sessions []*Session
	Total    int32
	Limit    int32
	Offset   int32
}

// GetUserActiveSessions lists all active sessions for a user
func (ss *SessionService) GetUserActiveSessions(ctx context.Context, userID int64) (*ListUserSessionsOutput, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	sessions, err := ss.repo.GetUserSessions(ctx, userID)
	if err != nil {
		ss.logger.Errorf("failed to list user sessions: %v", err)
		return nil, fmt.Errorf("failed to list user sessions: %w", err)
	}

	ss.logger.Debugf("retrieved %d active sessions for user %d", len(sessions), userID)

	return &ListUserSessionsOutput{
		Sessions: sessions,
		Total:    int32(len(sessions)),
		Limit:    100,
		Offset:   0,
	}, nil
}

// CleanupExpiredSessions removes expired and revoked sessions older than maxAge
func (ss *SessionService) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) error {
	deleted, err := ss.repo.CleanupExpiredSessions(ctx, maxAge)
	if err != nil {
		ss.logger.Errorf("failed to cleanup expired sessions: %v", err)
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	ss.logger.Infof("cleaned up %d expired sessions", deleted)
	return nil
}

// ===== Device Trust =====

// VerifyDeviceInput represents input for device verification
type VerifyDeviceInput struct {
	UserID            int64
	DeviceID          string
	DeviceFingerprint string
	IPAddress         string
	UserAgent         string
}

// VerifyDeviceOutput represents the result of device verification
type VerifyDeviceOutput struct {
	IsTrusted           bool
	RequireVerification bool
	TrustLevel          TrustLevel
}

// VerifyDevice checks if a device is trusted
func (ss *SessionService) VerifyDevice(ctx context.Context, input *VerifyDeviceInput) (*VerifyDeviceOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.DeviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	userSessions, err := ss.repo.GetUserSessions(ctx, input.UserID)
	if err != nil {
		ss.logger.Warnf("failed to check device trust: %v", err)
		return &VerifyDeviceOutput{
			IsTrusted:           false,
			RequireVerification: true,
			TrustLevel:          TrustLevelUntrusted,
		}, nil
	}

	for _, session := range userSessions {
		if session.DeviceID == input.DeviceID {
			return &VerifyDeviceOutput{
				IsTrusted:           session.TrustLevel == TrustLevelTrusted,
				RequireVerification: session.TrustLevel != TrustLevelTrusted,
				TrustLevel:          session.TrustLevel,
			}, nil
		}
	}

	ss.logger.Debugf("device %s not found in user %d's sessions", input.DeviceID, input.UserID)
	return &VerifyDeviceOutput{
		IsTrusted:           false,
		RequireVerification: true,
		TrustLevel:          TrustLevelUntrusted,
	}, nil
}

// TrustDeviceInput represents input for trusting a device
type TrustDeviceInput struct {
	UserID            int64
	DeviceID          string
	DeviceName        string
	DeviceType        string
	DeviceFingerprint string
	IPAddress         string
	UserAgent         string
}

// TrustDevice marks a device as trusted for a user's sessions
func (ss *SessionService) TrustDevice(ctx context.Context, input *TrustDeviceInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if input.DeviceID == "" {
		return fmt.Errorf("device ID is required")
	}

	userSessions, err := ss.repo.GetUserSessions(ctx, input.UserID)
	if err != nil {
		ss.logger.Errorf("failed to find sessions for device trust: %v", err)
		return fmt.Errorf("failed to find sessions: %w", err)
	}

	updatedCount := 0
	for _, session := range userSessions {
		if session.DeviceID == input.DeviceID {
			session.TrustLevel = TrustLevelTrusted
			if err := ss.repo.UpdateSession(ctx, session); err != nil {
				ss.logger.Warnf("failed to update session trust level: %v", err)
				continue
			}
			updatedCount++
		}
	}

	ss.logger.Infof("trusted device %s for user %d", input.DeviceID, input.UserID)
	return nil
}

// MarkDeviceSuspicious marks a device as suspicious/partially trusted
func (ss *SessionService) MarkDeviceSuspicious(ctx context.Context, userID int64, deviceID string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if deviceID == "" {
		return fmt.Errorf("device ID is required")
	}

	userSessions, err := ss.repo.GetUserSessions(ctx, userID)
	if err != nil {
		ss.logger.Errorf("failed to find sessions for device: %v", err)
		return fmt.Errorf("failed to find sessions: %w", err)
	}

	updatedCount := 0
	for _, session := range userSessions {
		if session.DeviceID == deviceID && session.TrustLevel != TrustLevelUntrusted {
			session.TrustLevel = TrustLevelPartial
			session.RequiresMFA = true
			if err := ss.repo.UpdateSession(ctx, session); err != nil {
				ss.logger.Warnf("failed to update session trust level: %v", err)
				continue
			}
			updatedCount++
		}
	}

	ss.logger.Warnf("marked device %s as suspicious for user %d", deviceID, userID)
	return nil
}

// ===== Helper Methods =====

// generateSessionID generates a unique session ID using cryptographic randomness
func (ss *SessionService) generateSessionID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("sess_%d_%d", time.Now().Unix(), time.Now().Nanosecond())
	}
	return fmt.Sprintf("sess_%x", b)
}

// generateDeviceFingerprint generates a fingerprint for a device using MD5
func (ss *SessionService) generateDeviceFingerprint(userAgent, ipAddress string) string {
	data := fmt.Sprintf("%s:%s", userAgent, ipAddress)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// enforceSessionLimit revokes oldest sessions if limit is exceeded
func (ss *SessionService) enforceSessionLimit(ctx context.Context, userID int64) error {
	sessions, err := ss.repo.GetUserSessions(ctx, userID)
	if err != nil {
		return err
	}

	if len(sessions) <= ss.maxSessions {
		return nil
	}

	sessionsToRevoke := len(sessions) - ss.maxSessions

	for i := 0; i < sessionsToRevoke && i < len(sessions); i++ {
		if err := ss.repo.RevokeSession(ctx, sessions[i].ID, "session_limit_exceeded"); err != nil {
			ss.logger.Warnf("failed to revoke session for limit enforcement: %v", err)
		}
	}

	ss.logger.Infof("enforced session limit for user %d: revoked %d sessions",
		userID, sessionsToRevoke)

	return nil
}
