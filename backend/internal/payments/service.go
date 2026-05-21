package payments

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"database/sql"
)

// Service handles all payment-related operations
type Service struct {
	queries PaymentQuerier // Interface for database operations
	logger  *log.Logger

	// Gateway providers
	stripeGateway   PaymentGatewayProvider
	razorpayGateway PaymentGatewayProvider
	paystackGateway PaymentGatewayProvider
	flutterGateway  PaymentGatewayProvider
}

// PaymentQuerier interface for database operations
type PaymentQuerier interface {
	// Payment Methods
	CreatePaymentMethod(ctx context.Context, userID int64, methodType, provider, providerID string) (*PaymentMethod, error)
	GetPaymentMethod(ctx context.Context, methodID int64) (*PaymentMethod, error)
	ListPaymentMethods(ctx context.Context, userID int64) ([]*PaymentMethod, error)
	SetDefaultPaymentMethod(ctx context.Context, userID, methodID int64) error
	DeletePaymentMethod(ctx context.Context, methodID int64) error

	// Payment Intents
	CreatePaymentIntent(ctx context.Context, orderID, userID int64, amount float64, currency, provider string) (*PaymentIntent, error)
	GetPaymentIntent(ctx context.Context, intentID int64) (*PaymentIntent, error)
	UpdatePaymentIntentStatus(ctx context.Context, intentID int64, status string) (*PaymentIntent, error)
	GetPaymentIntentByProviderID(ctx context.Context, providerID string) (*PaymentIntent, error)

	// Transactions
	CreatePaymentTransaction(ctx context.Context, intentID, orderID int64, transType string, amount float64) (*PaymentTransaction, error)
	GetPaymentTransaction(ctx context.Context, txID int64) (*PaymentTransaction, error)
	ListPaymentTransactions(ctx context.Context, intentID int64) ([]*PaymentTransaction, error)
	UpdateTransactionStatus(ctx context.Context, txID int64, status, errCode, errMsg string) (*PaymentTransaction, error)

	// Webhooks
	CreateWebhookLog(ctx context.Context, provider, eventType string, payload []byte) (int64, error)
	UpdateWebhookProcessed(ctx context.Context, logID int64, success bool, errorMsg string) error

	// Wallets
	GetUserWallet(ctx context.Context, userID int64) (*Wallet, error)
	CreateUserWallet(ctx context.Context, userID int64) (*Wallet, error)
	AddFundsToWallet(ctx context.Context, walletID int64, amount float64) (*WalletTransaction, error)
	DeductFromWallet(ctx context.Context, walletID int64, amount float64) (*WalletTransaction, error)
	GetWalletBalance(ctx context.Context, userID int64) (float64, error)

	// Refunds
	CreateRefund(ctx context.Context, intentID, orderID int64, amount float64, reason string, requestedBy int64) (*Refund, error)
	GetRefund(ctx context.Context, refundID int64) (*Refund, error)
	ListRefunds(ctx context.Context, criteria *RefundSearchCriteria) ([]*Refund, error)
	ApproveRefund(ctx context.Context, refundID int64, approvedBy int64) (*Refund, error)
	UpdateRefundStatus(ctx context.Context, refundID int64, status, providerRefundID string) (*Refund, error)

	// Escrow
	CreateEscrowAccount(ctx context.Context, orderID, buyerID, sellerID int64, amount float64) (*EscrowAccount, error)
	ReleaseEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error)
	RefundEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error)
}

// NewPaymentService creates a new payment service
func NewPaymentService(queries PaymentQuerier, logger *log.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger,
	}
}

// SetGatewayProviders sets the payment gateway providers
func (s *Service) SetGatewayProviders(
	stripe, razorpay, paystack, flutter PaymentGatewayProvider,
) {
	s.stripeGateway = stripe
	s.razorpayGateway = razorpay
	s.paystackGateway = paystack
	s.flutterGateway = flutter
}

// getGateway returns the appropriate gateway provider
func (s *Service) getGateway(provider string) PaymentGatewayProvider {
	switch provider {
	case "stripe":
		return s.stripeGateway
	case "razorpay":
		return s.razorpayGateway
	case "paystack":
		return s.paystackGateway
	case "flutterwave":
		return s.flutterGateway
	default:
		return nil
	}
}

// ===== PAYMENT INTENT OPERATIONS =====

