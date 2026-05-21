package payments

import (
	"context"
	"fmt"
	"time"
)

// PaymentGatewayProvider defines the interface for payment gateways
type PaymentGatewayProvider interface {
	// Create a payment intent
	CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error)

	// Process a payment
	ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error)

	// Refund a payment
	RefundPayment(ctx context.Context, input *RefundPaymentInput) (*RefundResult, error)

	// Verify webhook signature
	VerifyWebhookSignature(signature, timestamp, payload string) (bool, error)

	// Get payment status
	GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatus, error)
}

// ===== PAYMENT TYPES =====

// PaymentMethod represents a saved payment method
type PaymentMethod struct {
	ID              int64
	UserID          int64
	Type            string // card, wallet, bank_account, upi
	Provider        string // stripe, razorpay, paystack
	CardLastFour    *string
	CardBrand       *string
	CardExpiryMonth *int
	CardExpiryYear  *int
	CardholderName  *string
	AccountEmail    *string
	AccountPhone    *string
	ProviderID      string
	IsDefault       bool
	IsActive        bool
	Verified        bool
	VerifiedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PaymentIntent represents a payment intent
type PaymentIntent struct {
	ID                   int64
	OrderID              int64
	UserID               int64
	Amount               float64
	Currency             string
	Status               string // pending, authorized, processing, succeeded, failed, expired
	Provider             string
	ProviderIntentID     string
	ProviderClientSecret *string
	PaymentMethodID      *int64
	RiskLevel            *string
	Requires3DSecure     bool
	ThreeDSecureStatus   *string
	Description          *string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ExpiresAt            *time.Time
	SucceededAt          *time.Time
	FailedAt             *time.Time
}

// PaymentTransaction represents a transaction
type PaymentTransaction struct {
	ID                    int64
	PaymentIntentID       int64
	OrderID               int64
	Type                  string // charge, refund, dispute
	Status                string // pending, succeeded, failed
	Amount                float64
	Currency              string
	Provider              string
	ProviderTransactionID string
	ErrorCode             *string
	ErrorMessage          *string
	CreatedAt             time.Time
	ProcessedAt           *time.Time
}

// Refund represents a refund
type Refund struct {
	ID               int64
	PaymentIntentID  int64
	OrderID          int64
	Amount           float64
	Currency         string
	Reason           string
	Status           string // pending, processing, completed, failed
	Provider         string
	ProviderRefundID *string
	RequestedBy      int64
	RequestedAt      time.Time
	ApprovedBy       *int64
	ApprovedAt       *time.Time
	ProcessedAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Wallet represents user wallet/account balance
type Wallet struct {
	ID               int64
	UserID           int64
	Balance          float64
	Currency         string
	HeldAmount       float64
	AvailableBalance float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	ID              int64
	WalletID        int64
	UserID          int64
	Type            string // deposit, withdrawal, payment, refund
	Amount          float64
	PreviousBalance float64
	NewBalance      float64
	ReferenceType   string
	ReferenceID     string
	Description     string
	Status          string // pending, completed, failed
	CreatedAt       time.Time
}

// EscrowAccount represents an escrow account
type EscrowAccount struct {
	ID          int64
	OrderID     int64
	BuyerID     int64
	SellerID    int64
	Amount      float64
	Currency    string
	Status      string // held, released, refunded, disputed
	CreatedAt   time.Time
	ReleaseDate *time.Time
	ReleasedAt  *time.Time
	RefundedAt  *time.Time
}

// ===== INPUT TYPES =====

// CreatePaymentIntentInput represents input for creating a payment intent
type CreatePaymentIntentInput struct {
	OrderID         int64
	UserID          int64
	Amount          float64
	Currency        string
	Provider        string
	PaymentMethodID *int64
	Description     string
	Metadata        map[string]string
}

// ProcessPaymentInput represents input for processing a payment
type ProcessPaymentInput struct {
	PaymentIntentID int64
	Amount          float64
	Currency        string
	Provider        string
	ProviderToken   string // Token from payment gateway
	ReturnURL       string
	ConfirmURL      string
}

// RefundPaymentInput represents input for refunding a payment
type RefundPaymentInput struct {
	PaymentIntentID int64
	Amount          float64
	Reason          string
	RequestedBy     int64
}

// CreatePaymentMethodInput represents input for creating a payment method
type CreatePaymentMethodInput struct {
	UserID        int64
	Type          string
	Provider      string
	ProviderToken string
	IsDefault     bool
}

// AddFundsInput represents input for adding funds to wallet
type AddFundsInput struct {
	UserID   int64
	Amount   float64
	Currency string
	Method   string // card, bank_transfer, mobile_money
}

// ===== OUTPUT TYPES =====

// PaymentIntentResult represents the result of creating a payment intent
type PaymentIntentResult struct {
	PaymentIntentID      int64
	Provider             string
	ProviderIntentID     string
	ProviderClientSecret *string
	ClientToken          *string // For some gateways
	AuthURL              *string
	Status               string
}

// PaymentResult represents the result of processing a payment
type PaymentResult struct {
	Success               bool
	ProviderTransactionID string
	Status                string
	Message               string
	ErrorCode             *string
	Requires3DSecure      bool
	ThreeDSecureURL       *string
}

// RefundResult represents the result of a refund
type RefundResult struct {
	Success          bool
	ProviderRefundID string
	Status           string
	Message          string
	Amount           float64
}

// PaymentStatus represents payment status
type PaymentStatus struct {
	ProviderPaymentID string
	Status            string
	Amount            float64
	Currency          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ===== WEBHOOK TYPES =====

// WebhookPayload represents a webhook event
type WebhookPayload struct {
	Provider  string
	EventType string
	Data      map[string]interface{}
	Signature string
	Timestamp int64
}

// WebhookEvent represents a processed webhook event
type WebhookEvent struct {
	ID          int64
	Provider    string
	EventType   string
	PaymentID   *string
	OrderID     *int64
	Amount      *float64
	Status      string
	ProcessedAt *time.Time
	CreatedAt   time.Time
}

// ===== SEARCH & FILTER TYPES =====

// PaymentSearchCriteria represents search criteria for payments
type PaymentSearchCriteria struct {
	UserID    *int64
	OrderID   *int64
	Status    *string
	Provider  *string
	MinAmount *float64
	MaxAmount *float64
	StartDate *time.Time
	EndDate   *time.Time
	SortBy    string
	Limit     int64
	Offset    int64
}

// RefundSearchCriteria represents search criteria for refunds
type RefundSearchCriteria struct {
	UserID    *int64
	OrderID   *int64
	Status    *string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int64
	Offset    int64
}

// ===== PAYMENT GATEWAY IMPLEMENTATIONS =====

// StripeGateway implements PaymentGatewayProvider for Stripe
type StripeGateway struct {
	secretKey string
	publicKey string
}

// NewStripeGateway creates a new Stripe gateway
func NewStripeGateway(secretKey, publicKey string) *StripeGateway {
	return &StripeGateway{
		secretKey: secretKey,
		publicKey: publicKey,
	}
}

// CreatePaymentIntent creates a Stripe payment intent
func (sg *StripeGateway) CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error) {
	// TODO: Call Stripe API to create payment intent
	// stripe.PaymentIntents.New(...)
	return nil, fmt.Errorf("not yet implemented")
}

// ProcessPayment processes a Stripe payment
func (sg *StripeGateway) ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error) {
	// TODO: Call Stripe API to confirm payment intent
	return nil, fmt.Errorf("not yet implemented")
}

// RefundPayment refunds a Stripe payment
func (sg *StripeGateway) RefundPayment(ctx context.Context, input *RefundPaymentInput) (*RefundResult, error) {
	// TODO: Call Stripe API to create refund
	return nil, fmt.Errorf("not yet implemented")
}

// VerifyWebhookSignature verifies Stripe webhook signature
func (sg *StripeGateway) VerifyWebhookSignature(signature, timestamp, payload string) (bool, error) {
	// TODO: Verify using Stripe's signing secret
	return false, fmt.Errorf("not yet implemented")
}

// GetPaymentStatus gets payment status from Stripe
func (sg *StripeGateway) GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatus, error) {
	// TODO: Call Stripe API to retrieve payment intent
	return nil, fmt.Errorf("not yet implemented")
}

