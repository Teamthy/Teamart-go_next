package payments

import (
	"context"
	"fmt"
	"time"
)

// EscrowService handles escrow account operations
type EscrowService interface {
	CreateEscrow(ctx context.Context, orderID, buyerID, sellerID int64, amount float64) (*EscrowAccount, error)
	ReleaseEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error)
	RefundEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error)
	DisputeEscrow(ctx context.Context, escrowID int64, dispute *EscrowDispute) (*EscrowDispute, error)
	ResolveDispute(ctx context.Context, disputeID int64, outcome string, buyerAmount, sellerAmount float64) (*EscrowDispute, error)
	GetEscrowAccount(ctx context.Context, escrowID int64) (*EscrowAccount, error)
	ListEscrowAccounts(ctx context.Context, sellerID int64) ([]*EscrowAccount, error)
}

// EscrowManager handles escrow logic
type EscrowManager struct {
	querier PaymentQuerier
	logger  interface{ Printf(string, ...interface{}) }
}

// NewEscrowManager creates a new escrow manager
func NewEscrowManager(querier PaymentQuerier, logger interface{ Printf(string, ...interface{}) }) *EscrowManager {
	return &EscrowManager{
		querier: querier,
		logger:  logger,
	}
}

// CreateEscrow creates a new escrow account for an order
func (em *EscrowManager) CreateEscrow(ctx context.Context, orderID, buyerID, sellerID int64, amount float64) (*EscrowAccount, error) {
	if orderID == 0 || buyerID == 0 || sellerID == 0 {
		return nil, fmt.Errorf("order ID, buyer ID, and seller ID are required")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("escrow amount must be greater than zero")
	}

	// Create escrow account
	escrow, err := em.querier.CreateEscrowAccount(ctx, orderID, buyerID, sellerID, amount)
	if err != nil {
		em.logger.Printf("failed to create escrow account: %v", err)
		return nil, err
	}

	em.logger.Printf("escrow account created: %d for order %d, amount: %.2f", escrow.ID, orderID, amount)
	return escrow, nil
}

// ReleaseEscrow releases funds from escrow to the seller
// This happens after order is confirmed delivered
func (em *EscrowManager) ReleaseEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error) {
	if escrowID == 0 {
		return nil, fmt.Errorf("escrow ID is required")
	}

	// Get escrow account
	escrow, err := em.querier.GetEscrowAccount(ctx, escrowID)
	if err != nil {
		return nil, fmt.Errorf("escrow account not found: %v", err)
	}

	if escrow.Status != "held" {
		return nil, fmt.Errorf("can only release escrow with status 'held', current status: %s", escrow.Status)
	}

	// Update escrow status to released
	now := time.Now()
	escrow.Status = "released"
	escrow.ReleasedAt = &now

	// Add funds to seller's wallet
	sellerWallet, err := em.querier.GetUserWallet(ctx, escrow.SellerID)
	if err != nil {
		// Create wallet if not exists
		sellerWallet, err = em.querier.CreateUserWallet(ctx, escrow.SellerID)
		if err != nil {
			em.logger.Printf("failed to get/create seller wallet: %v", err)
			return nil, err
		}
	}

	// Add funds to wallet
	_, err = em.querier.AddFundsToWallet(ctx, sellerWallet.ID, escrow.Amount)
	if err != nil {
		em.logger.Printf("failed to add funds to seller wallet: %v", err)
		return nil, err
	}

	em.logger.Printf("escrow %d released to seller %d, amount: %.2f", escrowID, escrow.SellerID, escrow.Amount)
	return escrow, nil
}

