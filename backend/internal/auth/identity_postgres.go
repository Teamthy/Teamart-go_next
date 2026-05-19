package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

// IdentityRepositoryPostgres implements IdentityRepository using PostgreSQL
type IdentityRepositoryPostgres struct {
	db     *database.Pool
	logger *logger.Logger
}

// NewIdentityRepositoryPostgres creates a new PostgreSQL identity repository
func NewIdentityRepositoryPostgres(db *database.Pool, logger *logger.Logger) *IdentityRepositoryPostgres {
	return &IdentityRepositoryPostgres{
		db:     db,
		logger: logger,
	}
}

// CreateIdentity creates a new user identity in PostgreSQL
func (r *IdentityRepositoryPostgres) CreateIdentity(ctx context.Context, identity *UserIdentity) error {
	if identity == nil {
		return fmt.Errorf("identity cannot be nil")
	}
	if identity.Email == "" {
		return fmt.Errorf("email is required")
	}

	query := `
		INSERT INTO users (
			email, password_hash, onboarding_state, account_status, is_active,
			failed_login_attempts, locked_until, password_changed_at,
			requires_mfa, mfa_method, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12
		)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		identity.Email,
		identity.PasswordHash,
		identity.OnboardingState,
		identity.AccountStatus,
		identity.IsActive,
		identity.FailedLoginAttempts,
		identity.LockedUntil,
		identity.PasswordChangedAt,
		identity.RequiresMFA,
		identity.MFAMethod,
		now,
		now,
	).Scan(&identity.ID, &identity.CreatedAt, &identity.UpdatedAt)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return fmt.Errorf("failed to create identity: %w", err)
		}
		r.logger.Errorf("failed to create identity: %v", err)
		return fmt.Errorf("failed to create identity: %w", err)
	}

	r.logger.Debugf("created identity %d with email %s", identity.ID, identity.Email)
	return nil
}

// GetIdentityByID retrieves a user identity by ID
func (r *IdentityRepositoryPostgres) GetIdentityByID(ctx context.Context, userID int64) (*UserIdentity, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	query := `
		SELECT
			id, email, password_hash, onboarding_state, account_status, is_active,
			failed_login_attempts, failed_login_last_attempt, locked_until,
			password_changed_at, last_login_at, last_login_ip,
			recovery_email, phone_number, requires_mfa, mfa_method,
			created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	identity := &UserIdentity{}
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&identity.ID,
		&identity.Email,
		&identity.PasswordHash,
		&identity.OnboardingState,
		&identity.AccountStatus,
		&identity.IsActive,
		&identity.FailedLoginAttempts,
		&identity.FailedLoginLastAttempt,
		&identity.LockedUntil,
		&identity.PasswordChangedAt,
		&identity.LastLoginAt,
		&identity.LastLoginIP,
		&identity.RecoveryEmail,
		&identity.PhoneNumber,
		&identity.RequiresMFA,
		&identity.MFAMethod,
		&identity.CreatedAt,
		&identity.UpdatedAt,
		&identity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("identity not found: %d", userID)
	}
	if err != nil {
		r.logger.Errorf("failed to get identity by ID: %v", err)
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}

	return identity, nil
}

// GetIdentityByEmail retrieves a user identity by email
func (r *IdentityRepositoryPostgres) GetIdentityByEmail(ctx context.Context, email string) (*UserIdentity, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	query := `
		SELECT
			id, email, password_hash, onboarding_state, account_status, is_active,
			failed_login_attempts, failed_login_last_attempt, locked_until,
			password_changed_at, last_login_at, last_login_ip,
			recovery_email, phone_number, requires_mfa, mfa_method,
			created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	identity := &UserIdentity{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&identity.ID,
		&identity.Email,
		&identity.PasswordHash,
		&identity.OnboardingState,
		&identity.AccountStatus,
		&identity.IsActive,
		&identity.FailedLoginAttempts,
		&identity.FailedLoginLastAttempt,
		&identity.LockedUntil,
		&identity.PasswordChangedAt,
		&identity.LastLoginAt,
		&identity.LastLoginIP,
		&identity.RecoveryEmail,
		&identity.PhoneNumber,
		&identity.RequiresMFA,
		&identity.MFAMethod,
		&identity.CreatedAt,
		&identity.UpdatedAt,
		&identity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("identity not found for email: %s", email)
	}
	if err != nil {
		r.logger.Errorf("failed to get identity by email: %v", err)
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}

	return identity, nil
}

// UpdateIdentity updates an existing user identity
func (r *IdentityRepositoryPostgres) UpdateIdentity(ctx context.Context, identity *UserIdentity) error {
	if identity == nil {
		return fmt.Errorf("identity cannot be nil")
	}
	if identity.ID == 0 {
		return fmt.Errorf("identity ID is required")
	}

	query := `
		UPDATE users SET
			email = $1,
			password_hash = $2,
			onboarding_state = $3,
			account_status = $4,
			is_active = $5,
			failed_login_attempts = $6,
			failed_login_last_attempt = $7,
			locked_until = $8,
			password_changed_at = $9,
			last_login_at = $10,
			last_login_ip = $11,
			recovery_email = $12,
			phone_number = $13,
			requires_mfa = $14,
			mfa_method = $15,
			updated_at = $16
		WHERE id = $17 AND deleted_at IS NULL
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query,
		identity.Email,
		identity.PasswordHash,
		identity.OnboardingState,
		identity.AccountStatus,
		identity.IsActive,
		identity.FailedLoginAttempts,
		identity.FailedLoginLastAttempt,
		identity.LockedUntil,
		identity.PasswordChangedAt,
		identity.LastLoginAt,
		identity.LastLoginIP,
		identity.RecoveryEmail,
		identity.PhoneNumber,
		identity.RequiresMFA,
		identity.MFAMethod,
		now,
		identity.ID,
	)

	if err != nil {
		r.logger.Errorf("failed to update identity: %v", err)
		return fmt.Errorf("failed to update identity: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", identity.ID)
	}

	identity.UpdatedAt = now
	r.logger.Debugf("updated identity %d", identity.ID)
	return nil
}

