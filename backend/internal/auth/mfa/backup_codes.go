package mfa

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
)

// BackupCodesService manages MFA backup codes
type BackupCodesService struct {
	storage MFAStorage
	config  *BackupCodesConfig
}

// BackupCodesConfig holds backup codes configuration
type BackupCodesConfig struct {
	CodeCount  int // Default: 10
	CodeLength int // Default: 8
	GroupSize  int // Default: 4 (for display formatting like XXXX-XXXX)
}

// BackupCode represents a single backup code
type BackupCode struct {
	Code      string
	Used      bool
	UsedAt    *string
	CreatedAt string
}

// NewBackupCodesService creates a new backup codes service
func NewBackupCodesService(storage MFAStorage, config *BackupCodesConfig) *BackupCodesService {
	if config == nil {
		config = &BackupCodesConfig{
			CodeCount:  10,
			CodeLength: 8,
			GroupSize:  4,
		}
	}

	return &BackupCodesService{
		storage: storage,
		config:  config,
	}
}

// GenerateBackupCodes generates new backup codes for a user
func (s *BackupCodesService) GenerateBackupCodes(ctx context.Context, userID int64) ([]string, error) {
	codes := make([]string, s.config.CodeCount)

	for i := 0; i < s.config.CodeCount; i++ {
		code, err := s.generateRandomCode(s.config.CodeLength)
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}
		codes[i] = code
	}

	// Update MFA settings with new backup codes (encrypted)
	settings, err := s.storage.GetMFASettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get MFA settings: %w", err)
	}

	if settings == nil {
		return nil, fmt.Errorf("user MFA settings not found")
	}

	settings.BackupCodes = codes
	if err := s.storage.UpdateMFASettings(ctx, settings); err != nil {
		return nil, fmt.Errorf("failed to save backup codes: %w", err)
	}

	return codes, nil
}

// VerifyBackupCode verifies and consumes a backup code
func (s *BackupCodesService) VerifyBackupCode(ctx context.Context, userID int64, code string) error {
	if code == "" {
		return fmt.Errorf("code is required")
	}

	// Get user's MFA settings
	settings, err := s.storage.GetMFASettings(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get MFA settings: %w", err)
	}

	if settings == nil {
		return fmt.Errorf("user MFA settings not found")
	}

	// Normalize code (remove hyphens, convert to uppercase)
	normalizedCode := strings.ToUpper(strings.ReplaceAll(code, "-", ""))

	// Find and consume the code
	codeFound := false
	for i, backupCode := range settings.BackupCodes {
		normalizedBackupCode := strings.ToUpper(strings.ReplaceAll(backupCode, "-", ""))
		if normalizedBackupCode == normalizedCode {
			codeFound = true
			// Mark code as used (this should be handled separately in production)
			// For now, we remove the code from the list
			settings.BackupCodes = append(settings.BackupCodes[:i], settings.BackupCodes[i+1:]...)
			break
		}
	}

	if !codeFound {
		_ = s.storage.RecordMFAAttempt(ctx, userID, false, MFATypeBackupCode)
		return fmt.Errorf("invalid backup code")
	}

	// Update settings with consumed code
	if err := s.storage.UpdateMFASettings(ctx, settings); err != nil {
		return fmt.Errorf("failed to update MFA settings: %w", err)
	}

	_ = s.storage.RecordMFAAttempt(ctx, userID, true, MFATypeBackupCode)

	return nil
}

// GetBackupCodeCount returns the number of unused backup codes
func (s *BackupCodesService) GetBackupCodeCount(ctx context.Context, userID int64) (int, error) {
	settings, err := s.storage.GetMFASettings(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get MFA settings: %w", err)
	}

	if settings == nil {
		return 0, nil
	}

	return len(settings.BackupCodes), nil
}

// HasBackupCodes checks if user has any backup codes remaining
func (s *BackupCodesService) HasBackupCodes(ctx context.Context, userID int64) (bool, error) {
	count, err := s.GetBackupCodeCount(ctx, userID)
	return count > 0, err
}

// generateRandomCode generates a random backup code
func (s *BackupCodesService) generateRandomCode(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert to alphanumeric code (case-insensitive for user convenience)
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := 0; i < length; i++ {
		code[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	// Format as XXXX-XXXX or similar
	formatted := s.formatBackupCode(string(code))
	return formatted, nil
}

// formatBackupCode formats a backup code for display
func (s *BackupCodesService) formatBackupCode(code string) string {
	var formatted strings.Builder

	for i, char := range code {
		if i > 0 && i%s.config.GroupSize == 0 {
			formatted.WriteRune('-')
		}
		formatted.WriteRune(char)
	}

	return formatted.String()
}
