package auth

import (
	"context"
	"testing"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// TestCreateIdentity_Success tests successful identity creation
func TestCreateIdentity_Success(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}

	output, err := service.CreateIdentity(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Identity == nil {
		t.Fatal("expected identity, got nil")
	}
	if output.Identity.Email != "user@example.com" {
		t.Errorf("expected email user@example.com, got %s", output.Identity.Email)
	}
	if output.Identity.OnboardingState != StateNew {
		t.Errorf("expected onboarding state new, got %s", output.Identity.OnboardingState)
	}
	if output.Identity.AccountStatus != AccountStatusPending {
		t.Errorf("expected account status pending, got %s", output.Identity.AccountStatus)
	}
}

// TestCreateIdentity_DuplicateEmail tests creation with duplicate email
func TestCreateIdentity_DuplicateEmail(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create first identity
	input1 := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	service.CreateIdentity(context.Background(), input1)

	// Try to create second with same email
	input2 := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password456",
	}
	_, err := service.CreateIdentity(context.Background(), input2)

	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

// TestCreateIdentity_WeakPassword tests creation with weak password
func TestCreateIdentity_WeakPassword(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "weak", // Too short
	}

	_, err := service.CreateIdentity(context.Background(), input)

	if err == nil {
		t.Fatal("expected error for weak password")
	}
}

// TestGetIdentityByID retrieves identity by ID
func TestGetIdentityByID(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)
	userID := output.UserID

	// Retrieve
	identity, err := service.GetIdentityByID(context.Background(), userID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if identity.Email != "user@example.com" {
		t.Errorf("expected email user@example.com, got %s", identity.Email)
	}
}

// TestGetIdentityByEmail retrieves identity by email
func TestGetIdentityByEmail(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	service.CreateIdentity(context.Background(), input)

	// Retrieve by email
	identity, err := service.GetIdentityByEmail(context.Background(), "user@example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if identity.Email != "user@example.com" {
		t.Errorf("expected email user@example.com, got %s", identity.Email)
	}
}

// TestVerifyPassword_Success verifies correct password
func TestVerifyPassword_Success(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
		MaxLoginAttempts:       5,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Verify correct password
	verifyInput := &VerifyPasswordInput{
		UserID:        output.UserID,
		PlainPassword: "password123",
	}
	err := service.VerifyPassword(context.Background(), verifyInput)

	if err != nil {
		t.Fatalf("expected no error for correct password, got %v", err)
	}
}

// TestVerifyPassword_Failure verifies incorrect password
func TestVerifyPassword_Failure(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
		MaxLoginAttempts:       5,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Verify incorrect password
	verifyInput := &VerifyPasswordInput{
		UserID:        output.UserID,
		PlainPassword: "wrongpassword",
	}
	err := service.VerifyPassword(context.Background(), verifyInput)

	if err == nil {
		t.Fatal("expected error for incorrect password")
	}

	// Verify that failed attempt was recorded
	identity, _ := repo.GetIdentityByID(context.Background(), output.UserID)
	if identity.FailedLoginAttempts == 0 {
		t.Error("expected failed login attempt to be recorded")
	}
}

// TestRecordFailedLoginAttempt_AccountLocked tests account locking after too many attempts
func TestRecordFailedLoginAttempt_AccountLocked(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
		MaxLoginAttempts:       3,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Record failed attempts until account is locked
	for i := 0; i < int(config.MaxLoginAttempts); i++ {
		attemptInput := &RecordFailedLoginAttemptInput{
			UserID: output.UserID,
		}
		service.RecordFailedLoginAttempt(context.Background(), attemptInput)
	}

	// Check that account is locked
	identity, _ := repo.GetIdentityByID(context.Background(), output.UserID)
	if identity.LockedUntil == nil {
		t.Fatal("expected account to be locked")
	}
	if !time.Now().Before(*identity.LockedUntil) {
		t.Error("expected lock time to be in the future")
	}
}

