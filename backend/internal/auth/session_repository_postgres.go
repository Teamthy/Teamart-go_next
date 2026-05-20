package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

// SessionRepositoryPostgres is the PostgreSQL implementation of SessionRepository
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

// CreateSession creates a new session in the repository
func (r *SessionRepositoryPostgres) CreateSession(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			created_at, last_activity_at, expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
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
		false,                   // requires_password_verification
		nil, nil, nil, nil, nil, // geolocation fields
		session.CreatedAt,
		session.LastActivityAt,
		session.ExpiresAt,
	)

	if err != nil {
		r.logger.Errorf("failed to create session: %v", err)
		return fmt.Errorf("failed to create session: %w", err)
	}

	r.logger.Infof("session created: %s for user %d", session.ID, session.UserID)
	return nil
}

// GetSession retrieves a session by ID
func (r *SessionRepositoryPostgres) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	query := `
		SELECT id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, password_verified_at,
			revoked_at, revoke_reason, created_at, last_activity_at, expires_at
		FROM sessions
		WHERE id = $1
	`

	session := &Session{}
	var geoCountry, geoCity, geoTimezone *string
	var geoLat, geoLon *float64

	err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.DeviceFingerprint,
		&session.UserAgent,
		&session.IPAddress,
		&session.TrustLevel,
		&session.RequiresMFA,
		&session.RequiresMFA,
		&geoCountry,
		&geoCity,
		&geoLat,
		&geoLon,
		&geoTimezone,
		&session.MFAVerifiedAt,
		&session.LastActivityAt,
		&session.RevokedAt,
		&session.RevokeReason,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.ExpiresAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		r.logger.Errorf("failed to get session: %v", err)
		return nil, err
	}

	// Reconstruct geolocation
	if geoCountry != nil || geoCity != nil || geoLat != nil || geoLon != nil {
		session.GeoLocation = &GeoLocation{
			Country:   *geoCountry,
			City:      *geoCity,
			Latitude:  *geoLat,
			Longitude: *geoLon,
			Timezone:  *geoTimezone,
		}
	}

	return session, nil
}

