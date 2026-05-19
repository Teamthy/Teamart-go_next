package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

// SessionRepositoryPostgres implements SessionRepository using PostgreSQL
type SessionRepositoryPostgres struct {
	db     *database.Pool
	logger *logger.Logger
}

// NewSessionRepositoryPostgres creates a new PostgreSQL session repository
func NewSessionRepositoryPostgres(db *database.Pool, logger *logger.Logger) *SessionRepositoryPostgres {
	return &SessionRepositoryPostgres{
		db:     db,
		logger: logger,
	}
}

// CreateSession creates a new session in PostgreSQL
func (r *SessionRepositoryPostgres) CreateSession(ctx context.Context, session *Session) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	if session.ID == "" {
		return fmt.Errorf("session ID is required")
	}
	if session.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	query := `
		INSERT INTO sessions (
			id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, password_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18,
			$19, $20, $21
		)
	`

	_, err := r.db.Exec(ctx, query,
		session.ID,
		session.UserID,
		session.DeviceID,
		session.DeviceFingerprint,
		session.UserAgent,
		session.IPAddress,
		session.TrustLevel,
		session.RequiresMFAStep,
		session.RequiresPasswordVerification,
		session.GeoLocation.Country,
		session.GeoLocation.City,
		session.GeoLocation.Latitude,
		session.GeoLocation.Longitude,
		session.GeoLocation.Timezone,
		session.MFAVerifiedAt,
		session.PasswordVerifiedAt,
		session.RevokedAt,
		session.RevokeReason,
		session.CreatedAt,
		session.LastActivityAt,
		session.ExpiresAt,
	)

	if err != nil {
		r.logger.Errorf("failed to create session: %v", err)
		return fmt.Errorf("failed to create session: %w", err)
	}

	r.logger.Debugf("created session %s for user %d", session.ID, session.UserID)
	return nil
}

