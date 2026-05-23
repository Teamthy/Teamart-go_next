package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
// Ensures users can only access resources under their own tenant boundary.
func (t *TenantIsolationMiddleware) ValidateTenantAccess() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			claims, err := GetClaimsFromContext(req.Context())
			if err != nil {
				t.logger.Warnf("unauthorized: no claims in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tenantID := claims.UserID
			ctx := context.WithValue(req.Context(), "tenant_id", tenantID)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// EnforceStoreOwnership validates that user owns the store they're accessing
func (t *TenantIsolationMiddleware) EnforceStoreOwnership() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			claims, err := GetClaimsFromContext(req.Context())
			if err != nil {
				t.logger.Warnf("unauthorized: no claims in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			vars := mux.Vars(req)
			storeID, ok := vars["store_id"]
			if !ok || storeID == "" {
				t.logger.Warnf("missing store_id parameter")
				respondError(w, "Bad Request: missing store_id", http.StatusBadRequest)
				return
			}

			if _, err := strconv.ParseInt(storeID, 10, 64); err != nil {
				t.logger.Warnf("invalid store_id format: %s", storeID)
				respondError(w, "Bad Request: invalid store_id", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(req.Context(), "store_id", storeID)
			ctx = context.WithValue(ctx, "resource_owner_id", claims.UserID)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

// EnforceCreatorOwnership validates that user owns the creator account they're accessing
func (t *TenantIsolationMiddleware) EnforceCreatorOwnership() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			claims, err := GetClaimsFromContext(req.Context())
			if err != nil {
				t.logger.Warnf("unauthorized: no claims in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			vars := mux.Vars(req)
			creatorID, ok := vars["creator_id"]
			if !ok || creatorID == "" {
				t.logger.Warnf("missing creator_id parameter")
				respondError(w, "Bad Request: missing creator_id", http.StatusBadRequest)
				return
			}

			if _, err := strconv.ParseInt(creatorID, 10, 64); err != nil {
				t.logger.Warnf("invalid creator_id format: %s", creatorID)
				respondError(w, "Bad Request: invalid creator_id", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(req.Context(), "creator_id", creatorID)
			ctx = context.WithValue(ctx, "resource_owner_id", claims.UserID)
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
			claims, err := GetClaimsFromContext(req.Context())
			if err != nil {
				t.logger.Warnf("unauthorized: no claims in context")
				respondError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			vars := mux.Vars(req)
			resourceID := vars["resource_id"]
			if resourceID == "" {
				resourceID = vars["id"]
			}

			if resourceID == "" {
				t.logger.Warnf("missing resource ID parameter for type: %s", resourceType)
				respondError(w, "Bad Request: missing resource ID", http.StatusBadRequest)
				return
			}

			hasAccess, err := accessCheckFn(claims.UserID, resourceID)
			if err != nil {
				t.logger.Errorf("error checking resource access: %v", err)
				respondError(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !hasAccess {
				t.logger.Warnf("forbidden: user %d accessing %s %s", claims.UserID, resourceType, resourceID)
				respondError(w, fmt.Sprintf("Forbidden: you cannot access this %s", resourceType), http.StatusForbidden)
				return
			}

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
