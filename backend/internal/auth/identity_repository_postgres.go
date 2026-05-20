package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

// IdentityRepositoryPostgres is the PostgreSQL implementation of IdentityRepository
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
	query := `
		INSERT INTO users (
			email, password_hash, onboarding_state, account_status, is_active,
			failed_login_attempts, locked_until, password_changed_at,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		identity.Email,
		identity.PasswordHash,
		identity.OnboardingState,
		identity.AccountStatus,
		identity.IsActive,
		0,   // failed_login_attempts
		nil, // locked_until
		now,
		now,
		now,
	).Scan(&identity.ID)

	if err != nil {
		r.logger.Errorf("failed to create identity: %v", err)
		if err.Error() == "duplicate key value violates unique constraint \"idx_users_email\"" {
			return fmt.Errorf("email already registered")
		}
		return err
	}

	identity.CreatedAt = now
	identity.UpdatedAt = now
	r.logger.Infof("created new identity with ID %d", identity.ID)
	return nil
}

// GetIdentityByID retrieves a user identity by ID
func (r *IdentityRepositoryPostgres) GetIdentityByID(ctx context.Context, userID int64) (*UserIdentity, error) {
	query := `
		SELECT id, email, password_hash, onboarding_state, account_status, is_active,
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
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		r.logger.Errorf("failed to get identity by ID %d: %v", userID, err)
		return nil, err
	}

	return identity, nil
}

// GetIdentityByEmail retrieves a user identity by email
func (r *IdentityRepositoryPostgres) GetIdentityByEmail(ctx context.Context, email string) (*UserIdentity, error) {
	query := `
		SELECT id, email, password_hash, onboarding_state, account_status, is_active,
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
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		r.logger.Errorf("failed to get identity by email: %v", err)
		return nil, err
	}

	return identity, nil
}

// UpdateIdentity updates an existing user identity
func (r *IdentityRepositoryPostgres) UpdateIdentity(ctx context.Context, identity *UserIdentity) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, onboarding_state = $3, account_status = $4,
			is_active = $5, failed_login_attempts = $6, failed_login_last_attempt = $7,
			locked_until = $8, password_changed_at = $9, last_login_at = $10,
			last_login_ip = $11, recovery_email = $12, phone_number = $13,
			requires_mfa = $14, mfa_method = $15, updated_at = $16
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
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	identity.UpdatedAt = now
	return nil
}

// UpdateAccountStatus updates the account status
func (r *IdentityRepositoryPostgres) UpdateAccountStatus(ctx context.Context, userID int64, status AccountStatus) error {
	query := `UPDATE users SET account_status = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`
	now := time.Now()

	tag, err := r.db.Exec(ctx, query, status, now, userID)
	if err != nil {
		r.logger.Errorf("failed to update account status: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// LockAccount locks the account temporarily
func (r *IdentityRepositoryPostgres) LockAccount(ctx context.Context, userID int64, lockDuration time.Duration) error {
	query := `
		UPDATE users
		SET locked_until = $1, account_status = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	lockUntil := time.Now().Add(lockDuration)
	now := time.Now()

	tag, err := r.db.Exec(ctx, query, lockUntil, AccountStatusLocked, now, userID)
	if err != nil {
		r.logger.Errorf("failed to lock account: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	r.logger.Warnf("account locked for user %d until %v", userID, lockUntil)
	return nil
}

// UnlockAccount unlocks a locked account
func (r *IdentityRepositoryPostgres) UnlockAccount(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET locked_until = NULL, failed_login_attempts = 0, account_status = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	now := time.Now()

	tag, err := r.db.Exec(ctx, query, AccountStatusActive, now, userID)
	if err != nil {
		r.logger.Errorf("failed to unlock account: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	r.logger.Infof("account unlocked for user %d", userID)
	return nil
}

// RecordFailedLoginAttempt records a failed login attempt
func (r *IdentityRepositoryPostgres) RecordFailedLoginAttempt(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET failed_login_attempts = failed_login_attempts + 1,
			failed_login_last_attempt = $1,
			updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	now := time.Now()

	tag, err := r.db.Exec(ctx, query, now, now, userID)
	if err != nil {
		r.logger.Errorf("failed to record failed login attempt: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// RecordSuccessfulLogin records a successful login
func (r *IdentityRepositoryPostgres) RecordSuccessfulLogin(ctx context.Context, userID int64, ipAddress string) error {
	query := `
		UPDATE users
		SET last_login_at = $1,
			last_login_ip = $2,
			failed_login_attempts = 0,
			locked_until = NULL,
			updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	now := time.Now()

	tag, err := r.db.Exec(ctx, query, now, ipAddress, now, userID)
	if err != nil {
		r.logger.Errorf("failed to record successful login: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// UpdateOnboardingState updates the onboarding state
func (r *IdentityRepositoryPostgres) UpdateOnboardingState(ctx context.Context, userID int64, state OnboardingState) error {
	query := `
		UPDATE users
		SET onboarding_state = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	now := time.Now()

	tag, err := r.db.Exec(ctx, query, state, now, userID)
	if err != nil {
		r.logger.Errorf("failed to update onboarding state: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	r.logger.Infof("updated onboarding state for user %d to %s", userID, state)
	return nil
}

// UpdatePasswordHash updates the password hash
func (r *IdentityRepositoryPostgres) UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, password_changed_at = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	now := time.Now()

	tag, err := r.db.Exec(ctx, query, passwordHash, now, now, userID)
	if err != nil {
		r.logger.Errorf("failed to update password hash: %v", err)
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	r.logger.Infof("password updated for user %d", userID)
	return nil
}

// IdentityExists checks if an identity exists
func (r *IdentityRepositoryPostgres) IdentityExists(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		r.logger.Errorf("failed to check identity existence: %v", err)
		return false, err
	}

	return exists, nil
}

// EmailExists checks if an email is already registered
func (r *IdentityRepositoryPostgres) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		r.logger.Errorf("failed to check email existence: %v", err)
		return false, err
	}

	return exists, nil
}
