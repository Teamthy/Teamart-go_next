package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/teamart/commerce-api/internal/auth"
)

// AuthClaims represents the authenticated user identity for websocket connections.
type AuthClaims struct {
	UserID int64
	Email  string
}

// Authenticator validates incoming websocket requests.
type Authenticator interface {
	Authenticate(r *http.Request) (*AuthClaims, error)
}

// TokenAuthenticator validates access tokens using the existing auth token service.
type TokenAuthenticator struct {
	tokenService *auth.TokenService
}

// NewTokenAuthenticator creates a TokenAuthenticator.
func NewTokenAuthenticator(tokenService *auth.TokenService) *TokenAuthenticator {
	return &TokenAuthenticator{tokenService: tokenService}
}

// Authenticate checks the Authorization header or token query parameter.
func (a *TokenAuthenticator) Authenticate(r *http.Request) (*AuthClaims, error) {
	token := bearerToken(r.Header.Get("Authorization"))
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		return nil, fmt.Errorf("missing access token")
	}

	result, err := a.tokenService.ValidateToken(context.Background(), &auth.ValidateTokenInput{
		Token:     token,
		TokenType: auth.TokenTypeAccess,
	})
	if err != nil {
		return nil, err
	}
	if !result.IsValid {
		return nil, result.Error
	}

	return &AuthClaims{
		UserID: result.Claims.UserID,
		Email:  result.Claims.Email,
	}, nil
}

func bearerToken(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
