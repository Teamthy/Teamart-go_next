package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// OTPService manages OTP generation, verification, and expiration
type OTPService struct {
	config *AuthConfig
	logger *logger.Logger
}

// NewOTPService creates a new OTP service
func NewOTPService(config *AuthConfig, logger *logger.Logger) *OTPService {
	return &OTPService{
		config: config,
		logger: logger,
	}
}

// GenerateOTPInput represents input for OTP generation
type GenerateOTPInput struct {
	UserID      int64
	Type        OTPType
	Destination string // Email or phone number
}

// GenerateOTPOutput represents the generated OTP
type GenerateOTPOutput struct {
	ID          string
	UserID      int64
	Type        OTPType
	CodeHash    string // Hashed OTP code
	PlainCode   string // Actual code to send (only shown once)
	Destination string
	ExpiresAt   time.Time
}

// GenerateOTP generates a new OTP code
func (s *OTPService) GenerateOTP(ctx context.Context, input *GenerateOTPInput) (*GenerateOTPOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.Destination == "" {
		return nil, fmt.Errorf("destination is required")
	}

	// Generate random OTP code
	code, err := s.generateRandomCode(s.config.OTPLength)
	if err != nil {
		s.logger.Errorf("failed to generate OTP code: %v", err)
		return nil, err
	}

	// Hash the code for storage
	codeHash := s.hashOTP(code)
	
	// Generate unique ID
	id := s.generateRandomID()
	
	expiresAt := time.Now().Add(s.config.OTPTTL)

	s.logger.Infof("OTP generated for user %d via %s, expires at %v", 
		input.UserID, input.Type, expiresAt)

	return &GenerateOTPOutput{
		ID:          id,
		UserID:      input.UserID,
		Type:        input.Type,
		CodeHash:    codeHash,
		PlainCode:   code, // Only shown once
		Destination: input.Destination,
		ExpiresAt:   expiresAt,
	}, nil
}

// VerifyOTPInput represents input for OTP verification
type VerifyOTPInput struct {
	UserID int64
	OTPID  string
	Code   string // Plain OTP code from user
}

// VerifyOTPOutput represents the result of OTP verification
type VerifyOTPOutput struct {
	IsValid   bool
	UserID    int64
	Type      OTPType
	VerifiedAt time.Time
}

// VerifyOTP verifies an OTP code
func (s *OTPService) VerifyOTP(ctx context.Context, input *VerifyOTPInput) (*VerifyOTPOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.OTPID == "" {
		return nil, fmt.Errorf("OTP ID is required")
	}
	if input.Code == "" {
		return nil, fmt.Errorf("OTP code is required")
	}

	s.logger.Debugf("verifying OTP for user %d", input.UserID)

	// In a real implementation, this would:
	// 1. Fetch OTP from database by ID
	// 2. Check if OTP is expired
	// 3. Check if max attempts exceeded
	// 4. Compare provided code hash with stored hash
	// 5. Mark as verified if correct
	// 6. Return error if incorrect

	codeHash := s.hashOTP(input.Code)

	// Simulate verification
	return &VerifyOTPOutput{
		IsValid:    true,
		UserID:     input.UserID,
		VerifiedAt: time.Now(),
	}, nil
}

// RevokeOTPInput represents input for OTP revocation
type RevokeOTPInput struct {
	UserID int64
	OTPID  string
	Reason string
}

// RevokeOTP revokes an OTP code
func (s *OTPService) RevokeOTP(ctx context.Context, input *RevokeOTPInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if input.OTPID == "" {
		return fmt.Errorf("OTP ID is required")
	}

	s.logger.Infof("revoking OTP %s for user %d (reason: %s)", 
		input.OTPID, input.UserID, input.Reason)

	// In a real implementation, this would mark the OTP as revoked in the database

	return nil
}

// ValidateOTPAttempt checks if user has exceeded max OTP attempts
func (s *OTPService) ValidateOTPAttempt(ctx context.Context, userID int64, currentAttempts int32) error {
	if currentAttempts >= s.config.OTPMaxAttempts {
		s.logger.Warnf("user %d exceeded max OTP attempts (%d)", userID, s.config.OTPMaxAttempts)
		return ErrOTPMaxAttempts
	}
	return nil
}

// ===== Helper Methods =====

// generateRandomCode generates a random OTP code
func (s *OTPService) generateRandomCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)
	
	for i := 0; i < length; i++ {
		b := make([]byte, 1)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		code[i] = digits[int(b[0])%len(digits)]
	}
	
	return string(code), nil
}

// hashOTP hashes an OTP code using SHA256
func (s *OTPService) hashOTP(code string) string {
	hash := sha256.Sum256([]byte(code))
	return fmt.Sprintf("%x", hash)
}

// generateRandomID generates a random ID
func (s *OTPService) generateRandomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// ===== OTP Configuration Methods =====

// IsOTPExpired checks if OTP has expired
func (s *OTPService) IsOTPExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// GetTimeRemainingForOTP returns remaining time for OTP
func (s *OTPService) GetTimeRemainingForOTP(expiresAt time.Time) time.Duration {
	remaining := time.Until(expiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
