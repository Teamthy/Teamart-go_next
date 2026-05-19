package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// IdentityRepositoryMemory implements IdentityRepository using in-memory storage
type IdentityRepositoryMemory struct {
	mu         sync.RWMutex
	identities map[int64]*UserIdentity // By ID
	emails     map[string]int64        // Email -> ID mapping
	logger     *logger.Logger
	nextID     int64
}

// NewIdentityRepositoryMemory creates a new in-memory identity repository
func NewIdentityRepositoryMemory(logger *logger.Logger) *IdentityRepositoryMemory {
	return &IdentityRepositoryMemory{
		identities: make(map[int64]*UserIdentity),
		emails:     make(map[string]int64),
		logger:     logger,
		nextID:     1,
	}
}

// CreateIdentity creates a new user identity
func (r *IdentityRepositoryMemory) CreateIdentity(ctx context.Context, identity *UserIdentity) error {
	if identity == nil {
		return fmt.Errorf("identity cannot be nil")
	}
	if identity.Email == "" {
		return fmt.Errorf("email is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if email already exists
	if _, exists := r.emails[identity.Email]; exists {
		return fmt.Errorf("email already registered: %s", identity.Email)
	}

	// Assign ID if not already set
	if identity.ID == 0 {
		identity.ID = r.nextID
		r.nextID++
	}

	// Check if ID is already taken
	if _, exists := r.identities[identity.ID]; exists {
		return fmt.Errorf("identity with ID %d already exists", identity.ID)
	}

	// Set timestamps
	now := time.Now()
	identity.CreatedAt = now
	identity.UpdatedAt = now
	identity.AccountStatus = AccountStatusPending
	identity.IsActive = false

	// Store identity
	identityCopy := *identity
	r.identities[identity.ID] = &identityCopy
	r.emails[identity.Email] = identity.ID

	r.logger.Debugf("created identity %d with email %s", identity.ID, identity.Email)
	return nil
}

// GetIdentityByID retrieves a user identity by ID
func (r *IdentityRepositoryMemory) GetIdentityByID(ctx context.Context, userID int64) (*UserIdentity, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	identity, exists := r.identities[userID]
	if !exists {
		return nil, fmt.Errorf("identity not found: %d", userID)
	}

	// Return a copy
	identityCopy := *identity
	return &identityCopy, nil
}

// GetIdentityByEmail retrieves a user identity by email
func (r *IdentityRepositoryMemory) GetIdentityByEmail(ctx context.Context, email string) (*UserIdentity, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.emails[email]
	if !exists {
		return nil, fmt.Errorf("identity not found for email: %s", email)
	}

	identity, exists := r.identities[userID]
	if !exists {
		return nil, fmt.Errorf("identity not found: %d", userID)
	}

	// Return a copy
	identityCopy := *identity
	return &identityCopy, nil
}

// UpdateIdentity updates an existing user identity
func (r *IdentityRepositoryMemory) UpdateIdentity(ctx context.Context, identity *UserIdentity) error {
	if identity == nil {
		return fmt.Errorf("identity cannot be nil")
	}
	if identity.ID == 0 {
		return fmt.Errorf("identity ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.identities[identity.ID]
	if !exists {
		return fmt.Errorf("identity not found: %d", identity.ID)
	}

	// If email is being changed, update mapping
	if existing.Email != identity.Email {
		if _, emailExists := r.emails[identity.Email]; emailExists {
			return fmt.Errorf("email already in use: %s", identity.Email)
		}
		delete(r.emails, existing.Email)
		r.emails[identity.Email] = identity.ID
	}

	// Update timestamp
	identity.UpdatedAt = time.Now()

	// Store updated identity
	identityCopy := *identity
	r.identities[identity.ID] = &identityCopy

	r.logger.Debugf("updated identity %d", identity.ID)
	return nil
}

// UpdateAccountStatus updates account status
func (r *IdentityRepositoryMemory) UpdateAccountStatus(ctx context.Context, userID int64, status AccountStatus) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	identity, exists := r.identities[userID]
	if !exists {
		return fmt.Errorf("identity not found: %d", userID)
	}

	identity.AccountStatus = status
	identity.UpdatedAt = time.Now()

	// Update IsActive based on status
	identity.IsActive = (status == AccountStatusActive)

	r.logger.Debugf("updated account status for user %d to %s", userID, status)
	return nil
}

// LockAccount locks the account temporarily
func (r *IdentityRepositoryMemory) LockAccount(ctx context.Context, userID int64, lockDuration time.Duration) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	identity, exists := r.identities[userID]
	if !exists {
		return fmt.Errorf("identity not found: %d", userID)
	}

	lockUntil := time.Now().Add(lockDuration)
	identity.LockedUntil = &lockUntil
	identity.UpdatedAt = time.Now()

	r.logger.Warnf("locked account for user %d until %v", userID, lockUntil)
	return nil
}

// UnlockAccount unlocks a locked account
func (r *IdentityRepositoryMemory) UnlockAccount(ctx context.Context, userID int64) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	identity, exists := r.identities[userID]
	if !exists {
		return fmt.Errorf("identity not found: %d", userID)
	}

	identity.LockedUntil = nil
	identity.FailedLoginAttempts = 0
	identity.UpdatedAt = time.Now()

	r.logger.Infof("unlocked account for user %d", userID)
	return nil
}

// RecordFailedLoginAttempt records a failed login attempt
func (r *IdentityRepositoryMemory) RecordFailedLoginAttempt(ctx context.Context, userID int64) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	identity, exists := r.identities[userID]
	if !exists {
		return fmt.Errorf("identity not found: %d", userID)
	}

	now := time.Now()
	identity.FailedLoginAttempts++
	identity.FailedLoginLastAttempt = &now
	identity.UpdatedAt = now

	r.logger.Debugf("recorded failed login attempt for user %d (count: %d)", userID, identity.FailedLoginAttempts)
	return nil
}