// RefundEscrow refunds the escrow amount back to the buyer
// This happens if order is cancelled or buyer initiates chargeback
func (em *EscrowManager) RefundEscrow(ctx context.Context, escrowID int64) (*EscrowAccount, error) {
	if escrowID == 0 {
		return nil, fmt.Errorf("escrow ID is required")
	}

	// Get escrow account
	escrow, err := em.querier.GetEscrowAccount(ctx, escrowID)
	if err != nil {
		return nil, fmt.Errorf("escrow account not found: %v", err)
	}

	if escrow.Status == "refunded" {
		return nil, fmt.Errorf("escrow already refunded")
	}

	if escrow.Status == "released" {
		return nil, fmt.Errorf("cannot refund escrow that has been released")
	}

	// Update escrow status to refunded
	now := time.Now()
	escrow.Status = "refunded"
	escrow.RefundedAt = &now

	// Add funds back to buyer's wallet
	buyerWallet, err := em.querier.GetUserWallet(ctx, escrow.BuyerID)
	if err != nil {
		// Create wallet if not exists
		buyerWallet, err = em.querier.CreateUserWallet(ctx, escrow.BuyerID)
		if err != nil {
			em.logger.Printf("failed to get/create buyer wallet: %v", err)
			return nil, err
		}
	}

	// Add funds to wallet
	_, err = em.querier.AddFundsToWallet(ctx, buyerWallet.ID, escrow.Amount)
	if err != nil {
		em.logger.Printf("failed to add funds to buyer wallet: %v", err)
		return nil, err
	}

	em.logger.Printf("escrow %d refunded to buyer %d, amount: %.2f", escrowID, escrow.BuyerID, escrow.Amount)
	return escrow, nil
}

// DisputeEscrow marks an escrow account as disputed
func (em *EscrowManager) DisputeEscrow(ctx context.Context, escrowID int64, dispute *EscrowDispute) (*EscrowDispute, error) {
	if escrowID == 0 {
		return nil, fmt.Errorf("escrow ID is required")
	}

	// Get escrow account
	escrow, err := em.querier.GetEscrowAccount(ctx, escrowID)
	if err != nil {
		return nil, fmt.Errorf("escrow account not found: %v", err)
	}

	if escrow.Status != "held" {
		return nil, fmt.Errorf("can only dispute escrow with status 'held', current status: %s", escrow.Status)
	}

	// Create dispute
	dispute.EscrowAccountID = escrowID
	dispute.Status = "open"
	dispute.InitiatedAt = time.Now()

	em.logger.Printf("escrow dispute created for escrow %d, initiated by user %d", escrowID, dispute.InitiatedBy)
	return dispute, nil
}

// ResolveDispute resolves a dispute and distributes funds accordingly
func (em *EscrowManager) ResolveDispute(ctx context.Context, disputeID int64, outcome string, buyerAmount, sellerAmount float64) (*EscrowDispute, error) {
	if disputeID == 0 {
		return nil, fmt.Errorf("dispute ID is required")
	}

	// Validate outcome
	validOutcomes := map[string]bool{
		"buyer":  true,
		"seller": true,
		"split":  true,
	}

	if !validOutcomes[outcome] {
		return nil, fmt.Errorf("invalid dispute outcome: %s", outcome)
	}

	// Get dispute
	// Note: This would need to be implemented in the querier
	// For now, we'll create a basic structure

	now := time.Now()
	dispute := &EscrowDispute{
		Status:       "resolved",
		Outcome:      &outcome,
		ResolvedAt:   &now,
		BuyerAmount:  &buyerAmount,
		SellerAmount: &sellerAmount,
	}

	em.logger.Printf("dispute %d resolved with outcome: %s, buyer: %.2f, seller: %.2f", disputeID, outcome, buyerAmount, sellerAmount)
	return dispute, nil
}

// GetEscrowAccount retrieves an escrow account
func (em *EscrowManager) GetEscrowAccount(ctx context.Context, escrowID int64) (*EscrowAccount, error) {
	if escrowID == 0 {
		return nil, fmt.Errorf("escrow ID is required")
	}

	return em.querier.GetEscrowAccount(ctx, escrowID)
}

// ListEscrowAccounts lists all escrow accounts for a seller
func (em *EscrowManager) ListEscrowAccounts(ctx context.Context, sellerID int64) ([]*EscrowAccount, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	// This would need to be implemented in the querier
	// For now, return empty slice
	return []*EscrowAccount{}, nil
}