// CreatePaymentIntent creates a payment intent
func (s *Service) CreatePaymentIntent(ctx context.Context, input *CreatePaymentIntentInput) (*PaymentIntentResult, error) {
	// Validate input
	if input.OrderID == 0 || input.UserID == 0 {
		return nil, fmt.Errorf("order and user ID required")
	}
	if input.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %f", input.Amount)
	}
	if input.Provider == "" {
		input.Provider = "stripe" // Default to Stripe
	}

	// Get gateway
	gateway := s.getGateway(input.Provider)
	if gateway == nil {
		return nil, fmt.Errorf("unsupported payment provider: %s", input.Provider)
	}

	// Create payment intent with gateway
	result, err := gateway.CreatePaymentIntent(ctx, input)
	if err != nil {
		s.logger.Printf("failed to create payment intent with %s: %v", input.Provider, err)
		return nil, err
	}

	// Create database record
	dbIntent, err := s.queries.CreatePaymentIntent(ctx, input.OrderID, input.UserID, input.Amount, input.Currency, input.Provider)
	if err != nil {
		s.logger.Printf("failed to save payment intent: %v", err)
		return nil, err
	}

	result.PaymentIntentID = dbIntent.ID
	s.logger.Printf("created payment intent: %d for order %d", dbIntent.ID, input.OrderID)

	return result, nil
}

// GetPaymentIntent retrieves a payment intent
func (s *Service) GetPaymentIntent(ctx context.Context, intentID int64) (*PaymentIntent, error) {
	if intentID == 0 {
		return nil, fmt.Errorf("invalid payment intent ID")
	}

	intent, err := s.queries.GetPaymentIntent(ctx, intentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment intent not found")
		}
		return nil, err
	}

	return intent, nil
}

// UpdatePaymentIntentStatus updates payment intent status
func (s *Service) UpdatePaymentIntentStatus(ctx context.Context, intentID int64, status string) (*PaymentIntent, error) {
	if intentID == 0 {
		return nil, fmt.Errorf("invalid payment intent ID")
	}

	validStatuses := map[string]bool{
		"pending": true, "authorized": true, "processing": true,
		"succeeded": true, "failed": true, "expired": true,
	}

	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	updated, err := s.queries.UpdatePaymentIntentStatus(ctx, intentID, status)
	if err != nil {
		s.logger.Printf("failed to update payment intent %d: %v", intentID, err)
		return nil, err
	}

	s.logger.Printf("payment intent %d status updated to: %s", intentID, status)
	return updated, nil
}

// ===== PAYMENT PROCESSING =====

// ProcessPayment processes a payment
func (s *Service) ProcessPayment(ctx context.Context, input *ProcessPaymentInput) (*PaymentResult, error) {
	// Get payment intent
	intent, err := s.queries.GetPaymentIntent(ctx, input.PaymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("payment intent not found: %v", err)
	}

	if intent.Status != "pending" && intent.Status != "authorized" {
		return nil, fmt.Errorf("invalid payment intent status: %s", intent.Status)
	}

	// Get gateway
	gateway := s.getGateway(input.Provider)
	if gateway == nil {
		return nil, fmt.Errorf("unsupported payment provider: %s", input.Provider)
	}

	// Process payment
	result, err := gateway.ProcessPayment(ctx, input)
	if err != nil {
		s.logger.Printf("payment processing failed: %v", err)

		// Record failed transaction
		_, _ = s.queries.CreatePaymentTransaction(ctx, input.PaymentIntentID, intent.OrderID, "charge", input.Amount)

		return nil, err
	}

	if result.Success {
		// Record transaction
		tx, err := s.queries.CreatePaymentTransaction(ctx, input.PaymentIntentID, intent.OrderID, "charge", input.Amount)
		if err != nil {
			s.logger.Printf("failed to record payment transaction: %v", err)
			return nil, err
		}

		// Update transaction status
		_, err = s.queries.UpdateTransactionStatus(ctx, tx.ID, "succeeded", nil, nil)
		if err != nil {
			s.logger.Printf("failed to update transaction status: %v", err)
		}

		// Update payment intent status
		_, err = s.queries.UpdatePaymentIntentStatus(ctx, input.PaymentIntentID, "succeeded")
		if err != nil {
			s.logger.Printf("failed to update payment intent status: %v", err)
		}

		s.logger.Printf("payment processed successfully: %s", result.ProviderTransactionID)
	}

	return result, nil
}

// ===== REFUND OPERATIONS =====

