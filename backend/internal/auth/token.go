package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/teamart/commerce-api/pkg/logger"
)

// TokenService manages JWT token generation, validation, and lifecycle
type TokenService struct {
	config *AuthConfig
	logger *logger.Logger
}

// NewTokenService creates a new token service
func NewTokenService(config *AuthConfig, logger *logger.Logger) *TokenService {
	return &TokenService{
		config: config,
		logger: logger,
	}
}

// GenerateTokenPairInput represents input for token pair generation
type GenerateTokenPairInput struct {
	UserID       int64
	Email        string
	SessionID    string
	DeviceID     string
	Permissions  []string
}

// GenerateTokenPair generates an access token and refresh token pair
func (ts *TokenService) GenerateTokenPair(ctx context.Context, input *GenerateTokenPairInput) (*TokenPair, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	ts.logger.Debugf("generating token pair for user %d", input.UserID)

	// Generate access token
	accessToken, err := ts.generateAccessToken(input)
	if err != nil {
		ts.logger.Errorf("failed to generate access token: %v", err)
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := ts.generateRefreshToken(input)
	if err != nil {
		ts.logger.Errorf("failed to generate refresh token: %v", err)
		return nil, err
	}

	ts.logger.Infof("token pair generated for user %d", input.UserID)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(ts.config.JWTAccessTokenTTL.Seconds()),
	}, nil
}

// GenerateAccessTokenInput for access token generation
type GenerateAccessTokenInput struct {
	UserID      int64
	Email       string
	SessionID   string
	DeviceID    string
	Permissions []string
}

// generateAccessToken generates an access token
func (ts *TokenService) generateAccessToken(input *GenerateTokenPairInput) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ts.config.JWTAccessTokenTTL)
	jti := ts.generateJTI()

	claims := JWTClaims{
		UserID:      input.UserID,
		Email:       input.Email,
		TokenType:   TokenTypeAccess,
		SessionID:   input.SessionID,
		DeviceID:    input.DeviceID,
		Permissions: input.Permissions,
		IssuedAt:    now,
		ExpiresAt:   expiresAt,
		NotBefore:   now,
		JRTI:        jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(ts.config.JWTSecret))
	if err != nil {
		return "", err
	}

	ts.logger.Debugf("access token generated with JTI: %s", jti)
	return signedToken, nil
}

// generateRefreshToken generates a refresh token
func (ts *TokenService) generateRefreshToken(input *GenerateTokenPairInput) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ts.config.JWTRefreshTokenTTL)
	jti := ts.generateJTI()

	claims := JWTClaims{
		UserID:    input.UserID,
		Email:     input.Email,
		TokenType: TokenTypeRefresh,
		SessionID: input.SessionID,
		DeviceID:  input.DeviceID,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
		NotBefore: now,
		JRTI:      jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(ts.config.JWTSecret))
	if err != nil {
		return "", err
	}

	ts.logger.Debugf("refresh token generated with JTI: %s", jti)
	return signedToken, nil
}

// ValidateTokenInput represents input for token validation
type ValidateTokenInput struct {
	Token        string
	TokenType    TokenType
	ExpectJTI    string // Optional: for refresh rotation verification
}

// ValidateTokenOutput represents the result of token validation
type ValidateTokenOutput struct {
	IsValid bool
	Claims  *JWTClaims
	Error   error
}