// ===== RAZORPAY GATEWAY =====

// RazorpayGateway implements PaymentGatewayProvider for Razorpay
type RazorpayGateway struct {
	keyID     string
	keySecret string
}

// NewRazorpayGateway creates a new Razorpay gateway
func NewRazorpayGateway(keyID, keySecret string) *RazorpayGateway {
	return &RazorpayGateway{
		keyID:     keyID,
		keySecret: keySecret,
	}
}

// CreatePaymentIntent creates a Razorpay order (payment intent)
func (rg *RazorpayGateway) CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error) {
	// TODO: Call Razorpay API to create order
	return nil, fmt.Errorf("not yet implemented")
}

// ProcessPayment processes a Razorpay payment
func (rg *RazorpayGateway) ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error) {
	// TODO: Call Razorpay API
	return nil, fmt.Errorf("not yet implemented")
}

// RefundPayment refunds a Razorpay payment
func (rg *RazorpayGateway) RefundPayment(ctx context.Context, input *RefundPaymentInput) (*RefundResult, error) {
	// TODO: Call Razorpay API to create refund
	return nil, fmt.Errorf("not yet implemented")
}

// VerifyWebhookSignature verifies Razorpay webhook signature
func (rg *RazorpayGateway) VerifyWebhookSignature(signature, timestamp, payload string) (bool, error) {
	// TODO: Verify Razorpay signature
	return false, fmt.Errorf("not yet implemented")
}

