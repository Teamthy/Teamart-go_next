package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teamart/commerce-api/internal/auth"
	"github.com/teamart/commerce-api/pkg/logger"
)

func TestRequirePermission_AllowsWithPermission(t *testing.T) {
	log := logger.NewNoop()
	rpm := NewRequirePermissionMiddleware("admin:access", log)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	// Set context with user id and claims containing the permission
	ctx := req.Context()
	ctx = contextWithUserAndClaims(ctx, 42, &auth.JWTClaims{Permissions: []string{"admin:access"}})
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	rpm.Middleware(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rr.Code)
	}
}

func TestRequirePermission_DeniesWithoutPermission(t *testing.T) {
	log := logger.NewNoop()
	rpm := NewRequirePermissionMiddleware("admin:access", log)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	ctx := req.Context()
	ctx = contextWithUserAndClaims(ctx, 42, &auth.JWTClaims{Permissions: []string{"user:read"}})
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	rpm.Middleware(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rr.Code)
	}
}

func TestRequirePermission_DeniesWhenNoAuth(t *testing.T) {
	log := logger.NewNoop()
	rpm := NewRequirePermissionMiddleware("admin:access", log)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	// No auth info in context

	rr := httptest.NewRecorder()
	rpm.Middleware(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", rr.Code)
	}
}

// helper: put user id and claims into context using middleware keys
func contextWithUserAndClaims(ctx context.Context, userID int64, claims *auth.JWTClaims) context.Context {
	ctx = context.WithValue(ctx, ContextKeyUserID, userID)
	ctx = context.WithValue(ctx, ContextKeyClaims, claims)
	return ctx
}
