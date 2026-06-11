package payments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/pkg/logger"
)

// PaymentService handles payment intent lifecycle, provider communication,
// transaction recording, and refunds
type PaymentService struct {
	queries *queries.Queries
	logger  *logger.Logger
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(q *queries.Queries, logger *logger.Logger) *PaymentService {
	return &PaymentService{
		queries: q,
		logger:  logger,
	}
}

// ===== INPUT/OUTPUT TYPES =====

// CreatePaymentIntentRequest starts a new payment intent
type CreatePaymentIntentRequest struct {
	OrderID   int64                  `json:"order_id"`
	Amount    int64                  `json:"amount_cents"`
	Currency  string                 `json:"currency"`
	Provider  string                 `json:"provider"` // stripe, paystack, flutterwave, wallet, etc.
	Method    string                 `json:"method"`   // card, bank_transfer, wallet, mobile_money
	Metadata  map[string]interface{} `json:"metadata"`
	ReturnURL string                 `json:"return_url"`
}

// CreatePaymentIntentResponse returns payment details
type CreatePaymentIntentResponse struct {
	PaymentIntentID string            `json:"payment_intent_id"`
	ClientSecret    string            `json:"client_secret"`
	ProviderID      string            `json:"provider_id"`
	Status          string            `json:"status"`
	Amount          int64             `json:"amount_cents"`
	Currency        string            `json:"currency"`
	Provider        string            `json:"provider"`
	RequiresAction  bool              `json:"requires_action"`
	NextAction      *NextAction       `json:"next_action,omitempty"`
	AuthURL         string            `json:"auth_url,omitempty"`
}

// NextAction describes required next action for payment
type NextAction struct {
	Type        string `json:"type"` // use_stripe_sdk, redirect_to_url, display_form, etc.
	RedirectURL string `json:"redirect_url,omitempty"`
}

// ConfirmPaymentRequest confirms a payment with provider details
type ConfirmPaymentRequest struct {
	PaymentIntentID   string                 `json:"payment_intent_id"`
	ProviderToken     string                 `json:"provider_token"`
	BillingDetails    *BillingDetails        `json:"billing_details"`
	ExtraData         map[string]interface{} `json:"extra_data"`
}

// BillingDetails represents billing information
type BillingDetails struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZIP     string `json:"zip"`
	Country string `json:"country"`
}

// ConfirmPaymentResponse returns confirmation result
type ConfirmPaymentResponse struct {
	PaymentIntentID    string `json:"payment_intent_id"`
	TransactionID      string `json:"transaction_id"`
	Status             string `json:"status"` // succeeded, requires_action, failed
	Amount             int64  `json:"amount_cents"`
	AuthorizedAmount   int64  `json:"authorized_amount_cents"`
	CapturedAmount     int64  `json:"captured_amount_cents"`
	Message            string `json:"message"`
	FailureCode        string `json:"failure_code,omitempty"`
	FailureMessage     string `json:"failure_message,omitempty"`
}

// RefundRequest requests a refund
type RefundRequest struct {
	PaymentIntentID string `json:"payment_intent_id"`
	Amount          int64  `json:"amount_cents"`
	Reason          string `json:"reason"` // return, cancellation, fraud, etc.
	Description     string `json:"description"`
}

// RefundResponse returns refund result
type RefundResponse struct {
	RefundID       string `json:"refund_id"`
	Amount         int64  `json:"amount_cents"`
	Status         string `json:"status"` // succeeded, pending, failed
	CreatedAt      string `json:"created_at"`
	Message        string `json:"message"`
	FailureReason  string `json:"failure_reason,omitempty"`
}

// GetPaymentResponse returns full payment details
type GetPaymentResponse struct {
	PaymentIntentID string             `json:"payment_intent_id"`
	OrderID         int64              `json:"order_id"`
	Provider        string             `json:"provider"`
	Method          string             `json:"method"`
	Amount          int64              `json:"amount_cents"`
	Currency        string             `json:"currency"`
	Status          string             `json:"status"`
	AuthorizedAt    *time.Time         `json:"authorized_at,omitempty"`
	CapturedAt      *time.Time         `json:"captured_at,omitempty"`
	FailedAt        *time.Time         `json:"failed_at,omitempty"`
	Transactions    []PaymentTransaction `json:"transactions"`
	AuditLog        []PaymentAuditEntry  `json:"audit_log"`
}