// ValidateToken validates a JWT token
func (ts *TokenService) ValidateToken(ctx context.Context, input *ValidateTokenInput) (*ValidateTokenOutput, error) {
	if input.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	claims := &JWTClaims{}
	
	token, err := jwt.ParseWithClaims(input.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(ts.config.JWTSecret), nil
	})

	if err != nil {
		ts.logger.Debugf("token validation failed: %v", err)
		return &ValidateTokenOutput{
			IsValid: false,
			Error:   ErrInvalidToken,
		}, nil
	}

	if !token.Valid {
		ts.logger.Debugf("token is invalid")
		return &ValidateTokenOutput{
			IsValid: false,
			Error:   ErrInvalidToken,
		}, nil
	}

	// Verify token type matches expected
	if input.TokenType != "" && claims.TokenType != input.TokenType {
		ts.logger.Debugf("token type mismatch: expected %s, got %s", input.TokenType, claims.TokenType)
		return &ValidateTokenOutput{
			IsValid: false,
			Error:   ErrInvalidToken,
		}, nil
	}

	// Verify JTI if provided (for refresh rotation)
	if input.ExpectJTI != "" && claims.JRTI != input.ExpectJTI {
		ts.logger.Debugf("token JTI mismatch: expected %s, got %s", input.ExpectJTI, claims.JRTI)
		return &ValidateTokenOutput{
			IsValid: false,
			Error:   ErrInvalidToken,
		}, nil
	}

	ts.logger.Debugf("token validation successful for user %d", claims.UserID)

	return &ValidateTokenOutput{
		IsValid: true,
		Claims:  claims,
	}, nil
}

// RefreshTokenInput represents input for token refresh
type RefreshTokenInput struct {
	RefreshToken string
	SessionID    string
	DeviceID     string
}

// RefreshTokenOutput represents the result of token refresh
type RefreshTokenOutput struct {
	NewAccessToken string
	NewRefreshToken string // Token rotation: new refresh token
	ExpiresIn      int64
	OldRefreshToken string // For revocation
}

// RefreshToken refreshes an access token using a refresh token
func (ts *TokenService) RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error) {
	if input.RefreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	ts.logger.Debugf("refreshing token for session %s", input.SessionID)

	// Validate refresh token
	validateInput := &ValidateTokenInput{
		Token:     input.RefreshToken,
		TokenType: TokenTypeRefresh,
	}

	result, err := ts.ValidateToken(ctx, validateInput)
	if err != nil {
		ts.logger.Errorf("refresh token validation failed: %v", err)
		return nil, err
	}

	if !result.IsValid {
		ts.logger.Warnf("invalid refresh token for session %s", input.SessionID)
		return nil, result.Error
	}

	oldJTI := result.Claims.JRTI

	// Generate new token pair with rotation
	tokenInput := &GenerateTokenPairInput{
		UserID:      result.Claims.UserID,
		Email:       result.Claims.Email,
		SessionID:   input.SessionID,
		DeviceID:    input.DeviceID,
		Permissions: result.Claims.Permissions,
	}

	tokenPair, err := ts.GenerateTokenPair(ctx, tokenInput)
	if err != nil {
		return nil, err
	}

	ts.logger.Infof("token refreshed for user %d, old refresh JTI: %s", result.Claims.UserID, oldJTI)

	return &RefreshTokenOutput{
		NewAccessToken:  tokenPair.AccessToken,
		NewRefreshToken: tokenPair.RefreshToken,
		ExpiresIn:       tokenPair.ExpiresIn,
		OldRefreshToken: input.RefreshToken, // For revocation tracking
	}, nil
}

// RevokeTokenInput represents input for token revocation
type RevokeTokenInput struct {
	UserID int64
	JTI    string
	Reason string
}

// RevokeToken revokes a token by its JTI
func (ts *TokenService) RevokeToken(ctx context.Context, input *RevokeTokenInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if input.JTI == "" {
		return fmt.Errorf("JTI is required")
	}

	ts.logger.Infof("revoking token JTI %s for user %d (reason: %s)", 
		input.JTI, input.UserID, input.Reason)

	// In a real implementation, this would add the JTI to a blacklist/revocation list

	return nil
}

// ===== Helper Methods =====

// generateJTI generates a unique JWT Token ID
func (ts *TokenService) generateJTI() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// IsTokenExpired checks if a token is expired
func (ts *TokenService) IsTokenExpired(expiresAt time.Time) bool {
	return time.Now().After(expiresAt)
}

// GetTokenExpiryTime returns when a token expires
func (ts *TokenService) GetTokenExpiryTime(claims *JWTClaims) time.Time {
	return claims.ExpiresAt
}

// GetTimeRemainingForToken returns remaining time for token
func (ts *TokenService) GetTimeRemainingForToken(expiresAt time.Time) time.Duration {
	remaining := time.Until(expiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
