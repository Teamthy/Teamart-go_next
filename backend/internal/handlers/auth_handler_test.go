package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/auth"
	authoauth "github.com/teamart/commerce-api/internal/auth/oauth"
	"github.com/teamart/commerce-api/pkg/logger"
)

func TestAuthHandlerLoginReturnsTokenPairAndUser(t *testing.T) {
	log := logger.NewNoop()
	cfg := &auth.AuthConfig{
		JWTSecret:              "test-secret",
		JWTAccessTokenTTL:      time.Hour,
		JWTRefreshTokenTTL:     24 * time.Hour,
		SessionTTL:             24 * time.Hour,
		PasswordMinLength:      8,
		PasswordRequireSpecial: false,
		PasswordRequireNumbers: false,
		OTPLength:              6,
		OTPTTL:                 10 * time.Minute,
		OTPMaxAttempts:         5,
	}

	identityRepo := auth.NewIdentityRepositoryMemory(log)
	sessionRepo := auth.NewSessionRepositoryMemory(log)
	identityService := auth.NewIdentityService(cfg, log, identityRepo)
	sessionService := auth.NewSessionService(cfg, log, sessionRepo)
	tokenService := auth.NewTokenService(cfg, log)

	_, err := identityService.CreateIdentity(context.Background(), &auth.CreateIdentityInput{
		Email:    "login-user@example.com",
		Password: "StrongPass1",
	})
	if err != nil {
		t.Fatalf("create identity: %v", err)
	}

	handler := NewAuthHandler(identityService, sessionService, tokenService, auth.NewRedisService(nil, log), log)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"login-user@example.com","password":"StrongPass1","user_agent":"test-agent","ip_address":"127.0.0.1"}`))
	resp := httptest.NewRecorder()

	handler.HandleLogin(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, resp.Code, resp.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["access_token"] == nil || body["access_token"] == "" {
		t.Fatalf("expected access_token to be present, got %#v", body)
	}
	if body["refresh_token"] == nil || body["refresh_token"] == "" {
		t.Fatalf("expected refresh_token to be present, got %#v", body)
	}
	if _, ok := body["requires_mfa"]; !ok {
		t.Fatalf("expected requires_mfa to be present, got %#v", body)
	}

	user, ok := body["user"].(map[string]any)
	if !ok {
		t.Fatalf("expected user object, got %#v", body["user"])
	}
	if user["email"] != "login-user@example.com" {
		t.Fatalf("expected user email to be preserved, got %#v", user)
	}
}

func TestAuthHandlerStartsGoogleOAuthRedirect(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "google-client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "google-client-secret")
	t.Setenv("GOOGLE_REDIRECT_URI", "http://localhost:8000/auth/google/callback")

	handler := NewAuthHandler(nil, nil, nil, nil, logger.NewNoop())
	service := authoauth.NewOAuthService(authoauth.NewMemoryStateStorage())
	if err := service.RegisterProvider(&authoauth.OAuthConfig{
		Provider:     authoauth.ProviderGoogle,
		ClientID:     "google-client-id",
		ClientSecret: "google-client-secret",
		RedirectURI:  "http://localhost:8000/auth/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoints: authoauth.OAuthEndpoints{
			AuthorizationURL: "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:         "https://oauth2.googleapis.com/token",
			UserInfoURL:      "https://openidconnect.googleapis.com/v1/userinfo",
		},
	}); err != nil {
		t.Fatalf("register provider: %v", err)
	}
	handler.oauthService = service

	req := httptest.NewRequest(http.MethodGet, "/auth/google", nil)
	req = mux.SetURLVars(req, map[string]string{"provider": "google"})
	resp := httptest.NewRecorder()

	handler.HandleOAuthStart(resp, req)

	if resp.Code != http.StatusFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusFound, resp.Code, resp.Body.String())
	}

	location := resp.Header().Get("Location")
	if !strings.Contains(location, "https://accounts.google.com/o/oauth2/v2/auth") {
		t.Fatalf("expected google authorization redirect, got %q", location)
	}

	parsed, err := url.Parse(location)
	if err != nil {
		t.Fatalf("parse redirect location: %v", err)
	}

	if got := parsed.Query().Get("client_id"); got != "google-client-id" {
		t.Fatalf("expected client_id google-client-id, got %q", got)
	}
	if got := parsed.Query().Get("redirect_uri"); got != "http://localhost:8000/auth/google/callback" {
		t.Fatalf("expected redirect_uri to be preserved, got %q", got)
	}
}
