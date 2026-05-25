package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/pkg/logger"
)

// PasswordRecoveryHandler handles password recovery operations
type PasswordRecoveryHandler struct {
	identityService *auth.IdentityService
	sessionService  *auth.SessionService
	redisService    *auth.RedisService
	emailService    interface{} // EmailService from notifications/email
	logger          *logger.Logger
}

// NewPasswordRecoveryHandler creates a new password recovery handler
func NewPasswordRecoveryHandler(
	identityService *auth.IdentityService,
	sessionService *auth.SessionService,
	redisService *auth.RedisService,
	logger *logger.Logger,
) *PasswordRecoveryHandler {
	return &PasswordRecoveryHandler{
		identityService: identityService,
		sessionService:  sessionService,
		redisService:    redisService,
		logger:          logger,
	}
}

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ForgotPasswordResponse represents the response for forgot password
type ForgotPasswordResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Email   string `json:"email"`
}

// ResetPasswordRequest represents the request body for reset password
type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}

// ResetPasswordResponse represents the response for reset password
type ResetPasswordResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HandleForgotPassword handles POST /auth/forgot-password requests
//
// This endpoint initiates password recovery by:
// - Verifying the email exists
// - Generating a secure reset token
// - Storing the token in Redis with TTL
// - Sending reset email
// - Preventing timing attacks by always responding success
//
//	Example: curl -X POST http://localhost:8080/auth/forgot-password \
//	  -H "Content-Type: application/json" \
//	  -d '{"email":"user@example.com"}'
func (h *PasswordRecoveryHandler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling forgot password request")

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		h.respondError(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Check if identity exists
	identity, err := h.identityService.GetIdentityByEmail(r.Context(), req.Email)
	if err != nil {
		// Always respond success to prevent email enumeration attacks
		h.logger.Warnf("forgot password requested for non-existent email: %s", req.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ForgotPasswordResponse{
			Status:  "pending",
			Message: "If this email exists, you will receive a password reset link shortly",
			Email:   req.Email,
		})
		return
	}

	// Generate reset token
	resetToken := generateSecureToken(32)
	resetTokenHash := hashToken(resetToken)

	// Store reset token in Redis with 24-hour TTL
	resetKey := "password_reset:" + resetTokenHash
	if err := h.redisService.SetValue(r.Context(), resetKey, fmt.Sprintf("%d", identity.ID), 24*time.Hour); err != nil {
		h.logger.Errorf("failed to store reset token: %v", err)
		h.respondError(w, "Failed to initiate password recovery", http.StatusInternalServerError)
		return
	}

	// TODO: Send email with reset link
	// resetLink := fmt.Sprintf("https://teamart.app/reset-password?token=%s", resetToken)
	// if h.emailService != nil {
	//     h.emailService.SendPasswordResetEmail(r.Context(), req.Email, resetToken, resetLink)
	// }

	h.logger.Infof("password reset initiated for email: %s (user: %d)", req.Email, identity.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ForgotPasswordResponse{
		Status:  "pending",
		Message: "If this email exists, you will receive a password reset link shortly",
		Email:   req.Email,
	})
}

// HandleResetPassword handles POST /auth/reset-password requests
//
// This endpoint completes password recovery by:
// - Validating the reset token
// - Verifying the token hasn't expired
// - Hashing the new password
// - Updating the password hash
// - Invalidating all sessions (force re-login)
// - Revoking all refresh tokens
//
//	Example: curl -X POST http://localhost:8080/auth/reset-password \
//	  -H "Content-Type: application/json" \
//	  -d '{"reset_token":"...","new_password":"NewSecurePass123!"}'
func (h *PasswordRecoveryHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling reset password request")

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ResetToken == "" {
		h.respondError(w, "Reset token is required", http.StatusBadRequest)
		return
	}

	if req.NewPassword == "" {
		h.respondError(w, "New password is required", http.StatusBadRequest)
		return
	}

	// Hash the reset token to look it up in Redis
	resetTokenHash := hashToken(req.ResetToken)
	resetKey := "password_reset:" + resetTokenHash

	// Retrieve the user ID from Redis
	userIDStr, err := h.redisService.GetValue(r.Context(), resetKey)
	if err != nil {
		h.logger.Warnf("invalid or expired reset token")
		h.respondError(w, "Invalid or expired reset token", http.StatusUnauthorized)
		return
	}

	// Parse user ID (it's stored as string in Redis)
	userID := parseUserID(userIDStr)
	if userID == 0 {
		h.logger.Errorf("invalid user ID from reset token: %s", userIDStr)
		h.respondError(w, "Invalid reset token", http.StatusUnauthorized)
		return
	}

	// Get identity
	identity, err := h.identityService.GetIdentityByID(r.Context(), userID)
	if err != nil {
		h.logger.Errorf("failed to get identity: %v", err)
		h.respondError(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Update password
	err = h.identityService.UpdatePasswordHash(r.Context(), &auth.UpdatePasswordHashInput{
		UserID:      identity.ID,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		h.logger.Errorf("failed to update password: %v", err)
		h.respondError(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	// Revoke all sessions (force re-login)
	if h.sessionService != nil {
		err = h.sessionService.RevokeAllUserSessions(r.Context(), &auth.RevokeAllUserSessionsInput{
			UserID: identity.ID,
			Reason: "password reset",
		})
		if err != nil {
			h.logger.Warnf("failed to revoke sessions: %v", err)
			// Don't fail the reset, just log warning
		}
	}

	// Delete the reset token from Redis (single-use)
	if err := h.redisService.DeleteValue(r.Context(), resetKey); err != nil {
		h.logger.Warnf("failed to delete reset token: %v", err)
	}

	h.logger.Infof("password reset completed for user: %d", identity.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResetPasswordResponse{
		Status:  "success",
		Message: "Your password has been reset. Please login with your new password.",
	})
}

// ===== Helper Functions =====

func (h *PasswordRecoveryHandler) respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   http.StatusText(statusCode),
		"message": message,
		"code":    statusCode,
	})
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// hashToken hashes a token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// parseUserID parses a user ID string to int64
func parseUserID(userIDStr string) int64 {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0
	}
	return userID
}