// UpdateAccountStatus updates account status
func (r *IdentityRepositoryPostgres) UpdateAccountStatus(ctx context.Context, userID int64, status AccountStatus) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	query := `
		UPDATE users SET
			account_status = $1,
			is_active = $2,
			updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	isActive := (status == AccountStatusActive)
	tag, err := r.db.Exec(ctx, query, status, isActive, time.Now(), userID)

	if err != nil {
		r.logger.Errorf("failed to update account status: %v", err)
		return fmt.Errorf("failed to update account status: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", userID)
	}

	r.logger.Debugf("updated account status for user %d to %s", userID, status)
	return nil
}

// LockAccount locks the account temporarily
func (r *IdentityRepositoryPostgres) LockAccount(ctx context.Context, userID int64, lockDuration time.Duration) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	lockUntil := time.Now().Add(lockDuration)
	query := `
		UPDATE users SET
			locked_until = $1,
			updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	tag, err := r.db.Exec(ctx, query, lockUntil, time.Now(), userID)

	if err != nil {
		r.logger.Errorf("failed to lock account: %v", err)
		return fmt.Errorf("failed to lock account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", userID)
	}

	r.logger.Warnf("locked account for user %d until %v", userID, lockUntil)
	return nil
}

// UnlockAccount unlocks a locked account
func (r *IdentityRepositoryPostgres) UnlockAccount(ctx context.Context, userID int64) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	query := `
		UPDATE users SET
			locked_until = NULL,
			failed_login_attempts = 0,
			updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	tag, err := r.db.Exec(ctx, query, time.Now(), userID)

	if err != nil {
		r.logger.Errorf("failed to unlock account: %v", err)
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", userID)
	}

	r.logger.Infof("unlocked account for user %d", userID)
	return nil
}

// RecordFailedLoginAttempt records a failed login attempt
func (r *IdentityRepositoryPostgres) RecordFailedLoginAttempt(ctx context.Context, userID int64) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	query := `
		UPDATE users SET
			failed_login_attempts = failed_login_attempts + 1,
			failed_login_last_attempt = $1,
			updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, now, now, userID)

	if err != nil {
		r.logger.Errorf("failed to record login attempt: %v", err)
		return fmt.Errorf("failed to record login attempt: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", userID)
	}

	r.logger.Debugf("recorded failed login attempt for user %d", userID)
	return nil
}

// RecordSuccessfulLogin records a successful login
func (r *IdentityRepositoryPostgres) RecordSuccessfulLogin(ctx context.Context, userID int64, ipAddress string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	query := `
		UPDATE users SET
			failed_login_attempts = 0,
			failed_login_last_attempt = NULL,
			last_login_at = $1,
			last_login_ip = $2,
			updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, now, ipAddress, now, userID)

	if err != nil {
		r.logger.Errorf("failed to record successful login: %v", err)
		return fmt.Errorf("failed to record successful login: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", userID)
	}

	r.logger.Infof("recorded successful login for user %d from IP %s", userID, ipAddress)
	return nil
}

// UpdateOnboardingState updates the onboarding state
func (r *IdentityRepositoryPostgres) UpdateOnboardingState(ctx context.Context, userID int64, state OnboardingState) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	accountStatus := AccountStatusPending
	isActive := false

	if state == StateOnboarded {
		accountStatus = AccountStatusActive
		isActive = true
	}

	query := `
		UPDATE users SET
			onboarding_state = $1,
			account_status = $2,
			is_active = $3,
			updated_at = $4
		WHERE id = $5 AND deleted_at IS NULL
	`

	tag, err := r.db.Exec(ctx, query, state, accountStatus, isActive, time.Now(), userID)

	if err != nil {
		r.logger.Errorf("failed to update onboarding state: %v", err)
		return fmt.Errorf("failed to update onboarding state: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", userID)
	}

	r.logger.Infof("updated onboarding state for user %d to %s", userID, state)
	return nil
}

// UpdatePasswordHash updates the password hash
func (r *IdentityRepositoryPostgres) UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if passwordHash == "" {
		return fmt.Errorf("password hash is required")
	}

	query := `
		UPDATE users SET
			password_hash = $1,
			password_changed_at = $2,
			updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	now := time.Now()
	tag, err := r.db.Exec(ctx, query, passwordHash, now, now, userID)

	if err != nil {
		r.logger.Errorf("failed to update password hash: %v", err)
		return fmt.Errorf("failed to update password hash: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("identity not found: %d", userID)
	}

	r.logger.Debugf("updated password hash for user %d", userID)
	return nil
}

// IdentityExists checks if an identity exists
func (r *IdentityRepositoryPostgres) IdentityExists(ctx context.Context, userID int64) (bool, error) {
	if userID == 0 {
		return false, fmt.Errorf("user ID is required")
	}

	query := `SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL LIMIT 1`

	err := r.db.QueryRow(ctx, query, userID).Scan(nil)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		r.logger.Errorf("failed to check identity existence: %v", err)
		return false, fmt.Errorf("failed to check identity existence: %w", err)
	}

	return true, nil
}

// EmailExists checks if an email is already registered
func (r *IdentityRepositoryPostgres) EmailExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, fmt.Errorf("email is required")
	}

	query := `SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL LIMIT 1`

	err := r.db.QueryRow(ctx, query, email).Scan(nil)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		r.logger.Errorf("failed to check email existence: %v", err)
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return true, nil
}