// TestRecordSuccessfulLogin clears failed attempts
func TestRecordSuccessfulLogin(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
		MaxLoginAttempts:       5,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Record failed attempt
	attemptInput := &RecordFailedLoginAttemptInput{
		UserID: output.UserID,
	}
	service.RecordFailedLoginAttempt(context.Background(), attemptInput)

	// Record successful login
	successInput := &RecordSuccessfulLoginInput{
		UserID:    output.UserID,
		IPAddress: "192.168.1.1",
	}
	err := service.RecordSuccessfulLogin(context.Background(), successInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Check that failed attempts were cleared
	identity, _ := repo.GetIdentityByID(context.Background(), output.UserID)
	if identity.FailedLoginAttempts != 0 {
		t.Errorf("expected failed attempts to be cleared, got %d", identity.FailedLoginAttempts)
	}
	if identity.LastLoginAt == nil {
		t.Fatal("expected last login time to be set")
	}
}

// TestUpdateOnboardingState transitions state properly
func TestUpdateOnboardingState(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Update to email verified
	updateInput := &UpdateOnboardingStateInput{
		UserID: output.UserID,
		State:  StateEmailVerified,
	}
	err := service.UpdateOnboardingState(context.Background(), updateInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify state was updated
	identity, _ := repo.GetIdentityByID(context.Background(), output.UserID)
	if identity.OnboardingState != StateEmailVerified {
		t.Errorf("expected state email_verified, got %s", identity.OnboardingState)
	}
}

// TestUpdateOnboardingState_InvalidTransition tests invalid state transition
func TestUpdateOnboardingState_InvalidTransition(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Try invalid transition (new -> onboarded, skipping intermediate states)
	updateInput := &UpdateOnboardingStateInput{
		UserID: output.UserID,
		State:  StateOnboarded,
	}
	err := service.UpdateOnboardingState(context.Background(), updateInput)

	if err == nil {
		t.Fatal("expected error for invalid state transition")
	}
}

// TestLockAccount locks user account
func TestLockAccount(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Lock account
	lockInput := &LockAccountInput{
		UserID:       output.UserID,
		LockDuration: 1 * time.Hour,
		Reason:       "suspicious_activity",
	}
	err := service.LockAccount(context.Background(), lockInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify lock
	identity, _ := repo.GetIdentityByID(context.Background(), output.UserID)
	if identity.LockedUntil == nil {
		t.Fatal("expected account to be locked")
	}
	if !time.Now().Before(*identity.LockedUntil) {
		t.Error("expected lock time in the future")
	}
}

// TestUnlockAccount unlocks user account
func TestUnlockAccount(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create and lock identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	lockInput := &LockAccountInput{
		UserID:       output.UserID,
		LockDuration: 1 * time.Hour,
		Reason:       "test",
	}
	service.LockAccount(context.Background(), lockInput)

	// Unlock
	err := service.UnlockAccount(context.Background(), output.UserID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify unlock
	identity, _ := repo.GetIdentityByID(context.Background(), output.UserID)
	if identity.LockedUntil != nil {
		t.Fatal("expected account to be unlocked")
	}
	if identity.FailedLoginAttempts != 0 {
		t.Error("expected failed attempts to be cleared")
	}
}

// TestChangePassword changes user password
func TestChangePassword(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Change password
	changeInput := &ChangePasswordInput{
		UserID:      output.UserID,
		OldPassword: "password123",
		NewPassword: "newpassword456",
	}
	err := service.ChangePassword(context.Background(), changeInput)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify old password no longer works
	verifyInput := &VerifyPasswordInput{
		UserID:        output.UserID,
		PlainPassword: "password123",
	}
	err = service.VerifyPassword(context.Background(), verifyInput)
	if err == nil {
		t.Fatal("expected old password to not work")
	}

	// Verify new password works
	verifyInput.PlainPassword = "newpassword456"
	err = service.VerifyPassword(context.Background(), verifyInput)
	if err != nil {
		t.Fatalf("expected new password to work, got %v", err)
	}
}

// TestGetIdentityStatus returns proper status
func TestGetIdentityStatus(t *testing.T) {
	config := &AuthConfig{
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
	}
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)
	service := NewIdentityService(config, mockLogger, repo)

	// Create identity
	input := &CreateIdentityInput{
		Email:    "user@example.com",
		Password: "password123",
	}
	output, _ := service.CreateIdentity(context.Background(), input)

	// Get status
	status, err := service.GetIdentityStatus(context.Background(), output.UserID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status.IsActive {
		t.Error("expected new identity to not be active")
	}
	if status.IsLocked {
		t.Error("expected new identity to not be locked")
	}
	if status.CanLogin {
		t.Error("expected new identity to not be able to login yet")
	}
}

// TestIdentityRepository_InMemory tests repository directly
func TestIdentityRepository_InMemory(t *testing.T) {
	mockLogger := logger.NewLogger("test", false)
	repo := NewIdentityRepositoryMemory(mockLogger)

	identity := &UserIdentity{
		Email:           "user@example.com",
		PasswordHash:    "hashed_password",
		OnboardingState: StateNew,
		AccountStatus:   AccountStatusPending,
		IsActive:        false,
	}

	// Create
	err := repo.CreateIdentity(context.Background(), identity)
	if err != nil {
		t.Fatalf("expected no error creating identity, got %v", err)
	}

	// Get
	retrieved, err := repo.GetIdentityByEmail(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("expected no error getting identity, got %v", err)
	}
	if retrieved.Email != "user@example.com" {
		t.Errorf("expected email user@example.com, got %s", retrieved.Email)
	}

	// Update status
	err = repo.UpdateAccountStatus(context.Background(), retrieved.ID, AccountStatusActive)
	if err != nil {
		t.Fatalf("expected no error updating status, got %v", err)
	}

	// Verify update
	updated, _ := repo.GetIdentityByID(context.Background(), retrieved.ID)
	if updated.AccountStatus != AccountStatusActive {
		t.Errorf("expected status active, got %s", updated.AccountStatus)
	}
}