// GetUserSessions retrieves all active sessions for a user
func (r *SessionRepositoryPostgres) GetUserSessions(ctx context.Context, userID int64) ([]*Session, error) {
	query := `
		SELECT id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, password_verified_at,
			revoked_at, revoke_reason, created_at, last_activity_at, expires_at
		FROM sessions
		WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		ORDER BY last_activity_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Errorf("failed to get user sessions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		var geoCountry, geoCity, geoTimezone *string
		var geoLat, geoLon *float64

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.DeviceID,
			&session.DeviceFingerprint,
			&session.UserAgent,
			&session.IPAddress,
			&session.TrustLevel,
			&session.RequiresMFA,
			&session.RequiresMFA,
			&geoCountry,
			&geoCity,
			&geoLat,
			&geoLon,
			&geoTimezone,
			&session.MFAVerifiedAt,
			&session.LastActivityAt,
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

		// Reconstruct geolocation
		if geoCountry != nil || geoCity != nil {
			session.GeoLocation = &GeoLocation{
				Country:   *geoCountry,
				City:      *geoCity,
				Latitude:  *geoLat,
				Longitude: *geoLon,
				Timezone:  *geoTimezone,
			}
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf("error iterating sessions: %v", err)
		return nil, err
	}

	return sessions, nil
}

// GetUserSession retrieves a specific session for a user
func (r *SessionRepositoryPostgres) GetUserSession(ctx context.Context, userID int64, sessionID string) (*Session, error) {
	query := `
		SELECT id, user_id, device_id, device_fingerprint, user_agent, ip_address,
			trust_level, requires_mfa_step, requires_password_verification,
			geo_country, geo_city, geo_latitude, geo_longitude, geo_timezone,
			mfa_verified_at, password_verified_at,
			revoked_at, revoke_reason, created_at, last_activity_at, expires_at
		FROM sessions
		WHERE id = $1 AND user_id = $2
	`

	session := &Session{}
	var geoCountry, geoCity, geoTimezone *string
	var geoLat, geoLon *float64

	err := r.db.QueryRow(ctx, query, sessionID, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.DeviceFingerprint,
		&session.UserAgent,
		&session.IPAddress,
		&session.TrustLevel,
		&session.RequiresMFA,
		&session.RequiresMFA,
		&geoCountry,
		&geoCity,
		&geoLat,
		&geoLon,
		&geoTimezone,
		&session.MFAVerifiedAt,
		&session.LastActivityAt,
		&session.RevokedAt,
		&session.RevokeReason,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.ExpiresAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		r.logger.Errorf("failed to get user session: %v", err)
		return nil, err
	}

	return session, nil
}

// UpdateSession updates an existing session
func (r *SessionRepositoryPostgres) UpdateSession(ctx context.Context, session *Session) error {
	query := `
		UPDATE sessions
		SET trust_level = $1, requires_mfa_step = $2, requires_password_verification = $3,
			mfa_verified_at = $4, password_verified_at = $5
		WHERE id = $6
	`

	tag, err := r.db.Exec(ctx, query,
		session.TrustLevel,
		session.RequiresMFA,
		session.RequiresMFA,
		session.MFAVerifiedAt,
		session.LastActivityAt,
		session.ID,
	)

	if err != nil {
		r.logger.Errorf("failed to update session: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// TouchSession updates the last activity timestamp
func (r *SessionRepositoryPostgres) TouchSession(ctx context.Context, sessionID string) error {
	query := `UPDATE sessions SET last_activity_at = $1 WHERE id = $2`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, now, sessionID)
	if err != nil {
		r.logger.Errorf("failed to touch session: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// RevokeSession marks a session as revoked
func (r *SessionRepositoryPostgres) RevokeSession(ctx context.Context, sessionID string, reason string) error {
	query := `
		UPDATE sessions
		SET revoked_at = $1, revoke_reason = $2
		WHERE id = $3
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, now, reason, sessionID)
	if err != nil {
		r.logger.Errorf("failed to revoke session: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}

	r.logger.Infof("session revoked: %s (reason: %s)", sessionID, reason)
	return nil
}

// RevokeUserSessions revokes all sessions for a user except optionally one
func (r *SessionRepositoryPostgres) RevokeUserSessions(ctx context.Context, userID int64, exceptSessionID string, reason string) error {
	var query string
	var tag any

	now := time.Now()

	if exceptSessionID != "" {
		query = `
			UPDATE sessions
			SET revoked_at = $1, revoke_reason = $2
			WHERE user_id = $3 AND id != $4 AND revoked_at IS NULL
		`
		var err error
		tag, err = r.db.Exec(ctx, query, now, reason, userID, exceptSessionID)
		if err != nil {
			r.logger.Errorf("failed to revoke user sessions: %v", err)
			return err
		}
	} else {
		query = `
			UPDATE sessions
			SET revoked_at = $1, revoke_reason = $2
			WHERE user_id = $3 AND revoked_at IS NULL
		`
		var err error
		tag, err = r.db.Exec(ctx, query, now, reason, userID)
		if err != nil {
			r.logger.Errorf("failed to revoke user sessions: %v", err)
			return err
		}
	}

	r.logger.Infof("revoked all sessions for user %d (reason: %s)", userID, reason)
	return nil
}

// CleanupExpiredSessions deletes expired sessions older than given age
func (r *SessionRepositoryPostgres) CleanupExpiredSessions(ctx context.Context, maxAge time.Duration) (int64, error) {
	query := `
		DELETE FROM sessions
		WHERE (revoked_at IS NOT NULL OR expires_at < NOW())
		AND created_at < NOW() - INTERVAL '1 second' * $1
	`

	tag, err := r.db.Exec(ctx, query, maxAge.Seconds())
	if err != nil {
		r.logger.Errorf("failed to cleanup expired sessions: %v", err)
		return 0, err
	}

	deleted := tag.RowsAffected()
	r.logger.Infof("cleaned up %d expired sessions", deleted)
	return deleted, nil
}

// SessionExists checks if a session exists and is valid
func (r *SessionRepositoryPostgres) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM sessions
			WHERE id = $1 AND revoked_at IS NULL AND expires_at > NOW()
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, sessionID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("failed to check session existence: %v", err)
		return false, err
	}

	return exists, nil
}