// GetSession retrieves a session by ID
func (r *SessionRepositoryPostgres) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	query := `
		SELECT
			id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, password_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		FROM sessions
		WHERE id = $1
	`

	session := &Session{
		GeoLocation: &GeoLocation{},
	}

	err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.DeviceFingerprint,
		&session.UserAgent,
		&session.IPAddress,
		&session.TrustLevel,
		&session.RequiresMFAStep,
		&session.RequiresPasswordVerification,
		&session.GeoLocation.Country,
		&session.GeoLocation.City,
		&session.GeoLocation.Latitude,
		&session.GeoLocation.Longitude,
		&session.GeoLocation.Timezone,
		&session.MFAVerifiedAt,
		&session.PasswordVerifiedAt,
		&session.RevokedAt,
		&session.RevokeReason,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.ExpiresAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	if err != nil {
		r.logger.Errorf("failed to get session: %v", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// GetUserSessions retrieves all active sessions for a user
func (r *SessionRepositoryPostgres) GetUserSessions(ctx context.Context, userID int64) ([]*Session, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	query := `
		SELECT
			id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, password_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		FROM sessions
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > $2
		ORDER BY last_activity_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID, time.Now())
	if err != nil {
		r.logger.Errorf("failed to get user sessions: %v", err)
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{
			GeoLocation: &GeoLocation{},
		}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.DeviceID,
			&session.DeviceFingerprint,
			&session.UserAgent,
			&session.IPAddress,
			&session.TrustLevel,
			&session.RequiresMFAStep,
			&session.RequiresPasswordVerification,
			&session.GeoLocation.Country,
			&session.GeoLocation.City,
			&session.GeoLocation.Latitude,
			&session.GeoLocation.Longitude,
			&session.GeoLocation.Timezone,
			&session.MFAVerifiedAt,
			&session.PasswordVerifiedAt,
			&session.RevokedAt,
			&session.RevokeReason,
			&session.CreatedAt,
			&session.LastActivityAt,
			&session.ExpiresAt,
		)
		if err != nil {
			r.logger.Errorf("failed to scan session row: %v", err)
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf("error iterating sessions: %v", err)
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// GetUserSession retrieves a specific session for a user
func (r *SessionRepositoryPostgres) GetUserSession(ctx context.Context, userID int64, sessionID string) (*Session, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is required")
	}

	query := `
		SELECT
			id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, password_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		FROM sessions
		WHERE id = $1 AND user_id = $2
	`

	session := &Session{
		GeoLocation: &GeoLocation{},
	}

	err := r.db.QueryRow(ctx, query, sessionID, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.DeviceFingerprint,
		&session.UserAgent,
		&session.IPAddress,
		&session.TrustLevel,
		&session.RequiresMFAStep,
		&session.RequiresPasswordVerification,
		&session.GeoLocation.Country,
		&session.GeoLocation.City,
		&session.GeoLocation.Latitude,
		&session.GeoLocation.Longitude,
		&session.GeoLocation.Timezone,
		&session.MFAVerifiedAt,
		&session.PasswordVerifiedAt,
		&session.RevokedAt,
		&session.RevokeReason,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.ExpiresAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found for user %d: %s", userID, sessionID)
	}
	if err != nil {
		r.logger.Errorf("failed to get user session: %v", err)
		return nil, fmt.Errorf("failed to get user session: %w", err)
	}

	return session, nil
}

// UpdateSession updates an existing session
func (r *SessionRepositoryPostgres) UpdateSession(ctx context.Context, session *Session) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}
	if session.ID == "" {
		return fmt.Errorf("session ID is required")
	}

	query := `
		UPDATE sessions SET
			trust_level = $1,
			requires_mfa_step = $2,
			requires_password_verification = $3,
			geo_country = $4,
			geo_city = $5,
			geo_latitude = $6,
			geo_longitude = $7,
			geo_timezone = $8,
			mfa_verified_at = $9,
			password_verified_at = $10,
			revoked_at = $11,
			revoke_reason = $12,
			last_activity_at = $13,
			expires_at = $14
		WHERE id = $15
	`

	tag, err := r.db.Exec(ctx, query,
		session.TrustLevel,
		session.RequiresMFAStep,
		session.RequiresPasswordVerification,
		session.GeoLocation.Country,
		session.GeoLocation.City,
		session.GeoLocation.Latitude,
		session.GeoLocation.Longitude,
		session.GeoLocation.Timezone,
		session.MFAVerifiedAt,
		session.PasswordVerifiedAt,
		session.RevokedAt,
		session.RevokeReason,
		session.LastActivityAt,
		session.ExpiresAt,
		session.ID,
	)

	if err != nil {
		r.logger.Errorf("failed to update session: %v", err)
		return fmt.Errorf("failed to update session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session not found: %s", session.ID)
	}

	r.logger.Debugf("updated session %s", session.ID)
	return nil
}

// TouchSession updates the last activity timestamp
func (r *SessionRepositoryPostgres) TouchSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	query := `
		UPDATE sessions SET
			last_activity_at = $1
		WHERE id = $2 AND revoked_at IS NULL
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, now, sessionID)

	if err != nil {
		r.logger.Errorf("failed to touch session: %v", err)
		return fmt.Errorf("failed to touch session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session not found or revoked: %s", sessionID)
	}

	r.logger.Debugf("touched session %s", sessionID)
	return nil
}

// RevokeSession revokes a single session
func (r *SessionRepositoryPostgres) RevokeSession(ctx context.Context, sessionID string, reason string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	query := `
		UPDATE sessions SET
			revoked_at = $1,
			revoke_reason = $2
		WHERE id = $3
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, now, reason, sessionID)

	if err != nil {
		r.logger.Errorf("failed to revoke session: %v", err)
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	r.logger.Infof("revoked session %s: %s", sessionID, reason)
	return nil
}

// RevokeUserSessions revokes all user sessions except the provided one
func (r *SessionRepositoryPostgres) RevokeUserSessions(ctx context.Context, userID int64, exceptSessionID string, reason string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	query := `
		UPDATE sessions SET
			revoked_at = $1,
			revoke_reason = $2
		WHERE user_id = $3 AND id != $4 AND revoked_at IS NULL
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, now, reason, userID, exceptSessionID)

	if err != nil {
		r.logger.Errorf("failed to revoke user sessions: %v", err)
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}

	r.logger.Warnf("revoked %d sessions for user %d except %s: %s",
		tag.RowsAffected(), userID, exceptSessionID, reason)
	return nil
}

// CleanupExpiredSessions removes expired sessions older than maxAge
func (r *SessionRepositoryPostgres) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) error {
	query := `
		DELETE FROM sessions
		WHERE expires_at < $1 OR created_at < $2
	`

	cutoff := time.Now().Add(-maxAge)
	oldCutoff := time.Now().Add(-24 * time.Hour) // Also delete sessions older than 24h

	tag, err := r.db.Exec(ctx, query, cutoff, oldCutoff)

	if err != nil {
		r.logger.Errorf("failed to cleanup expired sessions: %v", err)
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	r.logger.Infof("cleaned up %d expired sessions", tag.RowsAffected())
	return nil
}

// SessionExists checks if a session exists and is valid
func (r *SessionRepositoryPostgres) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	if sessionID == "" {
		return false, fmt.Errorf("session ID is required")
	}

	query := `
		SELECT 1 FROM sessions
		WHERE id = $1 AND revoked_at IS NULL AND expires_at > $2
		LIMIT 1
	`

	err := r.db.QueryRow(ctx, query, sessionID, time.Now()).Scan(nil)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		r.logger.Errorf("failed to check session existence: %v", err)
		return false, fmt.Errorf("failed to check session existence: %w", err)
	}

	return true, nil
}
