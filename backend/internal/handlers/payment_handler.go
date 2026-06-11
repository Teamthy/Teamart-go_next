package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/teamart/commerce-api/internal/payments"
	"github.com/teamart/commerce-api/pkg/logger"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	paymentService *payments.PaymentService
	logger         *logger.Logger
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(ps *payments.PaymentService, logger *logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: ps,
		logger:         logger,
	}
}

// CreatePaymentIntent handles POST /api/payments/intents
// Request body: CreatePaymentIntentRequest
// Response: CreatePaymentIntentResponse
func (h *PaymentHandler) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/payments/intents")

	var req payments.CreatePaymentIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	resp, err := h.paymentService.CreatePaymentIntent(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to create payment intent", err)
		return
	}

	h.writeJSON(w, http.StatusCreated, resp)
}

// ConfirmPayment handles POST /api/payments/{payment_intent_id}/confirm
// Request body: ConfirmPaymentRequest
// Response: ConfirmPaymentResponse
func (h *PaymentHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/payments/{payment_intent_id}/confirm")

	paymentIntentID := mux.Vars(r)["payment_intent_id"]
	if paymentIntentID == "" {
		h.writeError(w, http.StatusBadRequest, "payment_intent_id is required", nil)
		return
	}

	var req payments.ConfirmPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	req.PaymentIntentID = paymentIntentID

	resp, err := h.paymentService.ConfirmPayment(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to confirm payment", err)
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// RefundPayment handles POST /api/payments/{payment_intent_id}/refund
// Request body: RefundRequest
// Response: RefundResponse
func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/payments/{payment_intent_id}/refund")

	paymentIntentID := mux.Vars(r)["payment_intent_id"]
	if paymentIntentID == "" {
		h.writeError(w, http.StatusBadRequest, "payment_intent_id is required", nil)
		return
	}

	var req payments.RefundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	req.PaymentIntentID = paymentIntentID

	resp, err := h.paymentService.RefundPayment(r.Context(), &req)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Failed to refund payment", err)
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// GetPayment handles GET /api/payments/{payment_intent_id}
// Response: GetPaymentResponse
func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/payments/{payment_intent_id}")

	paymentIntentID := mux.Vars(r)["payment_intent_id"]
	if paymentIntentID == "" {
		h.writeError(w, http.StatusBadRequest, "payment_intent_id is required", nil)
		return
	}

	resp, err := h.paymentService.GetPayment(r.Context(), paymentIntentID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Payment not found", err)
		return
	}

	h.writeJSON(w, http.StatusOK, resp)
}

// GetOrderPayments handles GET /api/orders/{order_id}/payments
// Response: []GetPaymentResponse
func (h *PaymentHandler) GetOrderPayments(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GET /api/orders/{order_id}/payments")

	orderIDStr := mux.Vars(r)["order_id"]
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid order_id", err)
		return
	}

	// TODO: Query all payments for order
	// For now return empty array
	payments := []interface{}{}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"order_id": orderID,
		"payments": payments,
	})
}

// HandleWebhook handles POST /api/payments/webhooks/{provider}
// This is called by payment providers to update payment status
func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("POST /api/payments/webhooks/{provider}")

	provider := mux.Vars(r)["provider"]
	if provider == "" {
		h.writeError(w, http.StatusBadRequest, "provider is required", nil)
		return
	}

	var payload json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid webhook payload", err)
		return
	}

	// Map provider name to event type
	eventType := provider + "_webhook"

	err := h.paymentService.HandleWebhook(r.Context(), eventType, payload)
	if err != nil {
		h.logger.Errorf("failed to handle webhook: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return 200 to acknowledge receipt
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "received",
	})
}

// Helper functions

func (h *PaymentHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *PaymentHandler) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		h.logger.Errorf("%s: %v", message, err)
	} else {
		h.logger.Error(message)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   message,
		"details": errorString(err),
	})
}

// RegisterPaymentRoutes registers payment routes on the provided router
func RegisterPaymentRoutes(r Router, h *PaymentHandler) {
	r.HandleFunc("/api/payments/intents", h.CreatePaymentIntent)
	r.HandleFunc("/api/payments/{payment_intent_id}/confirm", h.ConfirmPayment)
	r.HandleFunc("/api/payments/{payment_intent_id}/refund", h.RefundPayment)
	r.HandleFunc("/api/payments/{payment_intent_id}", h.GetPayment)
	r.HandleFunc("/api/orders/{order_id}/payments", h.GetOrderPayments)
	r.HandleFunc("/api/payments/webhooks/{provider}", h.HandleWebhook)
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