// RecordSuccessfulLogin records a successful login
func (r *IdentityRepositoryMemory) RecordSuccessfulLogin(ctx context.Context, userID int64, ipAddress string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	identity, exists := r.identities[userID]
	if !exists {
		return fmt.Errorf("identity not found: %d", userID)
	}

	now := time.Now()
	identity.FailedLoginAttempts = 0
	identity.FailedLoginLastAttempt = nil
	identity.LastLoginAt = &now
	identity.LastLoginIP = ipAddress
	identity.UpdatedAt = now

	r.logger.Infof("recorded successful login for user %d from IP %s", userID, ipAddress)
	return nil
}

// UpdateOnboardingState updates the onboarding state
func (r *IdentityRepositoryMemory) UpdateOnboardingState(ctx context.Context, userID int64, state OnboardingState) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	identity, exists := r.identities[userID]
	if !exists {
		return fmt.Errorf("identity not found: %d", userID)
	}

	// Validate transition
	if !ValidStateTransition(identity.OnboardingState, state) {
		return fmt.Errorf("invalid state transition from %s to %s", identity.OnboardingState, state)
	}

	identity.OnboardingState = state
	identity.UpdatedAt = time.Now()

	// Update account status based on onboarding state
	if state == StateOnboarded {
		identity.AccountStatus = AccountStatusActive
		identity.IsActive = true
	}

	r.logger.Infof("updated onboarding state for user %d to %s", userID, state)
	return nil
}

// UpdatePasswordHash updates the password hash
func (r *IdentityRepositoryMemory) UpdatePasswordHash(ctx context.Context, userID int64, passwordHash string) error {
	if userID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if passwordHash == "" {
		return fmt.Errorf("password hash is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	identity, exists := r.identities[userID]
	if !exists {
		return fmt.Errorf("identity not found: %d", userID)
	}

	identity.PasswordHash = passwordHash
	identity.PasswordChangedAt = time.Now()
	identity.UpdatedAt = time.Now()

	r.logger.Debugf("updated password hash for user %d", userID)
	return nil
}

// IdentityExists checks if an identity exists
func (r *IdentityRepositoryMemory) IdentityExists(ctx context.Context, userID int64) (bool, error) {
	if userID == 0 {
		return false, fmt.Errorf("user ID is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.identities[userID]
	return exists, nil
}

// EmailExists checks if an email is already registered
func (r *IdentityRepositoryMemory) EmailExists(ctx context.Context, email string) (bool, error) {
	if email == "" {
		return false, fmt.Errorf("email is required")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.emails[email]
	return exists, nil
}
