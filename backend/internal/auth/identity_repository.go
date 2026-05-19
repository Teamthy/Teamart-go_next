package auth

import (
	"context"
	"time"
)

// IdentityRepository defines the interface for identity persistence
type IdentityRepository interface {
	// CreateIdentity creates a new user identity
	CreateIdentity(ctx context.Context, identity *UserIdentity) error

	// GetIdentityByID retrieves a user identity by ID
	GetIdentityByID(ctx context.Context, userID int64) (*UserIdentity, error)

	// GetIdentityByEmail retrieves a user identity by email
	GetIdentityByEmail(ctx context.Context, email string) (*UserIdentity, error)

	// UpdateIdentity updates an existing user identity
	UpdateIdentity(ctx context.Context, identity *UserIdentity) error

	// UpdateAccountStatus updates account status
	UpdateAccountStatus(ctx context.Context, userID int64, status AccountStatus) error

	// LockAccount locks the account temporarily
	LockAccount(ctx context.Context, userID int64, lockDuration time.Duration) error

	// UnlockAccount unlocks a locked account
	UnlockAccount(ctx context.Context, userID int64) error

	// RecordFailedLoginAttempt records a failed login attempt
	RecordFailedLoginAttempt(ctx context.Context, userID int64) error

	// RecordSuccessfulLogin records a successful login
	RecordSuccessfulLogin(ctx context.Context, userID int64, ipAddress string) error

	// UpdateOnboardingState updates the onboarding state
	UpdateOnboardingState(ctx context.Context, userID int64, state OnboardingState) error

	// UpdatePasswordHash updates the password hash
	UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string) error

	// IdentityExists checks if an identity exists
	IdentityExists(ctx context.Context, userID int64) (bool, error)

	// EmailExists checks if an email is already registered
	EmailExists(ctx context.Context, email string) (bool, error)
}