// PaymentTransaction represents a single transaction
type PaymentTransaction struct {
	TransactionID string                 `json:"transaction_id"`
	Type          string                 `json:"type"` // auth, capture, refund, void
	Amount        int64                  `json:"amount_cents"`
	Status        string                 `json:"status"`
	CreatedAt     time.Time              `json:"created_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// PaymentAuditEntry represents an audit log entry
type PaymentAuditEntry struct {
	Action    string    `json:"action"`
	Actor     string    `json:"actor_type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// ===== SERVICE METHODS =====

// CreatePaymentIntent creates a new payment intent
func (s *PaymentService) CreatePaymentIntent(ctx context.Context, req *CreatePaymentIntentRequest) (*CreatePaymentIntentResponse, error) {
	s.logger.Debugf("creating payment intent: order_id=%d amount_cents=%d provider=%s", req.OrderID, req.Amount, req.Provider)

	// Validate request
	if req.OrderID == 0 {
		return nil, errors.New("order_id is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if req.Provider == "" {
		return nil, errors.New("provider is required")
	}

	// TODO: Call provider adapter to create payment intent
	// For now, return placeholder response

	return &CreatePaymentIntentResponse{
		PaymentIntentID: fmt.Sprintf("pi_%d_%d", req.OrderID, time.Now().Unix()),
		ClientSecret:    fmt.Sprintf("secret_%d", req.OrderID),
		ProviderID:      "test_provider_id",
		Status:          "pending",
		Amount:          req.Amount,
		Currency:        req.Currency,
		Provider:        req.Provider,
		RequiresAction:  true,
		NextAction: &NextAction{
			Type:        "redirect_to_url",
			RedirectURL: "https://example.com/pay",
		},
	}, nil
}

// ConfirmPayment confirms a payment with provider token
func (s *PaymentService) ConfirmPayment(ctx context.Context, req *ConfirmPaymentRequest) (*ConfirmPaymentResponse, error) {
	s.logger.Debugf("confirming payment: payment_intent_id=%s", req.PaymentIntentID)

	if req.PaymentIntentID == "" {
		return nil, errors.New("payment_intent_id is required")
	}

	// TODO: Implement payment confirmation
	// 1. Get payment intent from database
	// 2. Call provider to confirm payment
	// 3. Record transaction
	// 4. Update order payment status

	return &ConfirmPaymentResponse{
		PaymentIntentID: req.PaymentIntentID,
		TransactionID:   fmt.Sprintf("txn_%d", time.Now().Unix()),
		Status:          "succeeded",
		Amount:          0, // TODO: Get from intent
	}, nil
}

// AuthorizePayment authorizes a payment (pre-capture)
func (s *PaymentService) AuthorizePayment(ctx context.Context, paymentIntentID string) error {
	s.logger.Debugf("authorizing payment: payment_intent_id=%s", paymentIntentID)

	// TODO: Implement authorization flow
	// 1. Get payment intent
	// 2. Call provider to authorize
	// 3. Update status to 'authorized'
	// 4. Record audit log

	return nil
}

// CapturePayment captures a pre-authorized payment
func (s *PaymentService) CapturePayment(ctx context.Context, paymentIntentID string, amount *int64) error {
	s.logger.Debugf("capturing payment: payment_intent_id=%s", paymentIntentID)

	// TODO: Implement capture flow
	// 1. Get payment intent
	// 2. Call provider to capture
	// 3. Update status to 'captured'
	// 4. Trigger order confirmation

	return nil
}

// RefundPayment initiates a refund for a payment
func (s *PaymentService) RefundPayment(ctx context.Context, req *RefundRequest) (*RefundResponse, error) {
	s.logger.Debugf("refunding payment: payment_intent_id=%s amount_cents=%d reason=%s", req.PaymentIntentID, req.Amount, req.Reason)

	if req.PaymentIntentID == "" {
		return nil, errors.New("payment_intent_id is required")
	}
	if req.Amount <= 0 {
		return nil, errors.New("refund amount must be greater than 0")
	}

	// TODO: Implement refund flow
	// 1. Get payment intent
	// 2. Validate refund amount <= captured amount
	// 3. Call provider to refund
	// 4. Record refund transaction
	// 5. Update order/return status

	resp := &RefundResponse{
		RefundID:  fmt.Sprintf("ref_%d", time.Now().Unix()),
		Amount:    req.Amount,
		Status:    "pending",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	// Log refund created
	s.logger.Debugf("refund created: refund_id=%s amount_cents=%d", resp.RefundID, resp.Amount)
	return resp, nil
}

// GetPayment retrieves full payment details
func (s *PaymentService) GetPayment(ctx context.Context, paymentIntentID string) (*GetPaymentResponse, error) {
	if paymentIntentID == "" {
		return nil, errors.New("payment_intent_id is required")
	}

	s.logger.Debugf("getting payment details: payment_intent_id=%s", paymentIntentID)

	// TODO: Query payment intent and related transactions/audits

	return &GetPaymentResponse{
		PaymentIntentID: paymentIntentID,
		Status:          "unknown",
		Transactions:    []PaymentTransaction{},
		AuditLog:        []PaymentAuditEntry{},
	}, nil
}

// HandleWebhook processes payment provider webhooks
func (s *PaymentService) HandleWebhook(ctx context.Context, event string, payload json.RawMessage) error {
	// TODO: Implement webhook handling
	// Different providers have different webhook formats
	// Parse payload based on provider
	// Update payment status accordingly

	s.logger.Debugf("handling webhook: event=%s", event)
	return nil
}

// RecordAudit logs an audit entry for payment action
func (s *PaymentService) RecordAudit(ctx context.Context, paymentIntentID string, action string, message string, payload map[string]interface{}) error {
	s.logger.Debugf("recording payment audit: payment_intent_id=%s action=%s", paymentIntentID, action)

	// TODO: Insert audit log via SQLC

	return nil
}

// ValidateRefundAmount checks if refund amount is valid
func (s *PaymentService) ValidateRefundAmount(ctx context.Context, paymentIntentID string, amount int64) (bool, string, error) {
	// TODO: Implement validation
	// Check captured amount >= refund amount
	// Check for existing refunds
	// Check time limits

	return true, "valid", nil
}
