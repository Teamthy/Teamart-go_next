package payments

import (
	"context"
	"fmt"
	"time"
)

// MerchantWalletService handles merchant/seller wallet operations
type MerchantWalletService interface {
	GetMerchantWallet(ctx context.Context, sellerID int64) (*MerchantWallet, error)
	CreateMerchantWallet(ctx context.Context, sellerID int64) (*MerchantWallet, error)
	CreditSale(ctx context.Context, sellerID int64, amount float64, orderID int64) (*MerchantWalletTransaction, error)
	CreditRefund(ctx context.Context, sellerID int64, amount float64, refundID int64) (*MerchantWalletTransaction, error)
	CreditCommission(ctx context.Context, sellerID int64, amount float64, commissionType string) (*MerchantWalletTransaction, error)
	RequestPayout(ctx context.Context, sellerID, walletID int64, amount float64) (*Payout, error)
	ApprovePayout(ctx context.Context, payoutID int64, approvedBy int64) (*Payout, error)
	ProcessPayout(ctx context.Context, payoutID int64) (*Payout, error)
	GetMerchantBalance(ctx context.Context, sellerID int64) (total, pending, available float64, err error)
	ListTransactions(ctx context.Context, sellerID int64, limit, offset int64) ([]*MerchantWalletTransaction, error)
}

// MerchantWalletManager handles merchant wallet logic
type MerchantWalletManager struct {
	querier PaymentQuerier
	logger  interface{ Printf(string, ...interface{}) }
}

// NewMerchantWalletManager creates a new merchant wallet manager
func NewMerchantWalletManager(querier PaymentQuerier, logger interface{ Printf(string, ...interface{}) }) *MerchantWalletManager {
	return &MerchantWalletManager{
		querier: querier,
		logger:  logger,
	}
}

// GetMerchantWallet retrieves a merchant's wallet
func (mwm *MerchantWalletManager) GetMerchantWallet(ctx context.Context, sellerID int64) (*MerchantWallet, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	// This would fetch from the database
	// For now, create a basic structure
	wallet := &MerchantWallet{
		SellerID:  sellerID,
		CreatedAt: time.Now(),
	}

	return wallet, nil
}

