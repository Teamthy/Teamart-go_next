package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// IdentityService manages user identity and authentication
type IdentityService struct {
	config *AuthConfig
	logger *logger.Logger
	repo   IdentityRepository
}

// NewIdentityService creates a new identity service
func NewIdentityService(config *AuthConfig, logger *logger.Logger, repo IdentityRepository) *IdentityService {
	return &IdentityService{
		config: config,
		logger: logger,
		repo:   repo,
	}
}

// CreateIdentityInput represents input for creating a new identity
type CreateIdentityInput struct {
	Email    string
	Password string // Plain text password (should be hashed before storing)
}

// CreateIdentityOutput represents the output of creating an identity
type CreateIdentityOutput struct {
	Identity *UserIdentity
	UserID   int64
}

// CreateIdentity creates a new user identity
func (is *IdentityService) CreateIdentity(ctx context.Context, input *CreateIdentityInput) (*CreateIdentityOutput, error) {
	// Validate inputs
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if input.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Validate password strength
	if err := is.validatePasswordStrength(input.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := is.repo.EmailExists(ctx, input.Email)
	if err != nil {
		is.logger.Errorf("failed to check email existence: %v", err)
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already registered: %s", input.Email)
	}

	// Hash password (in production, use bcrypt)
	passwordHash := is.hashPassword(input.Password)

	// Create identity
	identity := &UserIdentity{
		Email:           input.Email,
		PasswordHash:    passwordHash,
		OnboardingState: StateNew,
		AccountStatus:   AccountStatusPending,
		IsActive:        false,
	}

	// Persist to repository
	if err := is.repo.CreateIdentity(ctx, identity); err != nil {
		is.logger.Errorf("failed to create identity: %v", err)
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	is.logger.Infof("created new identity with email %s (user ID: %d)", input.Email, identity.ID)

	return &CreateIdentityOutput{
		Identity: identity,
		UserID:   identity.ID,
	}, nil
}

// GetIdentityInput represents input for retrieving an identity
type GetIdentityInput struct {
	UserID *int64 // Either UserID or Email must be provided
	Email  *string
}

// GetIdentityByID retrieves an identity by user ID
func (is *IdentityService) GetIdentityByID(ctx context.Context, userID int64) (*UserIdentity, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	identity, err := is.repo.GetIdentityByID(ctx, userID)
	if err != nil {
		is.logger.Debugf("identity not found: %d", userID)
		return nil, fmt.Errorf("identity not found: %w", err)
	}

	return identity, nil
}

// GetIdentityByEmail retrieves an identity by email
func (is *IdentityService) GetIdentityByEmail(ctx context.Context, email string) (*UserIdentity, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	identity, err := is.repo.GetIdentityByEmail(ctx, email)
	if err != nil {
		is.logger.Debugf("identity not found for email: %s", email)
		return nil, fmt.Errorf("identity not found: %w", err)
	}

	return identity, nil
}

// VerifyPasswordInput represents input for password verification
type VerifyPasswordInput struct {
	UserID        int64
	PlainPassword string
}

// VerifyPassword verifies a password against stored hash
func (is *IdentityService) VerifyPassword(ctx context.Context, input *VerifyPasswordInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if input.PlainPassword == "" {
		return fmt.Errorf("password is required")
	}

	// Get identity
	identity, err := is.repo.GetIdentityByID(ctx, input.UserID)
	if err != nil {
		is.logger.Warnf("failed to verify password for user %d: %v", input.UserID, err)
		return fmt.Errorf("failed to verify password: %w", err)
	}

	// Verify password (in production, use bcrypt.CompareHashAndPassword)
	if !is.verifyPassword(input.PlainPassword, identity.PasswordHash) {
		is.logger.Warnf("password verification failed for user %d", input.UserID)

		// Record failed attempt
		if err := is.repo.RecordFailedLoginAttempt(ctx, input.UserID); err != nil {
			is.logger.Warnf("failed to record login attempt: %v", err)
		}

		return fmt.Errorf("invalid password")
	}

	return nil
}

// UpdateOnboardingStateInput represents input for updating onboarding state
type UpdateOnboardingStateInput struct {
	UserID int64
	State  OnboardingState
}

// UpdateOnboardingState updates the user's onboarding state
func (is *IdentityService) UpdateOnboardingState(ctx context.Context, input *UpdateOnboardingStateInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	// Validate state transition
	identity, err := is.repo.GetIdentityByID(ctx, input.UserID)
	if err != nil {
		return fmt.Errorf("failed to get identity: %w", err)
	}

	if !ValidStateTransition(identity.OnboardingState, input.State) {
		return fmt.Errorf("invalid state transition from %s to %s", identity.OnboardingState, input.State)
	}

	// Update state
	if err := is.repo.UpdateOnboardingState(ctx, input.UserID, input.State); err != nil {
		is.logger.Errorf("failed to update onboarding state: %v", err)
		return fmt.Errorf("failed to update onboarding state: %w", err)
	}

	is.logger.Infof("updated onboarding state for user %d to %s", input.UserID, input.State)
	return nil
}

// LockAccountInput represents input for locking an account
type LockAccountInput struct {
	UserID       int64
	LockDuration time.Duration
	Reason       string
}

// LockAccount locks a user account
func (is *IdentityService) LockAccount(ctx context.Context, input *LockAccountInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	lockDuration := input.LockDuration
	if lockDuration == 0 {
		lockDuration = 15 * time.Minute // Default lock duration
	}

	if err := is.repo.LockAccount(ctx, input.UserID, lockDuration); err != nil {
		is.logger.Errorf("failed to lock account: %v", err)
		return fmt.Errorf("failed to lock account: %w", err)
	}

	is.logger.Warnf("locked account for user %d for %v (reason: %s)",
		input.UserID, lockDuration, input.Reason)
	return nil
}

// UnlockAccount unlocks a locked account
func (is *IdentityService) UnlockAccount(ctx context.Context, userID int64) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	if err := is.repo.UnlockAccount(ctx, userID); err != nil {
		is.logger.Errorf("failed to unlock account: %v", err)
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	is.logger.Infof("unlocked account for user %d", userID)
	return nil
}

// RecordFailedLoginAttemptInput represents input for recording a failed login
type RecordFailedLoginAttemptInput struct {
	UserID int64
}

// RecordFailedLoginAttempt records a failed login attempt and locks account if needed
func (is *IdentityService) RecordFailedLoginAttempt(ctx context.Context, input *RecordFailedLoginAttemptInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	// Record attempt
	if err := is.repo.RecordFailedLoginAttempt(ctx, input.UserID); err != nil {
		is.logger.Errorf("failed to record login attempt: %v", err)
		return fmt.Errorf("failed to record login attempt: %w", err)
	}

	// Check if account should be locked
	identity, err := is.repo.GetIdentityByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if identity.FailedLoginAttempts >= is.config.MaxLoginAttempts {
		lockInput := &LockAccountInput{
			UserID:       input.UserID,
			LockDuration: 15 * time.Minute,
			Reason:       "too_many_failed_attempts",
		}
		return is.LockAccount(ctx, lockInput)
	}

	is.logger.Debugf("recorded failed login attempt for user %d (count: %d)",
		input.UserID, identity.FailedLoginAttempts)
	return nil
}

// RecordSuccessfulLoginInput represents input for recording a successful login
type RecordSuccessfulLoginInput struct {
	UserID    int64
	IPAddress string
}

// RecordSuccessfulLogin records a successful login
func (is *IdentityService) RecordSuccessfulLogin(ctx context.Context, input *RecordSuccessfulLoginInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if input.IPAddress == "" {
		return fmt.Errorf("IP address is required")
	}

	if err := is.repo.RecordSuccessfulLogin(ctx, input.UserID, input.IPAddress); err != nil {
		is.logger.Errorf("failed to record successful login: %v", err)
		return fmt.Errorf("failed to record successful login: %w", err)
	}

	is.logger.Infof("recorded successful login for user %d from IP %s", input.UserID, input.IPAddress)
	return nil
}

// ChangePasswordInput represents input for changing password
type ChangePasswordInput struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

// ChangePassword changes a user's password
func (is *IdentityService) ChangePassword(ctx context.Context, input *ChangePasswordInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if input.OldPassword == "" {
		return fmt.Errorf("old password is required")
	}
	if input.NewPassword == "" {
		return fmt.Errorf("new password is required")
	}

	// Verify old password
	verifyInput := &VerifyPasswordInput{
		UserID:        input.UserID,
		PlainPassword: input.OldPassword,
	}
	if err := is.VerifyPassword(ctx, verifyInput); err != nil {
		is.logger.Warnf("password change failed for user %d: invalid old password", input.UserID)
		return fmt.Errorf("invalid old password")
	}

	// Validate new password strength
	if err := is.validatePasswordStrength(input.NewPassword); err != nil {
		return err
	}

	// Hash new password
	newPasswordHash := is.hashPassword(input.NewPassword)

	// Update password
	if err := is.repo.UpdatePasswordHash(ctx, input.UserID, newPasswordHash); err != nil {
		is.logger.Errorf("failed to update password: %v", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	is.logger.Infof("password changed for user %d", input.UserID)
	return nil
}

// GetIdentityStatusInput represents input for checking identity status
type GetIdentityStatusInput struct {
	UserID int64
}

// GetIdentityStatus returns detailed status of an identity
type GetIdentityStatus struct {
	Identity         *UserIdentity
	IsActive         bool
	IsLocked         bool
	RequiresPassword bool
	CanLogin         bool
}

// GetIdentityStatus checks the status of an identity
func (is *IdentityService) GetIdentityStatus(ctx context.Context, userID int64) (*GetIdentityStatus, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	identity, err := is.repo.GetIdentityByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("identity not found: %w", err)
	}

	status := &GetIdentityStatus{
		Identity:         identity,
		IsActive:         identity.IsActive && identity.AccountStatus == AccountStatusActive,
		IsLocked:         identity.LockedUntil != nil && time.Now().Before(*identity.LockedUntil),
		RequiresPassword: identity.PasswordHash == "",
		CanLogin:         identity.CanLogin(),
	}

	return status, nil
}

// ===== Helper Methods =====

// validatePasswordStrength validates password meets strength requirements
func (is *IdentityService) validatePasswordStrength(password string) error {
	if len(password) < is.config.PasswordMinLength {
		return fmt.Errorf("password must be at least %d characters", is.config.PasswordMinLength)
	}

	if is.config.PasswordRequireSpecial {
		hasSpecial := false
		for _, char := range password {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
				hasSpecial = true
				break
			}
		}
		if !hasSpecial {
			return fmt.Errorf("password must contain special characters")
		}
	}

	if is.config.PasswordRequireNumbers {
		hasNumber := false
		for _, char := range password {
			if char >= '0' && char <= '9' {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return fmt.Errorf("password must contain numbers")
		}
	}

	return nil
}

// hashPassword hashes a password (simple implementation - use bcrypt in production)
func (is *IdentityService) hashPassword(password string) string {
	// In production, use: bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// This is a simple implementation for testing purposes
	return fmt.Sprintf("hashed_%s", password)
}

// verifyPassword verifies a password against a hash (simple implementation)
func (is *IdentityService) verifyPassword(password, hash string) bool {
	// In production, use: bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// This is a simple implementation for testing purposes
	return hash == fmt.Sprintf("hashed_%s", password)
}
