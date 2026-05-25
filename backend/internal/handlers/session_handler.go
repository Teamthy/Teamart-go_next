package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/internal/middleware"
	"github.com/teamart/commerce-api/pkg/logger"
)

// SessionHandler handles HTTP requests related to session management
type SessionHandler struct {
	sessionService *auth.SessionService
	logger         *logger.Logger
}

// NewSessionHandler creates a new session HTTP handler
func NewSessionHandler(sessionService *auth.SessionService, logger *logger.Logger) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		logger:         logger,
	}
}

// GetActiveSessionsResponse represents active sessions for a user
type GetActiveSessionsResponse struct {
	UserID   int64           `json:"user_id"`
	Sessions []SessionDetail `json:"sessions"`
	Count    int             `json:"count"`
	Message  string          `json:"message"`
}

// SessionDetail represents detailed session information
type SessionDetail struct {
	SessionID                    string  `json:"session_id"`
	DeviceID                     string  `json:"device_id"`
	DeviceFingerprint            string  `json:"device_fingerprint"`
	UserAgent                    string  `json:"user_agent"`
	IPAddress                    string  `json:"ip_address"`
	TrustLevel                   string  `json:"trust_level"`
	RequiresMFAStep              bool    `json:"requires_mfa_step"`
	RequiresPasswordVerification bool    `json:"requires_password_verification"`
	GeoCountry                   string  `json:"geo_country"`
	GeoCity                      string  `json:"geo_city"`
	CreatedAt                    string  `json:"created_at"`
	LastActivityAt               string  `json:"last_activity_at"`
	ExpiresAt                    string  `json:"expires_at"`
	RevokedAt                    *string `json:"revoked_at"`
}

