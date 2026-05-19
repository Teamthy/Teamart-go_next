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
	// Core Fields
	ID           int64  // Unique user ID
	Email        string // User email (unique)
	PasswordHash string // Bcrypt password hash

	// Status & State
	OnboardingState OnboardingState // Current onboarding phase
	AccountStatus   AccountStatus   // Current account status
	IsActive        bool            // Whether account is active

	// Security Fields
	FailedLoginAttempts    int32      // Count of failed login attempts
	FailedLoginLastAttempt *time.Time // Timestamp of last failed login
	LockedUntil            *time.Time // Account locked until this time
	PasswordChangedAt      time.Time  // When password was last changed
	LastLoginAt            *time.Time // Timestamp of last successful login
	LastLoginIP            string     // IP address of last login

	// Contact Fields
	RecoveryEmail *string // Alternative email for account recovery
	PhoneNumber   *string // Phone number for SMS/2FA

	// Flags & Preferences
	RequiresMFA bool    // Whether MFA is required/enabled
	MFAMethod   *string // "email", "sms", "authenticator"

	// Timestamps
	CreatedAt time.Time  // Account creation time
	UpdatedAt time.Time  // Last update time
	DeletedAt *time.Time // Soft delete timestamp (if applicable)
}

// IsActiveAndUnlocked returns true if account is active and not locked
func (u *UserIdentity) IsActiveAndUnlocked() bool {
	if !u.IsActive || u.AccountStatus != AccountStatusActive {
		return false
	}
	if u.LockedUntil != nil && time.Now().Before(*u.LockedUntil) {
		return false
	}
	return true
}

// CanLogin returns true if the user can attempt to login
func (u *UserIdentity) CanLogin() bool {
	return u.IsActiveAndUnlocked() && u.OnboardingState == StateOnboarded
}

// IsPasswordExpired returns true if password is older than provided duration
func (u *UserIdentity) IsPasswordExpired(maxAge time.Duration) bool {
	if u.PasswordChangedAt.IsZero() {
		return false
	}
	return time.Since(u.PasswordChangedAt) > maxAge
}

// LockAccount locks the account until lockDuration expires
func (u *UserIdentity) LockAccount(lockDuration time.Duration) {
	lockUntil := time.Now().Add(lockDuration)
	u.LockedUntil = &lockUntil
}

