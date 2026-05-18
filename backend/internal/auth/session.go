package auth

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// SessionService manages user sessions and device trust
type SessionService struct {
	config *AuthConfig
	logger *logger.Logger
}

// NewSessionService creates a new session service
func NewSessionService(config *AuthConfig, logger *logger.Logger) *SessionService {
	return &SessionService{
		config: config,
		logger: logger,
	}
}

// CreateSessionInput represents input for session creation
type CreateSessionInput struct {
	UserID        int64
	DeviceID      string
	DeviceName    string
	DeviceType    string
	IPAddress     string
	UserAgent     string
}

// CreateSessionOutput represents the result of session creation
type CreateSessionOutput struct {
	Session   *Session
	AccessTokenJTI   string
	RefreshTokenJTI  string
}

// CreateSession creates a new session
func (ss *SessionService) CreateSession(ctx context.Context, input *CreateSessionInput) (*CreateSessionOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.IPAddress == "" {
		return nil, fmt.Errorf("IP address is required")
	}

	// Generate device fingerprint
	fingerprint := ss.generateDeviceFingerprint(input.UserAgent, input.IPAddress)
	
	// Generate session ID
	sessionID := ss.generateSessionID()

	session := &Session{
		ID:                sessionID,
		UserID:            input.UserID,
		DeviceID:          input.DeviceID,
		DeviceFingerprint: fingerprint,
		IPAddress:         input.IPAddress,
		UserAgent:         input.UserAgent,
		ExpiresAt:         time.Now().Add(ss.config.SessionTTL),
		LastActivityAt:    time.Now(),
		IsRevoked:         false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	ss.logger.Infof("session created for user %d on device %s from IP %s", 
		input.UserID, input.DeviceID, input.IPAddress)

	return &CreateSessionOutput{
		Session: session,
	}, nil
}

// ValidateSessionInput represents input for session validation
type ValidateSessionInput struct {
	SessionID         string
	UserID            int64
	IPAddress         string
	UserAgent         string
}

// ValidateSessionOutput represents the result of session validation
type ValidateSessionOutput struct {
	IsValid      bool
	Session      *Session
	IsTrusted    bool
	Error        error
}

// ValidateSession validates a session
func (ss *SessionService) ValidateSession(ctx context.Context, input *ValidateSessionInput) (*ValidateSessionOutput, error) {
	if input.SessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	// In a real implementation, this would:
	// 1. Fetch session from database
	// 2. Check if session is revoked
	// 3. Check if session is expired
	// 4. Check if device fingerprint matches
	// 5. Check if idle timeout exceeded

	ss.logger.Debugf("validating session %s for user %d", input.SessionID, input.UserID)

	return &ValidateSessionOutput{
		IsValid: true,
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

	ss.logger.Debugf("updating activity for session %s", input.SessionID)

	// In a real implementation, this would update the last_activity_at timestamp

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

	ss.logger.Infof("revoking session %s for user %d (reason: %s)", 
		input.SessionID, input.UserID, input.Reason)

	// In a real implementation, this would mark the session as revoked

	return nil
}

// RevokeAllUserSessionsInput represents input for revoking all user sessions
type RevokeAllUserSessionsInput struct {
	UserID    int64
	ExceptSessionID string // Keep this session active
	Reason    string
}

// RevokeAllUserSessions revokes all sessions for a user
func (ss *SessionService) RevokeAllUserSessions(ctx context.Context, input *RevokeAllUserSessionsInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	ss.logger.Infof("revoking all sessions for user %d (reason: %s)", input.UserID, input.Reason)

	// In a real implementation, this would revoke all sessions except the specified one

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
	Sessions []Session
	Total    int32
	Limit    int32
	Offset   int32
}

// ListUserSessions lists all active sessions for a user
func (ss *SessionService) ListUserSessions(ctx context.Context, input *ListUserSessionsInput) (*ListUserSessionsOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	ss.logger.Debugf("listing sessions for user %d", input.UserID)

	// In a real implementation, this would fetch sessions from database

	return &ListUserSessionsOutput{
		Sessions: make([]Session, 0),
		Total:    0,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}, nil
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
	IsTrusted bool
	RequireVerification bool
}

// VerifyDevice checks if a device is trusted
func (ss *SessionService) VerifyDevice(ctx context.Context, input *VerifyDeviceInput) (*VerifyDeviceOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.DeviceID == "" {
		return nil, fmt.Errorf("device ID is required")
	}

	ss.logger.Debugf("verifying device %s for user %d", input.DeviceID, input.UserID)

	// In a real implementation, this would:
	// 1. Look up device trust record in database
	// 2. Check if device is marked as trusted
	// 3. Check if device fingerprint matches
	// 4. Check if device is revoked

	return &VerifyDeviceOutput{
		IsTrusted:          false,
		RequireVerification: ss.config.RequireDeviceVerification,
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

// TrustDevice marks a device as trusted
func (ss *SessionService) TrustDevice(ctx context.Context, input *TrustDeviceInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if input.DeviceID == "" {
		return fmt.Errorf("device ID is required")
	}

	ss.logger.Infof("trusting device %s (%s) for user %d", input.DeviceID, input.DeviceName, input.UserID)

	// In a real implementation, this would create/update device trust record

	return nil
}

// ===== Helper Methods =====

// generateSessionID generates a unique session ID
func (ss *SessionService) generateSessionID() string {
	b := make([]byte, 32)
	_, err := fmt.Sprintf("%x", b)
	if err != nil {
		// Fallback
		return fmt.Sprintf("session_%d_%d", time.Now().Unix(), time.Now().Nanosecond())
	}
	return fmt.Sprintf("sess_%x", b)
}

// generateDeviceFingerprint generates a fingerprint for a device
func (ss *SessionService) generateDeviceFingerprint(userAgent, ipAddress string) string {
	data := fmt.Sprintf("%s:%s", userAgent, ipAddress)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// IsSessionExpired checks if a session is expired
func (ss *SessionService) IsSessionExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// IsSessionIdleTimeoutExceeded checks if session idle timeout is exceeded
func (ss *SessionService) IsSessionIdleTimeoutExceeded(lastActivityAt time.Time) bool {
	return time.Since(lastActivityAt) > ss.config.SessionIdleTimeout
}
