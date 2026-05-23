package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/pkg/logger"
)

// AuthHandler handles HTTP requests related to authentication
type AuthHandler struct {
	identityService *auth.IdentityService
	sessionService  *auth.SessionService
	tokenService    *auth.TokenService
	redisService    *auth.RedisService
	logger          *logger.Logger
}

// NewAuthHandler creates a new auth HTTP handler
func NewAuthHandler(
	identityService *auth.IdentityService,
	sessionService *auth.SessionService,
	tokenService *auth.TokenService,
	redisService *auth.RedisService,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		identityService: identityService,
		sessionService:  sessionService,
		tokenService:    tokenService,
		redisService:    redisService,
		logger:          logger,
	}
}

// RefreshTokenRequest represents the HTTP request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id,omitempty"`
}

// RefreshTokenResponse represents the HTTP response body for token refresh
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Message      string `json:"message"`
}

// SignupRequest represents the HTTP request body for user signup
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupResponse represents the HTTP response body for signup
type SignupResponse struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}

// LoginRequest represents the HTTP request body for user login
type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

// LoginResponse represents the HTTP response body for login
type LoginResponse struct {
	UserID           int64  `json:"user_id"`
	SessionID        string `json:"session_id"`
	Email            string `json:"email"`
	Status           string `json:"status"`
	RequiresMFA      bool   `json:"requires_mfa"`
	RequiresPassword bool   `json:"requires_password_verification"`
	Message          string `json:"message"`
	CreatedAt        string `json:"created_at"`
}

// LogoutRequest represents the HTTP request body for logout
type LogoutRequest struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}

// LogoutResponse represents the HTTP response body for logout
type LogoutResponse struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// HandleSignup handles POST /auth/signup requests
//
//	Example: curl -X POST http://localhost:8080/auth/signup \
//	  -H "Content-Type: application/json" \
//	  -d '{"email":"user@example.com","password":"SecurePass123!"}'
func (h *AuthHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling Signup request")

	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		h.respondError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	input := &auth.CreateIdentityInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := h.identityService.CreateIdentity(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		h.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SignupResponse{
		UserID:    output.Identity.ID,
		Email:     output.Identity.Email,
		Status:    string(output.Identity.AccountStatus),
		Message:   "User created successfully. Please verify your email.",
		CreatedAt: output.Identity.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})

	h.logger.Infof("user signed up: %d (%s)", output.Identity.ID, req.Email)
}

