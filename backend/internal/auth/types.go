package auth

import "time"

// ===== State Management =====

// OnboardingState represents the current state of a user's onboarding
type OnboardingState string

const (
	// StateNew: User account created, awaiting email verification
	StateNew OnboardingState = "new"
	
	// StateEmailVerified: Email has been verified
	StateEmailVerified OnboardingState = "email_verified"
	
	// StateProfileComplete: Profile information filled in
	StateProfileComplete OnboardingState = "profile_complete"
	
	// StateOnboarded: Full onboarding complete, account is active
	StateOnboarded OnboardingState = "onboarded"
	
	// StateSuspended: Account suspended
	StateSuspended OnboardingState = "suspended"
	
	// StateDeactivated: Account deactivated by user
	StateDeactivated OnboardingState = "deactivated"
)

// ValidStateTransition checks if a state transition is valid
func ValidStateTransition(from, to OnboardingState) bool {
	transitions := map[OnboardingState][]OnboardingState{
		StateNew: {
			StateEmailVerified,
			StateSuspended,
		},
		StateEmailVerified: {
			StateProfileComplete,
			StateNew,
			StateSuspended,
		},
		StateProfileComplete: {
			StateOnboarded,
			StateEmailVerified,
			StateSuspended,
		},
		StateOnboarded: {
			StateSuspended,
			StateDeactivated,
		},
		StateSuspended: {
			StateOnboarded,
			StateNew,
		},
		StateDeactivated: {
			StateOnboarded,
		},
	}
	
	allowed, exists := transitions[from]
	if !exists {
		return false
	}
	
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// ===== Identity =====

// UserIdentity represents the core identity of a user
type UserIdentity struct {
	ID              int64
	Email           string
	PasswordHash    string
	OnboardingState OnboardingState
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ===== JWT Tokens =====

// TokenType represents the type of JWT token
type TokenType string

const (
	// TokenTypeAccess: Short-lived access token
	TokenTypeAccess TokenType = "access"
	
	// TokenTypeRefresh: Long-lived refresh token
	TokenTypeRefresh TokenType = "refresh"
	
	// TokenTypeEmailVerification: Email verification token
	TokenTypeEmailVerification TokenType = "email_verification"
	
	// TokenTypePasswordReset: Password reset token
	TokenTypePasswordReset TokenType = "password_reset"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID      int64     `json:"user_id"`
	Email       string    `json:"email"`
	TokenType   TokenType `json:"token_type"`
	SessionID   string    `json:"session_id,omitempty"`
	DeviceID    string    `json:"device_id,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	IssuedAt    time.Time `json:"iat"`
	ExpiresAt   time.Time `json:"exp"`
	NotBefore   time.Time `json:"nbf"`
	JRTI        string    `json:"jti"` // JWT Token ID for rotation tracking
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // Seconds until access token expires
}

// ===== Sessions =====

// Session represents an authenticated session
type Session struct {
	ID              string
	UserID          int64
	DeviceID        string
	DeviceFingerprint string
	IPAddress       string
	UserAgent       string
	AccessTokenJTI  string // Current access token's JTI
	RefreshTokenJTI string // Current refresh token's JTI
	ExpiresAt       time.Time
	LastActivityAt  time.Time
	IsRevoked       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ===== OTP (One-Time Password) =====

// OTPType represents the type of OTP
type OTPType string

const (
	// OTPTypeEmail: OTP sent via email
	OTPTypeEmail OTPType = "email"
	
	// OTPTypeSMS: OTP sent via SMS
	OTPTypeSMS OTPType = "sms"
	
	// OTPTypeAuthenticator: TOTP from authenticator app
	OTPTypeAuthenticator OTPType = "authenticator"
)

// OTPVerification represents an OTP verification request
type OTPVerification struct {
	ID            string
	UserID        int64
	Type          OTPType
	Code          string // Hashed OTP code
	Destination   string // Email or phone number
	Attempts      int32
	MaxAttempts   int32
	IsVerified    bool
	ExpiresAt     time.Time
	VerifiedAt    time.Time
	CreatedAt     time.Time
}

// ===== RBAC (Role-Based Access Control) =====

// Role represents a role in the system
type Role struct {
	ID          int64
	Name        string // "admin", "moderator", "user", etc.
	Description string
	Permissions []string // List of permission identifiers
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRole represents the assignment of a role to a user
type UserRole struct {
	ID        int64
	UserID    int64
	RoleID    int64
	Role      *Role // Denormalized for convenience
	GrantedAt time.Time
	GrantedBy int64 // User who granted this role
	ExpiresAt *time.Time // Optional expiration
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Permission represents a specific permission
type Permission struct {
	ID          string // "users:read", "users:write", etc.
	Name        string
	Description string
	Category    string // "users", "products", "orders", etc.
	Resource    string // "users", "products", etc.
	Action      string // "read", "write", "delete", etc.
}

// ===== Account Recovery =====

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID        string
	UserID    int64
	Token     string // Hashed token
	Email     string
	ExpiresAt time.Time
	IsUsed    bool
	UsedAt    *time.Time
	CreatedAt time.Time
}

// AccountRecovery represents account recovery process
type AccountRecovery struct {
	ID              string
	UserID          int64
	RecoveryCode    string // Hashed recovery code
	BackupCodes     []string // Hashed backup codes for account recovery
	ExpiresAt       time.Time
	IsRevoked       bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ===== Device Trust =====

// DeviceTrust represents a trusted device
type DeviceTrust struct {
	ID                string
	UserID            int64
	DeviceID          string
	DeviceFingerprint string
	DeviceName        string // "iPhone 14", "Chrome on MacBook", etc.
	DeviceType        string // "mobile", "desktop", "tablet"
	IPAddress         string
	UserAgent         string
	IsVerified        bool
	VerifiedAt        *time.Time
	LastSeenAt        time.Time
	ExpiresAt         time.Time
	IsRevoked         bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ===== Abuse Prevention =====

// LoginAttempt represents a login attempt record
type LoginAttempt struct {
	ID          string
	Email       string
	IPAddress   string
	UserAgent   string
	Success     bool
	FailureReason string // "invalid_password", "user_not_found", etc.
	Timestamp   time.Time
}

// SecurityEvent represents a security event
type SecurityEvent struct {
	ID          string
	UserID      int64
	EventType   string // "login_success", "login_failed", "password_changed", etc.
	Severity    string // "low", "medium", "high", "critical"
	Description string
	IPAddress   string
	UserAgent   string
	Metadata    map[string]interface{}
	Timestamp   time.Time
}

// ===== Auth Configuration =====

// AuthConfig represents authentication configuration
type AuthConfig struct {
	// JWT Configuration
	JWTSecret              string
	JWTAccessTokenTTL      time.Duration // Access token expiration
	JWTRefreshTokenTTL     time.Duration // Refresh token expiration
	JWTEmailTokenTTL       time.Duration // Email verification token expiration
	JWTPasswordResetTTL    time.Duration // Password reset token expiration
	
	// OTP Configuration
	OTPLength              int           // Number of digits in OTP
	OTPTTL                 time.Duration // OTP expiration
	OTPMaxAttempts         int32
	
	// Session Configuration
	SessionTTL             time.Duration // Session expiration
	SessionIdleTimeout     time.Duration // Idle session timeout
	
	// Security Configuration
	MaxLoginAttempts       int32
	LoginAttemptWindow     time.Duration
	PasswordMinLength      int
	PasswordRequireSpecial bool
	PasswordRequireNumbers bool
	
	// Device Trust
	RequireDeviceVerification bool
}

// ===== Error Types =====

// AuthError represents an authentication error
type AuthError struct {
	Code    string // "invalid_credentials", "user_not_found", etc.
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// Common auth errors
var (
	ErrInvalidCredentials      = &AuthError{Code: "invalid_credentials", Message: "Invalid email or password"}
	ErrUserNotFound            = &AuthError{Code: "user_not_found", Message: "User not found"}
	ErrUserNotActive           = &AuthError{Code: "user_not_active", Message: "User account is not active"}
	ErrInvalidToken            = &AuthError{Code: "invalid_token", Message: "Invalid or expired token"}
	ErrTokenExpired            = &AuthError{Code: "token_expired", Message: "Token has expired"}
	ErrInvalidOTP              = &AuthError{Code: "invalid_otp", Message: "Invalid OTP code"}
	ErrOTPExpired              = &AuthError{Code: "otp_expired", Message: "OTP has expired"}
	ErrOTPMaxAttempts          = &AuthError{Code: "otp_max_attempts", Message: "Maximum OTP attempts exceeded"}
	ErrInvalidPassword          = &AuthError{Code: "invalid_password", Message: "Invalid password"}
	ErrEmailAlreadyRegistered  = &AuthError{Code: "email_already_registered", Message: "Email is already registered"}
	ErrInvalidStateTransition  = &AuthError{Code: "invalid_state_transition", Message: "Invalid state transition"}
	ErrAccountSuspended        = &AuthError{Code: "account_suspended", Message: "Account is suspended"}
	ErrAccountDeactivated      = &AuthError{Code: "account_deactivated", Message: "Account is deactivated"}
	ErrPermissionDenied        = &AuthError{Code: "permission_denied", Message: "Permission denied"}
)