// UnlockAccount unlocks a locked account
func (u *UserIdentity) UnlockAccount() {
	u.LockedUntil = nil
	u.FailedLoginAttempts = 0
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

// AccountStatus represents the current status of a user account
type AccountStatus string

const (
	// AccountStatusActive: Account is active and can login
	AccountStatusActive AccountStatus = "active"

	// AccountStatusPending: Account is pending email verification
	AccountStatusPending AccountStatus = "pending"

	// AccountStatusLocked: Account is locked (too many failed attempts)
	AccountStatusLocked AccountStatus = "locked"

	// AccountStatusSuspended: Account is suspended by admin
	AccountStatusSuspended AccountStatus = "suspended"

	// AccountStatusDeactivated: Account is deactivated by user
	AccountStatusDeactivated AccountStatus = "deactivated"
)

// TrustLevel represents the trust level of a device/session
type TrustLevel string

const (
	// TrustLevelUntrusted: Device is not trusted
	TrustLevelUntrusted TrustLevel = "untrusted"

	// TrustLevelPartial: Device is partially trusted (needs MFA on sensitive operations)
	TrustLevelPartial TrustLevel = "partial"

	// TrustLevelTrusted: Device is fully trusted
	TrustLevelTrusted TrustLevel = "trusted"
)

// GeoLocation represents geographic information
type GeoLocation struct {
	Country   string
	City      string
	Latitude  float64
	Longitude float64
	Timezone  string
}

// Session represents an authenticated session
type Session struct {
	ID                string       // Unique session ID
	UserID            int64        // User who owns this session
	DeviceID          string       // Device identifier
	DeviceFingerprint string       // Hash of device characteristics
	IPAddress         string       // IP address of the session
	UserAgent         string       // User agent string
	GeoLocation       *GeoLocation // Geographic location of session
	AccessTokenJTI    string       // Current access token's JTI
	RefreshTokenJTI   string       // Current refresh token's JTI
	TrustLevel        TrustLevel   // Trust level of this session
	LastActivityAt    time.Time    // Last activity timestamp
	LastIPAddress     string       // Previous IP address
	LastGeoLocation   *GeoLocation // Previous location (for impossible travel detection)
	MFAVerifiedAt     *time.Time   // When MFA was last verified for this session
	RequiresMFA       bool         // Whether this session requires MFA
	ExpiresAt         time.Time    // Session expiration time
	IsRevoked         bool         // Whether session has been revoked
	RevokedAt         *time.Time   // When session was revoked
	RevokeReason      string       // Reason for revocation
	CreatedAt         time.Time    // Session creation time
	UpdatedAt         time.Time    // Last update time
}

// IsValid checks if a session is valid (not expired, not revoked, active)
func (s *Session) IsValid() bool {
	return !s.IsRevoked && time.Now().Before(s.ExpiresAt)
}

// IsExpired checks if a session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// RequiresMFAStep returns true if this session requires MFA verification
func (s *Session) RequiresMFAStep() bool {
	return s.RequiresMFA || s.TrustLevel == TrustLevelUntrusted
}

// IsInactive checks if session has been inactive for longer than provided duration
func (s *Session) IsInactive(timeout time.Duration) bool {
	return time.Since(s.LastActivityAt) > timeout
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
	ID          string
	UserID      int64
	Type        OTPType
	Code        string // Hashed OTP code
	Destination string // Email or phone number
	Attempts    int32
	MaxAttempts int32
	IsVerified  bool
	ExpiresAt   time.Time
	VerifiedAt  time.Time
	CreatedAt   time.Time
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
	GrantedBy int64      // User who granted this role
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
	ID           string
	UserID       int64
	RecoveryCode string   // Hashed recovery code
	BackupCodes  []string // Hashed backup codes for account recovery
	ExpiresAt    time.Time
	IsRevoked    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	ID            string
	Email         string
	IPAddress     string
	UserAgent     string
	Success       bool
	FailureReason string // "invalid_password", "user_not_found", etc.
	Timestamp     time.Time
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
	JWTSecret           string
	JWTAccessTokenTTL   time.Duration // Access token expiration
	JWTRefreshTokenTTL  time.Duration // Refresh token expiration
	JWTEmailTokenTTL    time.Duration // Email verification token expiration
	JWTPasswordResetTTL time.Duration // Password reset token expiration

	// OTP Configuration
	OTPLength      int           // Number of digits in OTP
	OTPTTL         time.Duration // OTP expiration
	OTPMaxAttempts int32

	// Session Configuration
	SessionTTL         time.Duration // Session expiration
	SessionIdleTimeout time.Duration // Idle session timeout

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
	ErrInvalidCredentials     = &AuthError{Code: "invalid_credentials", Message: "Invalid email or password"}
	ErrUserNotFound           = &AuthError{Code: "user_not_found", Message: "User not found"}
	ErrUserNotActive          = &AuthError{Code: "user_not_active", Message: "User account is not active"}
	ErrInvalidToken           = &AuthError{Code: "invalid_token", Message: "Invalid or expired token"}
	ErrTokenExpired           = &AuthError{Code: "token_expired", Message: "Token has expired"}
	ErrInvalidOTP             = &AuthError{Code: "invalid_otp", Message: "Invalid OTP code"}
	ErrOTPExpired             = &AuthError{Code: "otp_expired", Message: "OTP has expired"}
	ErrOTPMaxAttempts         = &AuthError{Code: "otp_max_attempts", Message: "Maximum OTP attempts exceeded"}
	ErrInvalidPassword        = &AuthError{Code: "invalid_password", Message: "Invalid password"}
	ErrEmailAlreadyRegistered = &AuthError{Code: "email_already_registered", Message: "Email is already registered"}
	ErrInvalidStateTransition = &AuthError{Code: "invalid_state_transition", Message: "Invalid state transition"}
	ErrAccountSuspended       = &AuthError{Code: "account_suspended", Message: "Account is suspended"}
	ErrAccountDeactivated     = &AuthError{Code: "account_deactivated", Message: "Account is deactivated"}
	ErrPermissionDenied       = &AuthError{Code: "permission_denied", Message: "Permission denied"}
)