// ValidateSessionRequest represents a request to validate a session
type ValidateSessionRequest struct {
	SessionID string `json:"session_id"`
	UserAgent string `json:"user_agent,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
}

// ValidateSessionResponse represents the response for session validation
type ValidateSessionResponse struct {
	SessionID                    string `json:"session_id"`
	IsValid                      bool   `json:"is_valid"`
	UserID                       int64  `json:"user_id"`
	TrustLevel                   string `json:"trust_level"`
	RequiresMFAStep              bool   `json:"requires_mfa_step"`
	RequiresPasswordVerification bool   `json:"requires_password_verification"`
	Message                      string `json:"message"`
}

// RevokeSessionRequest represents a request to revoke a session
type RevokeSessionRequest struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}

// RevokeSessionResponse represents the response for revoking a session
type RevokeSessionResponse struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// TrustDeviceRequest represents a request to mark a device as trusted
type TrustDeviceRequest struct {
	SessionID string `json:"session_id"`
	DeviceID  string `json:"device_id"`
}

// TrustDeviceResponse represents the response for trusting a device
type TrustDeviceResponse struct {
	DeviceID   string `json:"device_id"`
	TrustLevel string `json:"trust_level"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

// HandleGetActiveSessions handles GET /sessions/{user_id} requests
// Example: curl -X GET http://localhost:8080/sessions/123
func (h *SessionHandler) HandleGetActiveSessions(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling GetActiveSessions request")

	userIDStr := r.PathValue("user_id")
	if userIDStr == "" {
		h.respondError(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get active sessions
	output, err := h.sessionService.GetUserActiveSessions(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		h.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert sessions to response format
	var sessionDetails []SessionDetail
	for _, sess := range output.Sessions {
		detail := SessionDetail{
			SessionID:                    sess.ID,
			DeviceID:                     sess.DeviceID,
			DeviceFingerprint:            sess.DeviceFingerprint,
			UserAgent:                    sess.UserAgent,
			IPAddress:                    sess.IPAddress,
			TrustLevel:                   string(sess.TrustLevel),
			RequiresMFAStep:              sess.RequiresMFAStep(),
			RequiresPasswordVerification: sess.TrustLevel != auth.TrustLevelTrusted,
			GeoCountry:                   sess.GeoLocation.Country,
			GeoCity:                      sess.GeoLocation.City,
			CreatedAt:                    sess.CreatedAt.Format("2006-01-02T15:04:05Z"),
			LastActivityAt:               sess.LastActivityAt.Format("2006-01-02T15:04:05Z"),
			ExpiresAt:                    sess.ExpiresAt.Format("2006-01-02T15:04:05Z"),
		}
		if !sess.RevokedAt.IsZero() {
			revokedStr := sess.RevokedAt.Format("2006-01-02T15:04:05Z")
			detail.RevokedAt = &revokedStr
		}
		sessionDetails = append(sessionDetails, detail)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetActiveSessionsResponse{
		UserID:   userID,
		Sessions: sessionDetails,
		Count:    len(sessionDetails),
		Message:  "Active sessions retrieved successfully",
	})

	h.logger.Infof("retrieved %d active sessions for user %d", len(sessionDetails), userID)
}

// HandleValidateSession handles POST /sessions/validate requests
//
//	Example: curl -X POST http://localhost:8080/sessions/validate \
//	  -H "Content-Type: application/json" \
//	  -d '{"session_id":"abc123...","user_agent":"Mozilla/5.0...","ip_address":"192.168.1.1"}'
func (h *SessionHandler) HandleValidateSession(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling ValidateSession request")

	var req ValidateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" {
		h.respondError(w, "Session ID is required", http.StatusBadRequest)
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

	// Validate session
	validateOutput, err := h.sessionService.ValidateSession(r.Context(), &auth.ValidateSessionInput{
		SessionID: req.SessionID,
		UserAgent: userAgent,
		IPAddress: clientIP,
	})
	if err != nil {
		h.logger.Warnf("session validation failed: %v", err)
		h.respondError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ValidateSessionResponse{
		SessionID:                    validateOutput.Session.ID,
		IsValid:                      validateOutput.IsValid,
		UserID:                       validateOutput.Session.UserID,
		TrustLevel:                   string(validateOutput.Session.TrustLevel),
		RequiresMFAStep:              validateOutput.Session.RequiresMFAStep(),
		RequiresPasswordVerification: validateOutput.Session.TrustLevel != auth.TrustLevelTrusted,
		Message:                      "Session is valid",
	})

	h.logger.Infof("session validated: %s for user %d", req.SessionID, validateOutput.Session.UserID)
}

// HandleRevokeSession handles POST /sessions/revoke requests
//
//	Example: curl -X POST http://localhost:8080/sessions/revoke \
//	  -H "Content-Type: application/json" \
//	  -d '{"session_id":"abc123...","reason":"user logout"}'
func (h *SessionHandler) HandleRevokeSession(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling RevokeSession request")

	var req RevokeSessionRequest
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
		reason = "user revocation"
	}

	// Revoke session
	err := h.sessionService.RevokeSession(r.Context(), &auth.RevokeSessionInput{
		SessionID: req.SessionID,
		Reason:    reason,
	})
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		h.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RevokeSessionResponse{
		SessionID: req.SessionID,
		Status:    "revoked",
		Message:   "Session revoked successfully",
	})

	h.logger.Infof("session revoked: %s", req.SessionID)
}

// HandleTrustDevice handles POST /sessions/trust-device requests
//
//	Example: curl -X POST http://localhost:8080/sessions/trust-device \
//	  -H "Content-Type: application/json" \
//	  -d '{"session_id":"abc123...","device_id":"device456..."}'
func (h *SessionHandler) HandleTrustDevice(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling TrustDevice request")

	var req TrustDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" || req.DeviceID == "" {
		h.respondError(w, "Session ID and Device ID are required", http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Trust device
	err = h.sessionService.TrustDevice(r.Context(), &auth.TrustDeviceInput{
		UserID:   userID,
		DeviceID: req.DeviceID,
	})
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		h.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TrustDeviceResponse{
		DeviceID:   req.DeviceID,
		TrustLevel: string(auth.TrustLevelTrusted),
		Status:     "trusted",
		Message:    "Device marked as trusted",
	})

	h.logger.Infof("device trusted: %s", req.DeviceID)
}

// respondError writes an error response to the client
func (h *SessionHandler) respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    statusCode,
	})
}

// RegisterSessionRoutes registers all session management routes
func RegisterSessionRoutes(mux Router, handler *SessionHandler) {
	mux.HandleFunc("GET /sessions/{user_id}", handler.HandleGetActiveSessions)
	mux.HandleFunc("POST /sessions/validate", handler.HandleValidateSession)
	mux.HandleFunc("POST /sessions/revoke", handler.HandleRevokeSession)
	mux.HandleFunc("POST /sessions/trust-device", handler.HandleTrustDevice)
}
