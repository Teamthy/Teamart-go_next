package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/merchant"
	"github.com/teamart/commerce-api/internal/middleware"
	"github.com/teamart/commerce-api/internal/staff"
	"github.com/teamart/commerce-api/internal/tenant"
	"github.com/teamart/commerce-api/pkg/logger"
)

type MerchantHandler struct {
	merchantSvc *merchant.Service
	staffSvc    *staff.Service
	tenantSvc   *tenant.Service
	log         *logger.Logger
}

func NewMerchantHandler(merchantSvc *merchant.Service, staffSvc *staff.Service, tenantSvc *tenant.Service, log *logger.Logger) *MerchantHandler {
	return &MerchantHandler{merchantSvc: merchantSvc, staffSvc: staffSvc, tenantSvc: tenantSvc, log: log}
}

type CreateMerchantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	BillingPlan string `json:"billing_plan,omitempty"`
	Currency    string `json:"currency,omitempty"`
}

type CreateMerchantStoreRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Category    string          `json:"category,omitempty"`
	BannerURL   string          `json:"banner_url,omitempty"`
	Settings    json.RawMessage `json:"settings,omitempty"`
	CreatorID   *int64          `json:"creator_id,omitempty"`
}

type AddStaffRequest struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role,omitempty"`
}

func (h *MerchantHandler) HandleCreateMerchant(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.log.Warnf("unauthorized merchant create: %v", err)
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("invalid request body: %v", err)
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	merchantObj, err := h.merchantSvc.CreateMerchant(r.Context(), userID, &merchant.CreateMerchantInput{
		Name:        req.Name,
		Description: req.Description,
		BillingPlan: req.BillingPlan,
		Currency:    req.Currency,
	})
	if err != nil {
		h.log.Errorf("failed to create merchant: %v", err)
		h.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusCreated, merchantObj)
}

func (h *MerchantHandler) HandleGetMerchant(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := parseIDFromVars(r, "merchant_id")
	if !ok {
		h.respondError(w, "Missing merchant_id", http.StatusBadRequest)
		return
	}

	if !h.authorizeMerchantAccess(r, merchantID) {
		h.respondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	result, err := h.merchantSvc.GetMerchant(r.Context(), merchantID)
	if err != nil {
		h.respondError(w, "Merchant not found", http.StatusNotFound)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *MerchantHandler) HandleCreateStore(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := parseIDFromVars(r, "merchant_id")
	if !ok {
		h.respondError(w, "Missing merchant_id", http.StatusBadRequest)
		return
	}

	if !h.authorizeMerchantAccess(r, merchantID) {
		h.respondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateMerchantStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	store, err := h.merchantSvc.CreateStore(r.Context(), merchantID, userID, &merchant.CreateStoreInput{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		BannerURL:   req.BannerURL,
		Settings:    req.Settings,
		CreatorID:   req.CreatorID,
	})
	if err != nil {
		h.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusCreated, store)
}

func (h *MerchantHandler) HandleListStores(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := parseIDFromVars(r, "merchant_id")
	if !ok {
		h.respondError(w, "Missing merchant_id", http.StatusBadRequest)
		return
	}

	if !h.authorizeMerchantAccess(r, merchantID) {
		h.respondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	stores, err := h.merchantSvc.ListStoresForMerchant(r.Context(), merchantID)
	if err != nil {
		h.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"stores": stores})
}

func (h *MerchantHandler) HandleAddStaff(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := parseIDFromVars(r, "merchant_id")
	if !ok {
		h.respondError(w, "Missing merchant_id", http.StatusBadRequest)
		return
	}

	if !h.authorizeMerchantAccess(r, merchantID) {
		h.respondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req AddStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, err := h.staffSvc.CreateStaffAccount(r.Context(), merchantID, req.UserID, req.Role)
	if err != nil {
		h.respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusCreated, account)
}

func (h *MerchantHandler) HandleListStaff(w http.ResponseWriter, r *http.Request) {
	merchantID, ok := parseIDFromVars(r, "merchant_id")
	if !ok {
		h.respondError(w, "Missing merchant_id", http.StatusBadRequest)
		return
	}

	if !h.authorizeMerchantAccess(r, merchantID) {
		h.respondError(w, "Forbidden", http.StatusForbidden)
		return
	}

	staffMembers, err := h.staffSvc.ListStaffForMerchant(r.Context(), merchantID)
	if err != nil {
		h.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{"staff": staffMembers})
}

func (h *MerchantHandler) authorizeMerchantAccess(r *http.Request, merchantID int64) bool {
	claims, err := middleware.GetClaimsFromContext(r.Context())
	if err != nil {
		return false
	}

	merchantObj, err := h.merchantSvc.GetMerchant(r.Context(), merchantID)
	if err == nil && merchantObj.OwnerID == claims.UserID {
		return true
	}

	staffAccount, err := h.staffSvc.GetActiveStaffAccountByUserID(r.Context(), claims.UserID)
	if err == nil && staffAccount.MerchantID == merchantID {
		return true
	}

	return false
}

func parseIDFromVars(r *http.Request, key string) (int64, bool) {
	vars := mux.Vars(r)
	value, ok := vars[key]
	if !ok || value == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *MerchantHandler) respondError(w http.ResponseWriter, message string, code int) {
	h.log.Warnf("merchant handler error: %s", message)
	respondJSON(w, code, map[string]any{"error": message})
}

func RegisterMerchantRoutes(mux Router, handler *MerchantHandler) {
	mux.HandleFunc("POST /api/v1/merchants", handler.HandleCreateMerchant)
	mux.HandleFunc("GET /api/v1/merchants/{merchant_id}", handler.HandleGetMerchant)
	mux.HandleFunc("POST /api/v1/merchants/{merchant_id}/stores", handler.HandleCreateStore)
	mux.HandleFunc("GET /api/v1/merchants/{merchant_id}/stores", handler.HandleListStores)
	mux.HandleFunc("POST /api/v1/merchants/{merchant_id}/staff", handler.HandleAddStaff)
	mux.HandleFunc("GET /api/v1/merchants/{merchant_id}/staff", handler.HandleListStaff)
}
