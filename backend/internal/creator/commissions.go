package creator

import (
	"context"
	"fmt"
)

// CommissionService calculates creator commissions.
type CommissionService struct {
	Rate float64
}

// NewCommissionService creates a commission service.
func NewCommissionService(rate float64) *CommissionService {
	return &CommissionService{Rate: rate}
}

// CalculateCommission computes the commission for a sale.
func (c *CommissionService) CalculateCommission(ctx context.Context, saleAmount float64) float64 {
	revenue := saleAmount * c.Rate
	fmt.Printf("calculated commission %.2f on sale %.2f\n", revenue, saleAmount)
	return revenue
}

// DistributeCommission sends commission payments to creators.
func (c *CommissionService) DistributeCommission(ctx context.Context, creatorID int64, amount float64) error {
	fmt.Printf("distributing commission %.2f to creator %d\n", amount, creatorID)
	return nil
}
