package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/teamart/commerce-api/internal/checkout"
	"github.com/teamart/commerce-api/pkg/logger"
)

// CheckoutHandler handles checkout-related HTTP requests
type CheckoutHandler struct {
	checkoutService *checkout.CheckoutService
	logger          *logger.Logger
}

// NewCheckoutHandler creates a new checkout handler
func NewCheckoutHandler(cs *checkout.CheckoutService, logger *logger.Logger) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: cs,
		logger:          logger,
	}
}

// CalculateCheckout handles POST /api/checkout/calculate
// Request body: CalculateCheckoutRequest
// Response: CalculateCheckoutResponse
func (h *CheckoutHandler) CalculateCheckout(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/checkout/calculate")

	var req checkout.CalculateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.checkoutService.CalculateCheckout(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to calculate checkout", err)
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// CreateOrder handles POST /api/checkout/order
// Request body: CreateOrderRequest
// Response: CreateOrderResponse
func (h *CheckoutHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/checkout/order")

	var req checkout.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.checkoutService.CreateOrder(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to create order", err)
		return
	}

	h.writeJSON(w, http.StatusCreated, resp)
}

// ValidateAddress handles POST /api/checkout/validate-address
// Request body: AddressInput
// Response: {valid: bool, errors: []}
func (h *CheckoutHandler) ValidateAddress(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/checkout/validate-address")

	var addr checkout.AddressInput
	if err := json.NewDecoder(r.Body).Decode(&addr); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	err := h.checkoutService.ValidateShippingAddress(r.Context(), &addr)
	valid := err == nil

	response := map[string]interface{}{
		"valid": valid,
		"error": nil,
	}
	if err != nil {
		response["error"] = err.Error()
	}

	h.writeJSON(w, http.StatusOK, response)
}

// ApplyDiscount handles POST /api/checkout/apply-discount
// Request body: {code: string, subtotal_cents: int64}
// Response: AppliedDiscount
func (h *CheckoutHandler) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/checkout/apply-discount")

	var req struct {
		Code        string `json:"code"`
		SubtotalCents int64  `json:"subtotal_cents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.checkoutService.ApplyDiscountCode(r.Context(), req.Code, req.SubtotalCents)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to apply discount", err)
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// Helper functions

func (h *CheckoutHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *CheckoutHandler) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		h.logger.Errorf("%s: %v", message, err)
	} else {
		h.logger.Error(message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	details := ""
	if err != nil {
		details = err.Error()
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   message,
		"details": details,
	})
}

// RegisterCheckoutRoutes registers checkout routes on the provided router
func RegisterCheckoutRoutes(r Router, h *CheckoutHandler) {
	r.HandleFunc("/api/checkout/calculate", h.CalculateCheckout)
	r.HandleFunc("/api/checkout/order", h.CreateOrder)
	r.HandleFunc("/api/checkout/validate-address", h.ValidateAddress)
	r.HandleFunc("/api/checkout/apply-discount", h.ApplyDiscount)
}
