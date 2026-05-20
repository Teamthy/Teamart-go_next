package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/pkg/logger"
)

// TenantIsolationMiddleware ensures multi-tenant data isolation
type TenantIsolationMiddleware struct {
	logger *logger.Logger
}

// NewTenantIsolationMiddleware creates a new tenant isolation middleware
func NewTenantIsolationMiddleware(logger *logger.Logger) *TenantIsolationMiddleware {
	return &TenantIsolationMiddleware{
		logger: logger,
	}
}

// ValidateTenantAccess returns middleware that validates tenant access
// Ensures users can only access their own tenant's resources
func (t *TenantIsolationMiddleware) ValidateTenantAccess() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Get user from context (should be set by auth middleware)
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				t.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// User ID becomes the tenant ID (each user is their own tenant)
			tenantID := user.UserID

			// Store tenant info in context
			ctx := context.WithValue(req.Context(), "tenant_id", tenantID)
			ctx = context.WithValue(ctx, "user_id", user.UserID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// EnforceStoreOwnership validates that user owns the store they're accessing
func (t *TenantIsolationMiddleware) EnforceStoreOwnership(storeRepo interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				t.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get store ID from path
			storeID := req.PathValue("store_id")
			if storeID == "" {
				t.logger.Warnf("missing store_id parameter")
				respondError(w, "Bad Request: missing store_id", http.StatusBadRequest)
				return
			}

			// Parse store ID to int64 (assuming stores are owned by user IDs)
			// In production, query actual store ownership from database
			storeOwnerID, err := strconv.ParseInt(storeID, 10, 64)
			if err != nil {
				t.logger.Warnf("invalid store_id format: %s", storeID)
				respondError(w, "Bad Request: invalid store_id", http.StatusBadRequest)
				return
			}

			// Verify ownership
			if user.UserID != storeOwnerID {
				t.logger.Warnf("forbidden: user %d attempting to access store %d", user.UserID, storeOwnerID)
				respondError(w, "Forbidden: you do not own this store", http.StatusForbidden)
				return
			}

			// Store ownership context
			ctx := context.WithValue(req.Context(), "store_owner_id", user.UserID)
			ctx = context.WithValue(ctx, "store_id", storeID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// EnforceCreatorOwnership validates that user owns the creator account they're accessing
func (t *TenantIsolationMiddleware) EnforceCreatorOwnership() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				t.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get creator ID from path
			creatorID := req.PathValue("creator_id")
			if creatorID == "" {
				t.logger.Warnf("missing creator_id parameter")
				respondError(w, "Bad Request: missing creator_id", http.StatusBadRequest)
				return
			}

			// Parse creator ID to int64
			creatorOwnerID, err := strconv.ParseInt(creatorID, 10, 64)
			if err != nil {
				t.logger.Warnf("invalid creator_id format: %s", creatorID)
				respondError(w, "Bad Request: invalid creator_id", http.StatusBadRequest)
				return
			}

			// Verify ownership
			if user.UserID != creatorOwnerID {
				t.logger.Warnf("forbidden: user %d attempting to access creator %d", user.UserID, creatorOwnerID)
				respondError(w, "Forbidden: you do not own this creator account", http.StatusForbidden)
				return
			}

			// Store creator context
			ctx := context.WithValue(req.Context(), "creator_owner_id", user.UserID)
			ctx = context.WithValue(ctx, "creator_id", creatorID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// ValidateResourceAccess validates access to a specific resource
// Uses a callback function to determine if user can access the resource
func (t *TenantIsolationMiddleware) ValidateResourceAccess(
	resourceType string,
	accessCheckFn func(userID int64, resourceID string) (bool, error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			user, ok := req.Context().Value("user").(*auth.CustomClaims)
			if !ok || user == nil {
				t.logger.Warnf("unauthorized: no user in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get resource ID from path
			resourceID := req.PathValue("resource_id")
			if resourceID == "" {
				resourceID = req.PathValue("id")
			}

			if resourceID == "" {
				t.logger.Warnf("missing resource ID parameter for type: %s", resourceType)
				respondError(w, "Bad Request: missing resource ID", http.StatusBadRequest)
				return
			}

			// Check access
			hasAccess, err := accessCheckFn(user.UserID, resourceID)
			if err != nil {
				t.logger.Errorf("error checking resource access: %v", err)
				respondError(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !hasAccess {
				t.logger.Warnf("forbidden: user %d accessing %s %s", user.UserID, resourceType, resourceID)
				respondError(w, fmt.Sprintf("Forbidden: you cannot access this %s", resourceType), http.StatusForbidden)
				return
			}

			// Store resource context
			ctx := context.WithValue(req.Context(), "resource_type", resourceType)
			ctx = context.WithValue(ctx, "resource_id", resourceID)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// GetStoreOwnerID retrieves the store owner ID from context
func GetStoreOwnerID(req *http.Request) (int64, bool) {
	ownerID, ok := req.Context().Value("store_owner_id").(int64)
	return ownerID, ok
}

// GetCreatorOwnerID retrieves the creator owner ID from context
func GetCreatorOwnerID(req *http.Request) (int64, bool) {
	ownerID, ok := req.Context().Value("creator_owner_id").(int64)
	return ownerID, ok
}
