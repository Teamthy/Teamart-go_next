package creator

import (
	"context"
	"fmt"
)

// CreatorWallet manages creator earning balances and payouts.
type CreatorWallet struct {
	Balances map[int64]float64
}

// NewCreatorWallet creates a new creator wallet service.
func NewCreatorWallet() *CreatorWallet {
	return &CreatorWallet{Balances: make(map[int64]float64)}
}

// CreditBalance credits a creator wallet.
func (w *CreatorWallet) CreditBalance(ctx context.Context, creatorID int64, amount float64) error {
	w.Balances[creatorID] += amount
	fmt.Printf("credited %.2f to creator wallet %d\n", amount, creatorID)
	return nil
}

// WithdrawBalance processes a creator withdrawal request.
func (w *CreatorWallet) WithdrawBalance(ctx context.Context, creatorID int64, amount float64) error {
	if w.Balances[creatorID] < amount {
		return fmt.Errorf("insufficient balance")
	}
	w.Balances[creatorID] -= amount
	fmt.Printf("withdrew %.2f from creator wallet %d\n", amount, creatorID)
	return nil
}
