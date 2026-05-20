package auth

import (
	"time"
)

// AuthConfigManager manages authentication configuration
type AuthConfigManager struct {
	// JWT Configuration
	JWTSecret           string
	JWTAccessTokenTTL   time.Duration
	JWTRefreshTokenTTL  time.Duration
	JWTEmailTokenTTL    time.Duration
	JWTPasswordResetTTL time.Duration
	JWTIssuer           string
	JWTAudience         string

	// OTP Configuration
	OTPLength      int
	OTPTTL         time.Duration
	OTPMaxAttempts int32
	OTPResendDelay time.Duration

	// Session Configuration
	SessionTTL            time.Duration
	SessionIdleTimeout    time.Duration
	MaxConcurrentSessions int

	// Security Configuration
	MaxLoginAttempts         int32
	LoginAttemptWindow       time.Duration
	AccountLockDuration      time.Duration
	PasswordMinLength        int
	PasswordMaxLength        int
	PasswordRequireSpecial   bool
	PasswordRequireNumbers   bool
	PasswordRequireUppercase bool
	PasswordRequireLowercase bool
	PasswordExpirationDays   int

	// Device Trust Configuration
	RequireDeviceVerification bool
	DeviceTrustExpiry         time.Duration
	EnableFingerprinting      bool
	EnableGeoLocation         bool

	// Rate Limiting
	RateLimitPerMinute int
	RateLimitPerHour   int

	// Email Configuration
	EmailVerificationRequired bool
	EmailVerificationExpiry   time.Duration

	// MFA Configuration
	MFARequired      bool
	MFAWindow        time.Duration
	BackupCodesCount int

	// CORS Configuration
	AllowedOrigins []string

	// Audit Configuration
	EnableAuditLogging bool
	AuditRetentionDays int
}

// NewDefaultAuthConfig creates a default auth configuration for development
func NewDefaultAuthConfig() *AuthConfigManager {
	return &AuthConfigManager{
		// JWT
		JWTSecret:           "dev-secret-key-change-in-production",
		JWTAccessTokenTTL:   15 * time.Minute,
		JWTRefreshTokenTTL:  7 * 24 * time.Hour,
		JWTEmailTokenTTL:    24 * time.Hour,
		JWTPasswordResetTTL: 1 * time.Hour,
		JWTIssuer:           "teamart-commerce",
		JWTAudience:         "teamart-api",

		// OTP
		OTPLength:      6,
		OTPTTL:         10 * time.Minute,
		OTPMaxAttempts: 5,
		OTPResendDelay: 1 * time.Minute,

		// Session
		SessionTTL:            24 * time.Hour,
		SessionIdleTimeout:    30 * time.Minute,
		MaxConcurrentSessions: 5,

		// Security
		MaxLoginAttempts:         5,
		LoginAttemptWindow:       15 * time.Minute,
		AccountLockDuration:      30 * time.Minute,
		PasswordMinLength:        8,
		PasswordMaxLength:        128,
		PasswordRequireSpecial:   true,
		PasswordRequireNumbers:   true,
		PasswordRequireUppercase: true,
		PasswordRequireLowercase: true,
		PasswordExpirationDays:   90,

		// Device Trust
		RequireDeviceVerification: true,
		DeviceTrustExpiry:         90 * 24 * time.Hour,
		EnableFingerprinting:      true,
		EnableGeoLocation:         true,

		// Rate Limiting
		RateLimitPerMinute: 60,
		RateLimitPerHour:   1000,

		// Email
		EmailVerificationRequired: true,
		EmailVerificationExpiry:   24 * time.Hour,

		// MFA
		MFARequired:      false,
		MFAWindow:        30 * time.Second,
		BackupCodesCount: 10,

		// Audit
		EnableAuditLogging: true,
		AuditRetentionDays: 90,
	}
}

// NewProductionAuthConfig creates a production auth configuration
func NewProductionAuthConfig(jwtSecret string, origins []string) *AuthConfigManager {
	config := NewDefaultAuthConfig()
	config.JWTSecret = jwtSecret
	config.AllowedOrigins = origins
	config.MFARequired = true
	config.EmailVerificationRequired = true
	config.RequireDeviceVerification = true
	config.MaxLoginAttempts = 3
	config.AccountLockDuration = 1 * time.Hour
	return config
}

// Validate validates the auth configuration
func (c *AuthConfigManager) Validate() error {
	if c.JWTSecret == "" {
		return ErrInvalidCredentials
	}
	if c.JWTAccessTokenTTL == 0 {
		return ErrInvalidToken
	}
	if c.JWTRefreshTokenTTL == 0 {
		return ErrInvalidToken
	}
	if c.OTPLength < 4 || c.OTPLength > 12 {
		return ErrInvalidOTP
	}
	if c.PasswordMinLength < 6 || c.PasswordMinLength > 64 {
		return ErrInvalidPassword
	}
	return nil
}
