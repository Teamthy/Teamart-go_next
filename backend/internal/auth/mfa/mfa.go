package mfa

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"time"
)

// MFAType represents the type of MFA
type MFAType string

const (
	MFATypeTOTP       MFAType = "totp"
	MFATypeBackupCode MFAType = "backup_code"
)

// MFASettings holds user MFA configuration
type MFASettings struct {
	UserID         int64
	MFAEnabled     bool
	TOTPSecret     string
	TOTPVerified   bool
	TOTPVerifiedAt *time.Time
	BackupCodes    []string // Encrypted
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MFAService manages Multi-Factor Authentication
type MFAService struct {
	storage MFAStorage
	config  *MFAConfig
}

// MFAConfig holds MFA configuration
type MFAConfig struct {
	TOTPIssuer       string
	TOTPAccountName  string
	BackupCodeCount  int
	BackupCodeLength int
	MaxAttempts      int32
	LockoutDuration  time.Duration
}

// MFAStorage defines storage interface for MFA settings
type MFAStorage interface {
	SaveMFASettings(ctx context.Context, settings *MFASettings) error
	GetMFASettings(ctx context.Context, userID int64) (*MFASettings, error)
	UpdateMFASettings(ctx context.Context, settings *MFASettings) error
	DeleteMFASettings(ctx context.Context, userID int64) error
	RecordMFAAttempt(ctx context.Context, userID int64, success bool, attemptType MFAType) error
	GetFailedAttempts(ctx context.Context, userID int64, duration time.Duration) (int32, error)
}

// NewMFAService creates a new MFA service
func NewMFAService(storage MFAStorage, config *MFAConfig) *MFAService {
	if config == nil {
		config = &MFAConfig{
			TOTPIssuer:       "Teamart",
			TOTPAccountName:  "Teamart Account",
			BackupCodeCount:  10,
			BackupCodeLength: 8,
			MaxAttempts:      5,
			LockoutDuration:  30 * time.Minute,
		}
	}

	return &MFAService{
		storage: storage,
		config:  config,
	}
}

// GenerateTOTPSecret generates a new TOTP secret for a user
func (s *MFAService) GenerateTOTPSecret(ctx context.Context, userID int64) (string, error) {
	// Check if user already has TOTP enabled
	settings, err := s.storage.GetMFASettings(ctx, userID)
	if err == nil && settings != nil && settings.TOTPVerified {
		return "", errors.New("user already has TOTP enabled")
	}

	// Generate random secret
	secret, err := s.generateRandomSecret(32)
	if err != nil {
		return "", err
	}

	return secret, nil
}

// VerifyTOTPSecret verifies a TOTP secret setup
func (s *MFAService) VerifyTOTPSecret(ctx context.Context, userID int64, secret string, code string) error {
	if secret == "" {
		return errors.New("secret is required")
	}
	if code == "" {
		return errors.New("code is required")
	}

	// Verify the code matches the secret
	if !s.verifyTOTPCode(secret, code) {
		_ = s.storage.RecordMFAAttempt(ctx, userID, false, MFATypeTOTP)
		return errors.New("invalid code")
	}

	// Save MFA settings
	now := time.Now()
	settings := &MFASettings{
		UserID:         userID,
		MFAEnabled:     true,
		TOTPSecret:     secret,
		TOTPVerified:   true,
		TOTPVerifiedAt: &now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.storage.SaveMFASettings(ctx, settings); err != nil {
		return fmt.Errorf("failed to save MFA settings: %w", err)
	}

	_ = s.storage.RecordMFAAttempt(ctx, userID, true, MFATypeTOTP)

	return nil
}

// VerifyTOTPCode verifies a TOTP code for login
func (s *MFAService) VerifyTOTPCode(ctx context.Context, userID int64, code string) error {
	if code == "" {
		return errors.New("code is required")
	}

	// Get user's MFA settings
	settings, err := s.storage.GetMFASettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get MFA settings: %w", err)
	}

	if settings == nil || !settings.MFAEnabled {
		return errors.New("MFA not enabled for user")
	}

	// Check failed attempts
	failedAttempts, err := s.storage.GetFailedAttempts(ctx, userID, s.config.LockoutDuration)
	if err != nil {
		return fmt.Errorf("failed to get failed attempts: %w", err)
	}

	if failedAttempts >= s.config.MaxAttempts {
		return errors.New("too many failed attempts - account locked")
	}

	// Verify the code
	if !s.verifyTOTPCode(settings.TOTPSecret, code) {
		_ = s.storage.RecordMFAAttempt(ctx, userID, false, MFATypeTOTP)
		return errors.New("invalid code")
	}

	_ = s.storage.RecordMFAAttempt(ctx, userID, true, MFATypeTOTP)
	return nil
}

// DisableMFA disables MFA for a user
func (s *MFAService) DisableMFA(ctx context.Context, userID int64) error {
	return s.storage.DeleteMFASettings(ctx, userID)
}

// IsMFAEnabled checks if MFA is enabled for a user
func (s *MFAService) IsMFAEnabled(ctx context.Context, userID int64) (bool, error) {
	settings, err := s.storage.GetMFASettings(ctx, userID)
	if err != nil {
		return false, err
	}

	if settings == nil {
		return false, nil
	}

	return settings.MFAEnabled, nil
}

// GetMFASettings retrieves MFA settings for a user
func (s *MFAService) GetMFASettings(ctx context.Context, userID int64) (*MFASettings, error) {
	return s.storage.GetMFASettings(ctx, userID)
}

// verifyTOTPCode verifies a TOTP code against a secret
func (s *MFAService) verifyTOTPCode(secret, code string) bool {
	// Decode the secret from base32
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return false
	}

	// Get current time counter
	counter := time.Now().Unix() / 30

	// Check current and adjacent time windows (for clock skew tolerance)
	for i := -1; i <= 1; i++ {
		hash := s.generateHOTP(key, counter+int64(i))
		if hash == code {
			return true
		}
	}

	return false
}

// generateHOTP generates HOTP value
func (s *MFAService) generateHOTP(secret []byte, counter int64) string {
	// RFC 4226 HMAC-based One-Time Password implementation
	// This is a simplified version - production should use a proper TOTP library

	hash := sha1.New()
	// In production, use HMAC-SHA1
	// For now, this is a placeholder
	return fmt.Sprintf("%06d", int(hash.Sum(nil)[0])%1000000)
}

// generateRandomSecret generates a random TOTP secret
func (s *MFAService) generateRandomSecret(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode in base32 for TOTP compatibility
	return base32.StdEncoding.EncodeToString(randomBytes), nil
}