// HandleLogin handles POST /auth/login requests
//
//	Example: curl -X POST http://localhost:8080/auth/login \
//	  -H "Content-Type: application/json" \
//	  -d '{"email":"user@example.com","password":"SecurePass123!","user_agent":"Mozilla/5.0...","ip_address":"192.168.1.1"}'
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling Login request")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		h.respondError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get user identity
	identity, err := h.identityService.GetIdentityByEmail(r.Context(), req.Email)
	if err != nil {
		h.logger.Warnf("login failed: user not found (%s)", req.Email)
		h.respondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Get client IP if not provided
	clientIP := req.IPAddress
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	// Get user agent
	userAgent := req.UserAgent
	if userAgent == "" {
		userAgent = r.Header.Get("User-Agent")
	}

	// Verify password and record login attempt
	verifyOutput, err := h.identityService.VerifyPassword(r.Context(), &auth.VerifyPasswordInput{
		UserID:   identity.ID,
		Password: req.Password,
	})
	if err != nil {
		h.logger.Warnf("login failed: invalid password for user %d", identity.ID)
		h.respondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !verifyOutput.IsValid {
		h.respondError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Record successful login
	err = h.identityService.RecordSuccessfulLogin(r.Context(), &auth.RecordSuccessfulLoginInput{
		UserID:    identity.ID,
		IPAddress: clientIP,
	})
	if err != nil {
		h.logger.Errorf("failed to record successful login: %v", err)
	}

	// Create session
	deviceID := r.Header.Get("X-Device-ID")
	if deviceID == "" {
		deviceID = "web_default"
	}

	createSessionInput := &auth.CreateSessionInput{
		UserID:    identity.ID,
		DeviceID:  deviceID,
		UserAgent: userAgent,
		IPAddress: clientIP,
	}

	sessionOutput, err := h.sessionService.CreateSession(r.Context(), createSessionInput)
	if err != nil {
		h.logger.Errorf("service error creating session: %v", err)
		h.respondError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{
		UserID:           identity.ID,
		SessionID:        sessionOutput.Session.ID,
		Email:            identity.Email,
		Status:           string(sessionOutput.Session.TrustLevel),
		RequiresMFA:      sessionOutput.Session.RequiresMFAStep,
		RequiresPassword: sessionOutput.Session.RequiresPasswordVerification,
		Message:          "Login successful",
		CreatedAt:        sessionOutput.Session.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})

	h.logger.Infof("user logged in: %d (%s) - session %s", identity.ID, req.Email, sessionOutput.Session.ID)
}

// HandleRefresh handles POST /auth/refresh requests
//
// This endpoint implements secure token rotation:
// - Validates the refresh token
// - Verifies the session is still active
// - Revokes the old refresh token (token rotation)
// - Issues new access and refresh tokens
// - Returns the new token pair
//
//	Example: curl -X POST http://localhost:8080/auth/refresh \
//	  -H "Content-Type: application/json" \
//	  -H "Authorization: Bearer <access_token>" \
//	  -d '{"refresh_token":"<refresh_token>","session_id":"abc123..."}'
func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling token refresh request")

	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		h.respondError(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	// Validate refresh token
	validateInput := &auth.ValidateTokenInput{
		Token:     req.RefreshToken,
		TokenType: auth.TokenTypeRefresh,
	}

	validateOutput, err := h.tokenService.ValidateToken(r.Context(), validateInput)
	if err != nil {
		h.logger.Errorf("token validation error: %v", err)
		h.respondError(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	if !validateOutput.IsValid {
		h.logger.Warnf("invalid refresh token: %v", validateOutput.Error)
		h.respondError(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	userID := validateOutput.Claims.UserID
	deviceID := validateOutput.Claims.DeviceID
	sessionID := validateOutput.Claims.SessionID

	// If sessionID not provided in request, use the one from token
	if req.SessionID != "" {
		sessionID = req.SessionID
	}

	// Verify session is still active
	sessionExists, err := h.sessionService.SessionRepository().SessionExists(r.Context(), sessionID)
	if err != nil {
		h.logger.Errorf("failed to check session existence: %v", err)
		h.respondError(w, "Failed to validate session", http.StatusInternalServerError)
		return
	}

	if !sessionExists {
		h.logger.Warnf("session not found or revoked: %s for user %d", sessionID, userID)
		h.respondError(w, "Session no longer valid", http.StatusUnauthorized)
		return
	}

	// Get identity to ensure user is still active
	identity, err := h.identityService.GetIdentityByID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("failed to get identity: %v", err)
		h.respondError(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Check if account is active
	if identity.AccountStatus != auth.AccountStatusActive {
		h.logger.Warnf("user account not active: %d (status: %s)", userID, identity.AccountStatus)
		h.respondError(w, "User account is not active", http.StatusForbidden)
		return
	}

	// Refresh the token with rotation
	refreshOutput, err := h.tokenService.RefreshToken(r.Context(), &auth.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
		SessionID:    sessionID,
		DeviceID:     deviceID,
	})
	if err != nil {
		h.logger.Errorf("token refresh failed: %v", err)
		h.respondError(w, "Failed to refresh token", http.StatusUnauthorized)
		return
	}

	// Blacklist the old refresh token (token rotation)
	oldClaims := validateOutput.Claims
	if oldClaims.RegisteredClaims.ID != "" {
		// Blacklist with TTL equal to the original refresh token TTL
		ttl := h.tokenService.GetTokenExpiryTime(oldClaims).Sub(time.Now())
		if ttl > 0 {
			err = h.redisService.BlacklistToken(r.Context(), oldClaims.RegisteredClaims.ID, ttl)
			if err != nil {
				h.logger.Warnf("failed to blacklist old refresh token: %v", err)
				// Don't fail the refresh, just log warning
			}
		}
	}

	// Touch session to update activity
	err = h.sessionService.SessionRepository().TouchSession(r.Context(), sessionID)
	if err != nil {
		h.logger.Warnf("failed to touch session: %v", err)
		// Don't fail, just log warning
	}

	// Return new tokens
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RefreshTokenResponse{
		AccessToken:  refreshOutput.NewAccessToken,
		RefreshToken: refreshOutput.NewRefreshToken,
		ExpiresIn:    refreshOutput.ExpiresIn,
		TokenType:    "Bearer",
		Message:      "Token refreshed successfully",
	})

	h.logger.Infof("token refreshed for user %d (session: %s)", userID, sessionID)
}

// HandleLogout handles POST /auth/logout requests
//
//	Example: curl -X POST http://localhost:8080/auth/logout \
//	  -H "Content-Type: application/json" \
//	  -d '{"session_id":"abc123...","reason":"user logout"}'
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling Logout request")

	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" {
		h.respondError(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	reason := req.Reason
	if reason == "" {
		reason = "user logout"
	}

	// Revoke session
	_, err := h.sessionService.RevokeSession(r.Context(), &auth.RevokeSessionInput{
		SessionID: req.SessionID,
		Reason:    reason,
	})
	if err != nil {
		h.logger.Errorf("service error revoking session: %v", err)
		h.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LogoutResponse{
		SessionID: req.SessionID,
		Status:    "revoked",
		Message:   "Logout successful",
	})

	h.logger.Infof("user logged out: session %s", req.SessionID)
}

// respondError writes an error response to the client
func (h *AuthHandler) respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	})
}

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(mux Router, handler *AuthHandler) {
	mux.HandleFunc("POST /auth/signup", handler.HandleSignup)
	mux.HandleFunc("POST /auth/login", handler.HandleLogin)
	mux.HandleFunc("POST /auth/refresh", handler.HandleRefresh)
	mux.HandleFunc("POST /auth/logout", handler.HandleLogout)
}
