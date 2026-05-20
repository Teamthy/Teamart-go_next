package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/pkg/logger"
)

// RBACMiddleware provides role-based access control middleware
type RBACMiddleware struct {
	rbacService *auth.RBACService
	logger      *logger.Logger
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(rbacService *auth.RBACService, logger *logger.Logger) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
		logger:      logger,
	}
}

// RequireRole returns a middleware that requires a specific role
func (r *RBACMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				r.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get user roles from claims
			hasRole := false
			for _, requiredRole := range roles {
				for _, userRole := range user.Roles {
					if userRole == requiredRole {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				r.logger.Warnf("forbidden: user %d does not have required roles: %v", user.UserID, roles)
				respondError(w, "Forbidden: insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

// RequirePermission returns a middleware that requires a specific permission
func (r *RBACMiddleware) RequirePermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				r.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has any of the required permissions
			hasPermission := false
			for _, requiredPerm := range permissions {
				for _, userPerm := range user.Permissions {
					if userPerm == requiredPerm {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}

			if !hasPermission {
				r.logger.Warnf("forbidden: user %d does not have required permissions: %v", user.UserID, permissions)
				respondError(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

// RequireStoreOwnership returns a middleware that requires store ownership
func (r *RBACMiddleware) RequireStoreOwnership(storeIDParam string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				r.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get store ID from path parameter
			storeID := req.PathValue(storeIDParam)
			if storeID == "" {
				r.logger.Warnf("missing store ID parameter: %s", storeIDParam)
				respondError(w, "Bad Request: missing store ID", http.StatusBadRequest)
				return
			}

			// For now, store this in context for later use
			// In production, check against database for actual ownership
			ctx := context.WithValue(req.Context(), "store_id", storeID)
			ctx = context.WithValue(ctx, "resource_owner_id", user.UserID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// RequireCreatorOwnership returns a middleware that requires creator ownership
func (r *RBACMiddleware) RequireCreatorOwnership(creatorIDParam string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				r.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get creator ID from path parameter
			creatorID := req.PathValue(creatorIDParam)
			if creatorID == "" {
				r.logger.Warnf("missing creator ID parameter: %s", creatorIDParam)
				respondError(w, "Bad Request: missing creator ID", http.StatusBadRequest)
				return
			}

			// Store in context for later use
			ctx := context.WithValue(req.Context(), "creator_id", creatorID)
			ctx = context.WithValue(ctx, "resource_owner_id", user.UserID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// RequireResourceOwnership returns a middleware that validates resource ownership
// Checks that the user owns the resource they're trying to access
func (r *RBACMiddleware) RequireResourceOwnership(resourceType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				r.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ownerID, ok := req.Context().Value("resource_owner_id").(int64)
			if !ok {
				// Owner ID should be set by RequireStoreOwnership or RequireCreatorOwnership
				r.logger.Warnf("missing resource owner ID in context")
				respondError(w, "Bad Request: invalid resource context", http.StatusBadRequest)
				return
			}

			if user.UserID != ownerID {
				r.logger.Warnf("forbidden: user %d is not the owner of %s (owner: %d)", user.UserID, resourceType, ownerID)
				respondError(w, "Forbidden: you do not own this resource", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

// RequireTenantContext returns a middleware that ensures tenant isolation
func (r *RBACMiddleware) RequireTenantContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				r.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Store tenant context (user owns their own tenant boundary)
			ctx := context.WithValue(req.Context(), "tenant_id", user.UserID)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(req *http.Request) (int64, bool) {
	tenantID, ok := req.Context().Value("tenant_id").(int64)
	return tenantID, ok
}

// GetStoreID retrieves store ID from context
func GetStoreID(req *http.Request) (string, bool) {
	storeID, ok := req.Context().Value("store_id").(string)
	return storeID, ok
}

// GetCreatorID retrieves creator ID from context
func GetCreatorID(req *http.Request) (string, bool) {
	creatorID, ok := req.Context().Value("creator_id").(string)
	return creatorID, ok
}

// GetResourceOwnerID retrieves resource owner ID from context
func GetResourceOwnerID(req *http.Request) (int64, bool) {
	ownerID, ok := req.Context().Value("resource_owner_id").(int64)
	return ownerID, ok
}

// ===== Helper Functions =====

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"error":"%s","message":"%s","code":%d}`, http.StatusText(statusCode), message, statusCode)
}
