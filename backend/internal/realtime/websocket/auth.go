package websocket

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// WebsocketAuthenticator validates websocket upgrade requests.
type WebsocketAuthenticator struct{}

// NewWebsocketAuthenticator creates a new authenticator instance.
func NewWebsocketAuthenticator() *WebsocketAuthenticator {
	return &WebsocketAuthenticator{}
}

// AuthenticateRequest validates the request and returns the user ID.
func (a *WebsocketAuthenticator) AuthenticateRequest(r *http.Request) (int64, error) {
	if idHeader := r.Header.Get("X-User-ID"); idHeader != "" {
		userID, err := strconv.ParseInt(idHeader, 10, 64)
		if err == nil && userID > 0 {
			return userID, nil
		}
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return 0, fmt.Errorf("missing authorization header")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return 0, fmt.Errorf("invalid authorization header")
	}

	// TODO: Replace this placeholder with real token parsing and verification.
	return 0, fmt.Errorf("token validation not implemented")
}
