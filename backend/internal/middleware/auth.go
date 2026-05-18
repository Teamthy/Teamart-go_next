package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/pkg/logger"
)

// AuthContext context key for storing auth information
type contextKey string

const (
	ContextKeyUserID    contextKey = "user_id"
	ContextKeyEmail     contextKey = "email"
	ContextKeyClaims    contextKey = "claims"
	ContextKeySessionID contextKey = "session_id"
	ContextKeyDeviceID  contextKey = "device_id"
)

// AuthMiddleware is middleware for authentication
type AuthMiddleware struct {
	tokenService   *auth.TokenService
	sessionService *auth.SessionService
	logger         *logger.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(
	tokenService *auth.TokenService,
	sessionService *auth.SessionService,
	logger *logger.Logger,
) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService:   tokenService,
		sessionService: sessionService,
		logger:         logger,
	}
}

// Middleware is the HTTP middleware function
func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract bearer token from header
		token, err := extractBearerToken(r)
		if err != nil {
			am.logger.Debugf("failed to extract bearer token: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate token
		result, err := am.tokenService.ValidateToken(r.Context(), &auth.ValidateTokenInput{
			Token:     token,
			TokenType: auth.TokenTypeAccess,
		})
		if err != nil {
			am.logger.Errorf("token validation error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !result.IsValid {
			am.logger.Debugf("invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate session
		sessionResult, err := am.sessionService.ValidateSession(r.Context(), &auth.ValidateSessionInput{
			SessionID: result.Claims.SessionID,
			UserID:    result.Claims.UserID,
			IPAddress: getClientIP(r),
			UserAgent: r.Header.Get("User-Agent"),
		})
		if err != nil || !sessionResult.IsValid {
			am.logger.Warnf("session validation failed for user %d", result.Claims.UserID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Update session activity
		_ = am.sessionService.UpdateSessionActivity(r.Context(), &auth.UpdateSessionActivityInput{
			SessionID: result.Claims.SessionID,
			UserID:    result.Claims.UserID,
			IPAddress: getClientIP(r),
			UserAgent: r.Header.Get("User-Agent"),
		})

		// Store auth info in context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, result.Claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyEmail, result.Claims.Email)
		ctx = context.WithValue(ctx, ContextKeyClaims, result.Claims)
		ctx = context.WithValue(ctx, ContextKeySessionID, result.Claims.SessionID)
		ctx = context.WithValue(ctx, ContextKeyDeviceID, result.Claims.DeviceID)

		am.logger.Debugf("authentication successful for user %d", result.Claims.UserID)

		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware is middleware for optional authentication
type OptionalAuthMiddleware struct {
	tokenService *auth.TokenService
	logger       *logger.Logger
}

// NewOptionalAuthMiddleware creates a new optional auth middleware
func NewOptionalAuthMiddleware(
	tokenService *auth.TokenService,
	logger *logger.Logger,
) *OptionalAuthMiddleware {
	return &OptionalAuthMiddleware{
		tokenService: tokenService,
		logger:       logger,
	}
}

// Middleware is the HTTP middleware function for optional auth
func (oam *OptionalAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract bearer token from header
		token, err := extractBearerToken(r)
		if err != nil {
			// No token, continue without authentication
			oam.logger.Debugf("no bearer token provided, continuing without auth")
			next.ServeHTTP(w, r)
			return
		}

		// Validate token
		result, err := oam.tokenService.ValidateToken(r.Context(), &auth.ValidateTokenInput{
			Token:     token,
			TokenType: auth.TokenTypeAccess,
		})
		if err != nil || !result.IsValid {
			// Invalid token, continue without authentication
			oam.logger.Debugf("invalid token, continuing without auth")
			next.ServeHTTP(w, r)
			return
		}

		// Store auth info in context
		ctx := context.WithValue(r.Context(), ContextKeyUserID, result.Claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyEmail, result.Claims.Email)
		ctx = context.WithValue(ctx, ContextKeyClaims, result.Claims)
		ctx = context.WithValue(ctx, ContextKeySessionID, result.Claims.SessionID)
		ctx = context.WithValue(ctx, ContextKeyDeviceID, result.Claims.DeviceID)

		oam.logger.Debugf("optional authentication successful for user %d", result.Claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequirePermissionMiddleware is middleware for checking permissions
type RequirePermissionMiddleware struct {
	permission string
	logger     *logger.Logger
}

// NewRequirePermissionMiddleware creates a new permission middleware
func NewRequirePermissionMiddleware(permission string, logger *logger.Logger) *RequirePermissionMiddleware {
	return &RequirePermissionMiddleware{
		permission: permission,
		logger:     logger,
	}
}

// Middleware is the HTTP middleware function for permission checking
func (rpm *RequirePermissionMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context
		userID, ok := r.Context().Value(ContextKeyUserID).(int64)
		if !ok || userID == 0 {
			rpm.logger.Debugf("user ID not found in context")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Get claims from context
		claims, ok := r.Context().Value(ContextKeyClaims).(*auth.JWTClaims)
		if !ok {
			rpm.logger.Debugf("claims not found in context")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Check if permission is in claims
		hasPermission := false
		for _, p := range claims.Permissions {
			if p == rpm.permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			rpm.logger.Warnf("permission denied for user %d: missing %s", userID, rpm.permission)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		rpm.logger.Debugf("permission check passed for user %d: %s", userID, rpm.permission)

		next.ServeHTTP(w, r)
	})
}

// ===== Helper Functions =====

// extractBearerToken extracts the bearer token from the Authorization header
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header not found")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
}

// getClientIP gets the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (proxy)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(ContextKeyUserID).(int64)
	if !ok || userID == 0 {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetEmailFromContext retrieves email from context
func GetEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(ContextKeyEmail).(string)
	if !ok || email == "" {
		return "", fmt.Errorf("email not found in context")
	}
	return email, nil
}

// GetClaimsFromContext retrieves JWT claims from context
func GetClaimsFromContext(ctx context.Context) (*auth.JWTClaims, error) {
	claims, ok := ctx.Value(ContextKeyClaims).(*auth.JWTClaims)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}

// GetSessionIDFromContext retrieves session ID from context
func GetSessionIDFromContext(ctx context.Context) (string, error) {
	sessionID, ok := ctx.Value(ContextKeySessionID).(string)
	if !ok || sessionID == "" {
		return "", fmt.Errorf("session ID not found in context")
	}
	return sessionID, nil
}
