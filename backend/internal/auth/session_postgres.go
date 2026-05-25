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
			trust_level, requires_mfa,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8,
			$9, $10, $11, $12, $13,
			$14, $15, $16,
			$17, $18, $19
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
		session.RequiresMFA,
		session.GeoLocation.Country,
		session.GeoLocation.City,
		session.GeoLocation.Latitude,
		session.GeoLocation.Longitude,
		session.GeoLocation.Timezone,
		session.MFAVerifiedAt,
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
			trust_level, requires_mfa,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		FROM sessions
		WHERE id = $1 AND revoked_at IS NULL
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
		&session.RequiresMFA,
		&session.GeoLocation.Country,
		&session.GeoLocation.City,
		&session.GeoLocation.Latitude,
		&session.GeoLocation.Longitude,
		&session.GeoLocation.Timezone,
		&session.MFAVerifiedAt,
		&session.RevokedAt,
		&session.RevokeReason,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		r.logger.Warnf("failed to get session: %v", err)
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
			trust_level, requires_mfa,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		FROM sessions
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		ORDER BY last_activity_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Warnf("failed to get user sessions: %v", err)
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]*Session, 0)
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
			&session.RequiresMFA,
			&session.GeoLocation.Country,
			&session.GeoLocation.City,
			&session.GeoLocation.Latitude,
			&session.GeoLocation.Longitude,
			&session.GeoLocation.Timezone,
			&session.MFAVerifiedAt,
			&session.RevokedAt,
			&session.RevokeReason,
			&session.CreatedAt,
			&session.LastActivityAt,
			&session.ExpiresAt,
		)

		if err != nil {
			r.logger.Warnf("failed to scan session: %v", err)
			continue
		}

		sessions = append(sessions, session)
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
			trust_level, requires_mfa,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, revoked_at, revoke_reason,
			created_at, last_activity_at, expires_at
		FROM sessions
		WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL
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
		&session.RequiresMFA,
		&session.GeoLocation.Country,
		&session.GeoLocation.City,
		&session.GeoLocation.Latitude,
		&session.GeoLocation.Longitude,
		&session.GeoLocation.Timezone,
		&session.MFAVerifiedAt,
		&session.RevokedAt,
		&session.RevokeReason,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
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
			requires_mfa = $2,
			mfa_verified_at = $3,
			last_activity_at = $4,
			updated_at = NOW()
		WHERE id = $5 AND revoked_at IS NULL
	`

	result, err := r.db.Exec(ctx, query,
		session.TrustLevel,
		session.RequiresMFA,
		session.MFAVerifiedAt,
		session.LastActivityAt,
		session.ID,
	)

	if err != nil {
		r.logger.Errorf("failed to update session: %v", err)
		return fmt.Errorf("failed to update session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// TouchSession updates the last activity time for a session
func (r *SessionRepositoryPostgres) TouchSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	query := `
		UPDATE sessions SET
			last_activity_at = NOW(),
			updated_at = NOW()
		WHERE id = $1 AND revoked_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, sessionID)
	if err != nil {
		r.logger.Errorf("failed to touch session: %v", err)
		return fmt.Errorf("failed to touch session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// RevokeSession revokes a session
func (r *SessionRepositoryPostgres) RevokeSession(ctx context.Context, sessionID string, reason string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	query := `
		UPDATE sessions SET
			revoked_at = NOW(),
			revoke_reason = $1,
			is_revoked = true,
			updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, reason, sessionID)
	if err != nil {
		r.logger.Errorf("failed to revoke session: %v", err)
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	r.logger.Infof("revoked session %s: %s", sessionID, reason)
	return nil
}

// RevokeUserSessions revokes all sessions for a user
func (r *SessionRepositoryPostgres) RevokeUserSessions(ctx context.Context, userID int64, exceptSessionID string, reason string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	query := `
		UPDATE sessions SET
			revoked_at = NOW(),
			revoke_reason = $1,
			is_revoked = true,
			updated_at = NOW()
		WHERE user_id = $2 AND id != $3 AND revoked_at IS NULL
	`

	result, err := r.db.Exec(ctx, query, reason, userID, exceptSessionID)
	if err != nil {
		r.logger.Errorf("failed to revoke user sessions: %v", err)
		return fmt.Errorf("failed to revoke user sessions: %w", err)
	}

	r.logger.Infof("revoked %d sessions for user %d: %s", result.RowsAffected(), userID, reason)
	return nil
}

// CleanupExpiredSessions removes expired and revoked sessions older than maxAge
func (r *SessionRepositoryPostgres) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) (int64, error) {
	query := `
		DELETE FROM sessions
		WHERE (revoked_at IS NOT NULL AND revoked_at < NOW() - INTERVAL '1 second' * $1)
		   OR (expires_at < NOW() AND created_at < NOW() - INTERVAL '1 second' * $1)
	`

	result, err := r.db.Exec(ctx, query, int64(maxAge.Seconds()))
	if err != nil {
		r.logger.Errorf("failed to cleanup expired sessions: %v", err)
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}

	return result.RowsAffected(), nil
}

// SessionExists checks whether a session exists and is still valid.
func (r *SessionRepositoryPostgres) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	if sessionID == "" {
		return false, fmt.Errorf("session ID is required")
	}

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM sessions
			WHERE id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		)
	`

	var exists bool
	if err := r.db.QueryRow(ctx, query, sessionID).Scan(&exists); err != nil {
		r.logger.Warnf("failed to check session existence: %v", err)
		return false, fmt.Errorf("failed to check session existence: %w", err)
	}

	return exists, nil
}