// GetPaymentStatus gets payment status from Razorpay
func (rg *RazorpayGateway) GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatus, error) {
	// TODO: Call Razorpay API
	return nil, fmt.Errorf("not yet implemented")
}

// ===== PAYSTACK GATEWAY =====

// PaystackGateway implements PaymentGatewayProvider for Paystack
type PaystackGateway struct {
	secretKey string
	publicKey string
}

// NewPaystackGateway creates a new Paystack gateway
func NewPaystackGateway(secretKey, publicKey string) *PaystackGateway {
	return &PaystackGateway{
		secretKey: secretKey,
		publicKey: publicKey,
	}
}

// CreatePaymentIntent creates a Paystack transaction
func (pg *PaystackGateway) CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error) {
	// TODO: Call Paystack API
	return nil, fmt.Errorf("not yet implemented")
}

// ProcessPayment processes a Paystack payment
func (pg *PaystackGateway) ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error) {
	// TODO: Call Paystack API
	return nil, fmt.Errorf("not yet implemented")
}

// RefundPayment refunds a Paystack payment
func (pg *PaystackGateway) RefundPayment(ctx context.Context, input *RefundPaymentInput) (*RefundResult, error) {
	// TODO: Call Paystack API
	return nil, fmt.Errorf("not yet implemented")
}

// VerifyWebhookSignature verifies Paystack webhook signature
func (pg *PaystackGateway) VerifyWebhookSignature(signature, timestamp, payload string) (bool, error) {
	// TODO: Verify Paystack signature
	return false, fmt.Errorf("not yet implemented")
}

// GetPaymentStatus gets payment status from Paystack
func (pg *PaystackGateway) GetPaymentStatus(ctx context.Context, providerPaymentID string) (*PaymentStatus, error) {
	// TODO: Call Paystack API
	return nil, fmt.Errorf("not yet implemented")
}

// ===== SPLIT PAYMENTS =====