// CreateRefund initiates a refund
func (s *Service) CreateRefund(ctx context.Context, input *RefundPaymentInput) (*Refund, error) {
	// Get payment intent
	intent, err := s.queries.GetPaymentIntent(ctx, input.PaymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("payment intent not found: %v", err)
	}

	if intent.Status != "succeeded" {
		return nil, fmt.Errorf("cannot refund payment with status: %s", intent.Status)
	}

	if input.Amount > intent.Amount {
		return nil, fmt.Errorf("refund amount exceeds payment amount")
	}

	// Create refund record
	refund, err := s.queries.CreateRefund(ctx, input.PaymentIntentID, intent.OrderID, input.Amount, input.Reason, input.RequestedBy)
	if err != nil {
		s.logger.Printf("failed to create refund: %v", err)
		return nil, err
	}

	s.logger.Printf("refund initiated: %d for amount: %.2f", refund.ID, input.Amount)
	return refund, nil
}

// ProcessRefund processes a refund through the payment gateway
func (s *Service) ProcessRefund(ctx context.Context, refundID int64) (*RefundResult, error) {
	// Get refund
	refund, err := s.queries.GetRefund(ctx, refundID)
	if err != nil {
		return nil, fmt.Errorf("refund not found: %v", err)
	}

	if refund.Status != "pending" {
		return nil, fmt.Errorf("refund status is not pending: %s", refund.Status)
	}

	// Get payment intent
	intent, err := s.queries.GetPaymentIntent(ctx, refund.PaymentIntentID)
	if err != nil {
		return nil, fmt.Errorf("payment intent not found: %v", err)
	}

	// Get gateway
	gateway := s.getGateway(intent.Provider)
	if gateway == nil {
		return nil, fmt.Errorf("unsupported payment provider: %s", intent.Provider)
	}

	// Process refund with gateway
	refundInput := &RefundPaymentInput{
		PaymentIntentID: refund.PaymentIntentID,
		Amount:          refund.Amount,
		Reason:          refund.Reason,
		RequestedBy:     refund.RequestedBy,
	}

	result, err := gateway.RefundPayment(ctx, refundInput)
	if err != nil {
		s.logger.Printf("gateway refund failed: %v", err)
		return nil, err
	}

	// Update refund status
	status := "completed"
	if !result.Success {
		status = "failed"
	}

	_, err = s.queries.UpdateRefundStatus(ctx, refundID, status, result.ProviderRefundID)
	if err != nil {
		s.logger.Printf("failed to update refund status: %v", err)
		return nil, err
	}

	s.logger.Printf("refund processed: %d with status: %s", refundID, status)
	return result, nil
}

// ===== WALLET OPERATIONS =====

// GetWalletBalance gets current wallet balance
func (s *Service) GetWalletBalance(ctx context.Context, userID int64) (float64, error) {
	if userID == 0 {
		return 0, fmt.Errorf("invalid user ID")
	}

	balance, err := s.queries.GetWalletBalance(ctx, userID)
	if err != nil {
		s.logger.Printf("failed to get wallet balance: %v", err)
		return 0, err
	}

	return balance, nil
}

// AddFundsToWallet adds funds to user's wallet
func (s *Service) AddFundsToWallet(ctx context.Context, input *AddFundsInput) (*WalletTransaction, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if input.Amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %.2f", input.Amount)
	}

	// Get or create wallet
	wallet, err := s.queries.GetUserWallet(ctx, input.UserID)
	if err != nil {
		wallet, err = s.queries.CreateUserWallet(ctx, input.UserID)
		if err != nil {
			return nil, err
		}
	}

	// Add funds
	tx, err := s.queries.AddFundsToWallet(ctx, wallet.ID, input.Amount)
	if err != nil {
		s.logger.Printf("failed to add funds to wallet: %v", err)
		return nil, err
	}

	s.logger.Printf("added %.2f to wallet for user %d", input.Amount, input.UserID)
	return tx, nil
}

// DeductFromWallet deducts funds from wallet
func (s *Service) DeductFromWallet(ctx context.Context, userID int64, amount float64) (*WalletTransaction, error) {
	if userID == 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %.2f", amount)
	}

	// Get wallet
	wallet, err := s.queries.GetUserWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("wallet not found: %v", err)
	}

	// Check balance
	if wallet.AvailableBalance < amount {
		return nil, fmt.Errorf("insufficient wallet balance")
	}

	// Deduct funds
	tx, err := s.queries.DeductFromWallet(ctx, wallet.ID, amount)
	if err != nil {
		s.logger.Printf("failed to deduct from wallet: %v", err)
		return nil, err
	}

	s.logger.Printf("deducted %.2f from wallet for user %d", amount, userID)
	return tx, nil
}