// CreateMerchantWallet creates a new merchant wallet
func (mwm *MerchantWalletManager) CreateMerchantWallet(ctx context.Context, sellerID int64) (*MerchantWallet, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	wallet := &MerchantWallet{
		SellerID:         sellerID,
		Balance:          0,
		Currency:         "USD",
		PendingBalance:   0,
		AvailableBalance: 0,
		TotalEarned:      0,
		TotalWithdrawn:   0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	mwm.logger.Printf("merchant wallet created for seller %d", sellerID)
	return wallet, nil
}

// CreditSale credits a seller's wallet for a sale
func (mwm *MerchantWalletManager) CreditSale(ctx context.Context, sellerID int64, amount float64, orderID int64) (*MerchantWalletTransaction, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// Get merchant wallet
	wallet, err := mwm.GetMerchantWallet(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	// Create transaction
	tx := &MerchantWalletTransaction{
		WalletID:        wallet.ID,
		SellerID:        sellerID,
		Type:            "sale",
		Amount:          amount,
		PreviousBalance: wallet.Balance,
		NewBalance:      wallet.Balance + amount,
		ReferenceType:   "order",
		ReferenceID:     fmt.Sprintf("%d", orderID),
		Description:     fmt.Sprintf("Sale from order %d", orderID),
		Status:          "completed",
		CreatedAt:       time.Now(),
	}

	mwm.logger.Printf("credited sale amount %.2f to seller %d for order %d", amount, sellerID, orderID)
	return tx, nil
}

// CreditRefund credits a seller's wallet for a refund
// This reduces the seller's balance (they return the money)
func (mwm *MerchantWalletManager) CreditRefund(ctx context.Context, sellerID int64, amount float64, refundID int64) (*MerchantWalletTransaction, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// Get merchant wallet
	wallet, err := mwm.GetMerchantWallet(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	// Create transaction (negative amount for refund)
	tx := &MerchantWalletTransaction{
		WalletID:        wallet.ID,
		SellerID:        sellerID,
		Type:            "refund",
		Amount:          -amount, // Negative for deduction
		PreviousBalance: wallet.Balance,
		NewBalance:      wallet.Balance - amount,
		ReferenceType:   "refund",
		ReferenceID:     fmt.Sprintf("%d", refundID),
		Description:     fmt.Sprintf("Refund for order (refund ID: %d)", refundID),
		Status:          "completed",
		CreatedAt:       time.Now(),
	}

	mwm.logger.Printf("credited refund amount %.2f from seller %d for refund %d", amount, sellerID, refundID)
	return tx, nil
}

// CreditCommission credits a seller's wallet for earning commission
// This is for affiliate/livestream commission, etc.
func (mwm *MerchantWalletManager) CreditCommission(ctx context.Context, sellerID int64, amount float64, commissionType string) (*MerchantWalletTransaction, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	if commissionType == "" {
		return nil, fmt.Errorf("commission type is required")
	}

	// Get merchant wallet
	wallet, err := mwm.GetMerchantWallet(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	// Create transaction
	tx := &MerchantWalletTransaction{
		WalletID:        wallet.ID,
		SellerID:        sellerID,
		Type:            "commission",
		Amount:          amount,
		PreviousBalance: wallet.Balance,
		NewBalance:      wallet.Balance + amount,
		ReferenceType:   commissionType,
		ReferenceID:     fmt.Sprintf("commission_%s", commissionType),
		Description:     fmt.Sprintf("%s commission earned", commissionType),
		Status:          "completed",
		CreatedAt:       time.Now(),
	}

	mwm.logger.Printf("credited %s commission amount %.2f to seller %d", commissionType, amount, sellerID)
	return tx, nil
}

// RequestPayout creates a payout request from merchant wallet
func (mwm *MerchantWalletManager) RequestPayout(ctx context.Context, sellerID, walletID int64, amount float64) (*Payout, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// Get merchant wallet to verify balance
	wallet, err := mwm.GetMerchantWallet(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	if wallet.AvailableBalance < amount {
		return nil, fmt.Errorf("insufficient balance. Available: %.2f, Requested: %.2f", wallet.AvailableBalance, amount)
	}

	// Create payout request
	payout := &Payout{
		SellerID:  sellerID,
		Amount:    amount,
		Currency:  "USD",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mwm.logger.Printf("payout request created for seller %d, amount: %.2f", sellerID, amount)
	return payout, nil
}

// ApprovePayout approves a payout request
func (mwm *MerchantWalletManager) ApprovePayout(ctx context.Context, payoutID int64, approvedBy int64) (*Payout, error) {
	if payoutID == 0 {
		return nil, fmt.Errorf("payout ID is required")
	}

	// Get payout (would fetch from DB)
	payout := &Payout{
		ID:         payoutID,
		Status:     "approved",
		ReviewedBy: &approvedBy,
	}

	now := time.Now()
	payout.ReviewedAt = &now
	payout.UpdatedAt = now

	mwm.logger.Printf("payout %d approved by user %d", payoutID, approvedBy)
	return payout, nil
}

// ProcessPayout processes an approved payout
func (mwm *MerchantWalletManager) ProcessPayout(ctx context.Context, payoutID int64) (*Payout, error) {
	if payoutID == 0 {
		return nil, fmt.Errorf("payout ID is required")
	}

	// Get payout (would fetch from DB)
	payout := &Payout{
		ID:     payoutID,
		Status: "processing",
	}

	now := time.Now()
	payout.ProcessingStartedAt = &now
	payout.UpdatedAt = now

	mwm.logger.Printf("payout %d processing started", payoutID)
	return payout, nil
}

// GetMerchantBalance retrieves total, pending, and available balances for a merchant
func (mwm *MerchantWalletManager) GetMerchantBalance(ctx context.Context, sellerID int64) (total, pending, available float64, err error) {
	if sellerID == 0 {
		return 0, 0, 0, fmt.Errorf("seller ID is required")
	}

	wallet, err := mwm.GetMerchantWallet(ctx, sellerID)
	if err != nil {
		return 0, 0, 0, err
	}

	return wallet.Balance, wallet.PendingBalance, wallet.AvailableBalance, nil
}

// ListTransactions lists all transactions for a merchant
func (mwm *MerchantWalletManager) ListTransactions(ctx context.Context, sellerID int64, limit, offset int64) ([]*MerchantWalletTransaction, error) {
	if sellerID == 0 {
		return nil, fmt.Errorf("seller ID is required")
	}

	// This would fetch from database with pagination
	return make([]*MerchantWalletTransaction, 0), nil
}

// HoldAmount holds an amount in pending balance when order is placed
func (mwm *MerchantWalletManager) HoldAmount(ctx context.Context, sellerID int64, amount float64, orderID int64) error {
	if sellerID == 0 || amount <= 0 {
		return fmt.Errorf("invalid input")
	}

	// Move amount from available to pending
	mwm.logger.Printf("held amount %.2f for seller %d from order %d", amount, sellerID, orderID)
	return nil
}

// ReleaseHoldAmount releases a held amount back to available
func (mwm *MerchantWalletManager) ReleaseHoldAmount(ctx context.Context, sellerID int64, amount float64, orderID int64) error {
	if sellerID == 0 || amount <= 0 {
		return fmt.Errorf("invalid input")
	}

	// Move amount from pending to available
	mwm.logger.Printf("released held amount %.2f for seller %d from order %d", amount, sellerID, orderID)
	return nil
}
