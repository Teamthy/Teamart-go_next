package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/teamart/commerce-api/internal/merchant"
	"github.com/teamart/commerce-api/internal/middleware"
	"github.com/teamart/commerce-api/internal/staff"
	"github.com/teamart/commerce-api/internal/tenant"
	"github.com/teamart/commerce-api/pkg/logger"
)

type TenantHandler struct {
	merchantSvc *merchant.Service
	staffSvc    *staff.Service
	tenantSvc   *tenant.Service
	log         *logger.Logger
}

func NewTenantHandler(merchantSvc *merchant.Service, staffSvc *staff.Service, tenantSvc *tenant.Service, log *logger.Logger) *TenantHandler {
	return &TenantHandler{merchantSvc: merchantSvc, staffSvc: staffSvc, tenantSvc: tenantSvc, log: log}
}

type UpsertTenantSettingRequest struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

func (h *TenantHandler) HandleGetTenantSettings(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := parseIDFromVars(r, "tenant_id")
	if !ok {
		h.respondError(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	if !h.authorizeTenantAccess(r, tenantID) {
		h.respondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	settings, err := h.tenantSvc.ListSettings(r.Context(), tenantID)
	if err != nil {
		h.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"settings": settings})
}

func (h *TenantHandler) HandleUpsertTenantSetting(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := parseIDFromVars(r, "tenant_id")
	if !ok {
		h.respondError(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	if !h.authorizeTenantAccess(r, tenantID) {
		h.respondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req UpsertTenantSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Key == "" {
		h.respondError(w, "Setting key is required", http.StatusBadRequest)
		return
	}

	setting, err := h.tenantSvc.UpsertSetting(r.Context(), tenantID, req.Key, req.Value)
	if err != nil {
		h.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, http.StatusOK, setting)
}

func (h *TenantHandler) authorizeTenantAccess(r *http.Request, tenantID int64) bool {
	claims, err := middleware.GetClaimsFromContext(r.Context())
	if err != nil {
		return false
	}

	merchantObj, err := h.merchantSvc.GetMerchant(r.Context(), tenantID)
	if err == nil && merchantObj.OwnerID == claims.UserID {
		return true
	}

	staffAccount, err := h.staffSvc.GetActiveStaffAccountByUserID(r.Context(), claims.UserID)
	if err == nil && staffAccount.MerchantID == tenantID {
		return true
	}

	return false
}

func RegisterTenantRoutes(mux Router, handler *TenantHandler) {
	mux.HandleFunc("GET /api/v1/tenants/{tenant_id}/settings", handler.HandleGetTenantSettings)
	mux.HandleFunc("PUT /api/v1/tenants/{tenant_id}/settings", handler.HandleUpsertTenantSetting)
}

func (h *TenantHandler) respondError(w http.ResponseWriter, message string, code int) {
	h.log.Warnf("tenant handler error: %s", message)
	respondJSON(w, code, map[string]any{"error": message})
}