// ===== ESCROW OPERATIONS =====

// CreateEscrow creates an escrow account for an order
func (s *Service) CreateEscrow(ctx context.Context, orderID, buyerID, sellerID int64, amount float64) (*EscrowAccount, error) {
	if orderID == 0 || buyerID == 0 || sellerID == 0 {
		return nil, fmt.Errorf("order, buyer, and seller IDs required")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("invalid amount: %.2f", amount)
	}

	escrow, err := s.queries.CreateEscrowAccount(ctx, orderID, buyerID, sellerID, amount)
	if err != nil {
		s.logger.Printf("failed to create escrow: %v", err)
		return nil, err
	}

	s.logger.Printf("created escrow account: %d for order %d", escrow.ID, orderID)
	return escrow, nil
}

// ReleaseEscrow releases funds from escrow to seller
func (s *Service) ReleaseEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error) {
	if escrowID == 0 {
		return nil, fmt.Errorf("invalid escrow ID")
	}

	escrow, err := s.queries.ReleaseEscrow(ctx, escrowID)
	if err != nil {
		s.logger.Printf("failed to release escrow: %v", err)
		return nil, err
	}

	s.logger.Printf("released escrow: %d", escrowID)
	return escrow, nil
}

// RefundEscrow refunds escrow funds to buyer
func (s *Service) RefundEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error) {
	if escrowID == 0 {
		return nil, fmt.Errorf("invalid escrow ID")
	}

	escrow, err := s.queries.RefundEscrow(ctx, escrowID)
	if err != nil {
		s.logger.Printf("failed to refund escrow: %v", err)
		return nil, err
	}

	s.logger.Printf("refunded escrow: %d", escrowID)
	return escrow, nil
}

// ===== WEBHOOK OPERATIONS =====

// ProcessWebhook processes a webhook from a payment gateway
func (s *Service) ProcessWebhook(ctx context.Context, payload *WebhookPayload) (*WebhookEvent, error) {
	if payload.Provider == "" {
		return nil, fmt.Errorf("provider required")
	}

	// Log webhook
	logID, err := s.queries.CreateWebhookLog(ctx, payload.Provider, payload.EventType, nil)
	if err != nil {
		s.logger.Printf("failed to log webhook: %v", err)
		return nil, err
	}

	// Verify signature
	gateway := s.getGateway(payload.Provider)
	if gateway == nil {
		s.queries.UpdateWebhookProcessed(ctx, logID, false, "unsupported provider")
		return nil, fmt.Errorf("unsupported provider: %s", payload.Provider)
	}

	// Verify webhook signature (not implemented yet - TODO)
	// valid, err := gateway.VerifyWebhookSignature(payload.Signature, fmt.Sprintf("%d", payload.Timestamp), "")
	// if !valid {
	//     s.queries.UpdateWebhookProcessed(ctx, logID, false, "invalid signature")
	//     return nil, fmt.Errorf("invalid webhook signature")
	// }

	// Process webhook based on event type
	// TODO: Implement event-specific processing

	s.queries.UpdateWebhookProcessed(ctx, logID, true, "")

	return &WebhookEvent{
		ID:        logID,
		Provider:  payload.Provider,
		EventType: payload.EventType,
		Status:    "processed",
	}, nil
}

// ===== HELPER FUNCTIONS =====

// GeneratePaymentReference generates a unique payment reference
func GeneratePaymentReference() string {
	timestamp := time.Now().UnixNano()
	random := rand.Intn(100000)
	return fmt.Sprintf("PAY-%d-%d", timestamp, random)
}

// CalculateFees calculates payment fees based on amount and provider
func CalculateFees(amount float64, provider string) float64 {
	var feeRate float64

	switch provider {
	case "stripe":
		feeRate = 0.029 // 2.9% + $0.30 (plus fixed amount)
	case "razorpay":
		feeRate = 0.018 // 1.8%
	case "paystack":
		feeRate = 0.015 // 1.5% + ₦100 (plus fixed amount)
	case "flutterwave":
		feeRate = 0.015 // 1.5% + ₦50 (plus fixed amount)
	default:
		feeRate = 0.02 // 2% default
	}

	return amount * feeRate
}