// SplitPayment represents a payment that's split between multiple recipients
type SplitPayment struct {
	ID              int64
	PaymentIntentID int64
	OrderID         int64
	Status          string // pending, completed, failed
	Splits          []*SplitPaymentLine
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SplitPaymentLine represents a single line item in a split payment
type SplitPaymentLine struct {
	ID            int64
	SplitID       int64
	RecipientID   int64  // seller, platform, affiliate, etc
	RecipientType string // seller, platform, affiliate, livestream_host, creator
	Amount        float64
	Currency      string
	Percentage    float64
	Status        string // pending, completed, failed
	ProviderID    *string
	CreatedAt     time.Time
}

// CreateSplitPaymentInput represents input for creating a split payment
type CreateSplitPaymentInput struct {
	OrderID     int64
	TotalAmount float64
	Currency    string
	Splits      []*SplitLine
}

// SplitLine represents a single split
type SplitLine struct {
	RecipientID   int64
	RecipientType string // seller, platform, affiliate, livestream_host
	Amount        float64
	Percentage    float64 // Alternative to amount
}

// ===== ESCROW DISPUTES =====

// EscrowDispute represents a dispute in an escrow transaction
type EscrowDispute struct {
	ID              int64
	EscrowAccountID int64
	InitiatedBy     int64
	InitiatedAt     time.Time
	Reason          string
	Status          string // open, investigating, resolved, closed
	ResolvedBy      *int64
	ResolvedAt      *time.Time
	Resolution      *string
	Outcome         *string // buyer, seller, split
	BuyerAmount     *float64
	SellerAmount    *float64
	Evidence        *string // URLs or text evidence
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ===== PAYOUT ENGINE =====

// Payout represents a payout to a seller/creator
type Payout struct {
	ID                  int64
	SellerID            int64
	Amount              float64
	Currency            string
	Status              string // pending, approved, processing, completed, failed, reversed
	PayoutMethod        string // bank_transfer, mobile_money, paypal, crypto, wallet
	PaymentMethodID     int64
	PeriodStart         *time.Time
	PeriodEnd           *time.Time
	Provider            string
	ProviderPayoutID    *string
	ReviewedBy          *int64
	ReviewedAt          *time.Time
	ReviewNotes         *string
	ScheduledAt         *time.Time
	ProcessingStartedAt *time.Time
	CompletedAt         *time.Time
	FailureReason       *string
	Metadata            map[string]interface{}
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// PayoutSchedule represents a payout schedule for a seller
type PayoutSchedule struct {
	ID            int64
	SellerID      int64
	Frequency     string // daily, weekly, biweekly, monthly
	NextPayoutAt  time.Time
	IsActive      bool
	MinimumAmount float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CreatePayoutInput represents input for creating a payout
type CreatePayoutInput struct {
	SellerID     int64
	Amount       float64
	Currency     string
	PayoutMethod string
	MethodID     int64
	PeriodStart  time.Time
	PeriodEnd    time.Time
}

// ===== MERCHANT WALLET SYSTEM =====

// MerchantWallet represents a seller/merchant's wallet
type MerchantWallet struct {
	ID               int64
	SellerID         int64
	Balance          float64 // total balance
	Currency         string
	PendingBalance   float64 // awaiting release from escrow
	AvailableBalance float64 // ready for withdrawal
	TotalEarned      float64 // lifetime earnings
	TotalWithdrawn   float64 // lifetime withdrawals
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// MerchantWalletTransaction represents a transaction in merchant wallet
type MerchantWalletTransaction struct {
	ID              int64
	WalletID        int64
	SellerID        int64
	Type            string // sale, refund, payout, adjustment, reversal, commission
	Amount          float64
	PreviousBalance float64
	NewBalance      float64
	ReferenceType   string // order, refund, payout, dispute
	ReferenceID     string
	Description     string
	Status          string // completed, pending, failed
	CreatedAt       time.Time
}

// ===== RECONCILIATION =====

// PaymentReconciliation represents a reconciliation record
type PaymentReconciliation struct {
	ID                int64
	Provider          string
	ReportDate        time.Time
	PeriodStart       time.Time
	PeriodEnd         time.Time
	ReceivedAmount    float64
	ProcessedAmount   float64
	DiscrepancyAmount float64
	Discrepancies     []*ReconciliationDiscrepancy
	Status            string // pending, verified, flagged, adjusted
	ReviewedBy        *int64
	ReviewedAt        *time.Time
	Notes             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ReconciliationDiscrepancy represents a discrepancy in reconciliation
type ReconciliationDiscrepancy struct {
	ID               int64
	ReconciliationID int64
	PaymentIntentID  *int64
	Amount           float64
	DiscrepancyType  string // missing, extra, amount_mismatch, duplicate
	Status           string // open, investigating, resolved
	Resolution       *string
	CreatedAt        time.Time
}

// ===== APPLE PAY & GOOGLE PAY =====

// ApplePayToken represents an Apple Pay token
type ApplePayToken struct {
	Token       string // Encrypted token from Apple
	LastFour    string
	Brand       string
	ExpiryMonth int
	ExpiryYear  int
}

// GooglePayToken represents a Google Pay token
type GooglePayToken struct {
	Token       string // Encrypted token from Google
	LastFour    string
	Brand       string
	ExpiryMonth int
	ExpiryYear  int
}

// ===== PAYMENT ANALYTICS =====

// PaymentMetrics represents payment metrics for a time period
type PaymentMetrics struct {
	Period             string // daily, weekly, monthly
	Date               time.Time
	TotalTransactions  int64
	SuccessfulPayments int64
	FailedPayments     int64
	TotalAmount        float64
	AvgTransactionSize float64
	SuccessRate        float64 // percentage
	DeclinedRate       float64 // percentage
	PaymentsByProvider map[string]int64
	PaymentsByGateway  map[string]float64
}

// ===== ERROR CODES =====

// PaymentError represents payment error details
type PaymentError struct {
	Code      string
	Message   string
	Retryable bool
}
