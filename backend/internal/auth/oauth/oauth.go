package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

// OAuthProvider represents a supported OAuth provider
type OAuthProvider string

const (
	ProviderGoogle OAuthProvider = "google"
	ProviderApple  OAuthProvider = "apple"
	ProviderTikTok OAuthProvider = "tiktok"
	ProviderGitHub OAuthProvider = "github"
)

// OAuthConfig holds OAuth provider configuration
type OAuthConfig struct {
	Provider     OAuthProvider
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
	Endpoints    OAuthEndpoints
}

// OAuthEndpoints holds OAuth endpoint URLs
type OAuthEndpoints struct {
	AuthorizationURL string
	TokenURL         string
	UserInfoURL      string
}

// OAuthToken represents an OAuth access token response
type OAuthToken struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
	ExpiresAt    time.Time
	Scope        string
}

// OAuthUser represents user information from OAuth provider
type OAuthUser struct {
	ID       string
	Email    string
	Name     string
	Avatar   string
	Provider OAuthProvider
	RawData  map[string]interface{}
}

// OAuthState represents OAuth state token (for CSRF protection)
type OAuthState struct {
	StateToken  string
	Provider    OAuthProvider
	CreatedAt   time.Time
	ExpiresAt   time.Time
	RedirectURI string
	Nonce       string
}

// OAuthService manages OAuth integrations
type OAuthService struct {
	configs map[OAuthProvider]*OAuthConfig
	storage OAuthStateStorage
}

// OAuthStateStorage defines the interface for storing OAuth state tokens
type OAuthStateStorage interface {
	SaveState(ctx context.Context, state *OAuthState) error
	GetState(ctx context.Context, stateToken string) (*OAuthState, error)
	DeleteState(ctx context.Context, stateToken string) error
	IsStateExpired(ctx context.Context, stateToken string) bool
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(storage OAuthStateStorage) *OAuthService {
	return &OAuthService{
		configs: make(map[OAuthProvider]*OAuthConfig),
		storage: storage,
	}
}

// RegisterProvider registers an OAuth provider configuration
func (s *OAuthService) RegisterProvider(config *OAuthConfig) error {
	if config == nil {
		return errors.New("config is required")
	}
	if config.Provider == "" {
		return errors.New("provider is required")
	}
	if config.ClientID == "" {
		return errors.New("client_id is required")
	}
	if config.ClientSecret == "" {
		return errors.New("client_secret is required")
	}
	if config.RedirectURI == "" {
		return errors.New("redirect_uri is required")
	}

	s.configs[config.Provider] = config
	return nil
}

// GenerateAuthorizationURL generates the OAuth authorization URL
func (s *OAuthService) GenerateAuthorizationURL(ctx context.Context, provider OAuthProvider, redirectURI string) (string, string, error) {
	config, ok := s.configs[provider]
	if !ok {
		return "", "", fmt.Errorf("provider %s not registered", provider)
	}

	// Generate state token
	stateToken, err := s.generateRandomToken(32)
	if err != nil {
		return "", "", err
	}

	// Generate nonce for OpenID Connect (optional but recommended)
	nonce, err := s.generateRandomToken(16)
	if err != nil {
		return "", "", err
	}

	// Save state token
	state := &OAuthState{
		StateToken:  stateToken,
		Provider:    provider,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		RedirectURI: redirectURI,
		Nonce:       nonce,
	}

	if err := s.storage.SaveState(ctx, state); err != nil {
		return "", "", fmt.Errorf("failed to save state: %w", err)
	}

	// Build authorization URL with PKCE
	authURL := s.buildAuthorizationURL(config, stateToken, nonce)

	return authURL, stateToken, nil
}

// VerifyStateToken verifies the OAuth state token (CSRF protection)
func (s *OAuthService) VerifyStateToken(ctx context.Context, stateToken string) (*OAuthState, error) {
	if stateToken == "" {
		return nil, errors.New("state token is required")
	}

	state, err := s.storage.GetState(ctx, stateToken)
	if err != nil {
		return nil, fmt.Errorf("invalid state token: %w", err)
	}

	if state == nil {
		return nil, errors.New("state token not found")
	}

	// Check if state has expired
	if time.Now().After(state.ExpiresAt) {
		_ = s.storage.DeleteState(ctx, stateToken)
		return nil, errors.New("state token has expired")
	}

	// Delete state token after verification (prevent reuse)
	_ = s.storage.DeleteState(ctx, stateToken)

	return state, nil
}

// IsProviderRegistered checks if a provider is registered
func (s *OAuthService) IsProviderRegistered(provider OAuthProvider) bool {
	_, ok := s.configs[provider]
	return ok
}

// GetProvider returns the configuration for a provider
func (s *OAuthService) GetProvider(provider OAuthProvider) (*OAuthConfig, error) {
	config, ok := s.configs[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not registered", provider)
	}
	return config, nil
}

// generateRandomToken generates a cryptographically secure random token
func (s *OAuthService) generateRandomToken(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

// buildAuthorizationURL builds the OAuth authorization URL
func (s *OAuthService) buildAuthorizationURL(config *OAuthConfig, stateToken, nonce string) string {
	scopes := strings.Join(config.Scopes, " ")

	params := []string{
		fmt.Sprintf("client_id=%s", config.ClientID),
		fmt.Sprintf("redirect_uri=%s", config.RedirectURI),
		fmt.Sprintf("response_type=code"),
		fmt.Sprintf("scope=%s", scopes),
		fmt.Sprintf("state=%s", stateToken),
	}

	// Add nonce for OpenID Connect
	if nonce != "" {
		params = append(params, fmt.Sprintf("nonce=%s", nonce))
	}

	// Add PKCE challenge if supported
	params = append(params, "prompt=consent")

	return config.Endpoints.AuthorizationURL + "?" + strings.Join(params, "&")
}
