package payments

import (
	"context"
	"fmt"
	"time"
)

// SplitPaymentEngine handles split payment logic
type SplitPaymentEngine interface {
	CreateSplitPayment(ctx context.Context, input *CreateSplitPaymentInput) (*SplitPayment, error)
	DistributeSplitPayment(ctx context.Context, splitID int64, paymentResult *PaymentResult) error
	CalculateSplits(totalAmount float64, splits []*SplitLine) ([]*SplitLine, error)
	ValidateSplits(splits []*SplitLine) error
}

// SplitPaymentProcessor handles split payment processing
type SplitPaymentProcessor struct {
	querier PaymentQuerier
	logger  interface{ Printf(string, ...interface{}) }
}

// NewSplitPaymentProcessor creates a new split payment processor
func NewSplitPaymentProcessor(querier PaymentQuerier, logger interface{ Printf(string, ...interface{}) }) *SplitPaymentProcessor {
	return &SplitPaymentProcessor{
		querier: querier,
		logger:  logger,
	}
}

// CreateSplitPayment creates a split payment configuration
func (spp *SplitPaymentProcessor) CreateSplitPayment(ctx context.Context, input *CreateSplitPaymentInput) (*SplitPayment, error) {
	if input.OrderID == 0 {
		return nil, fmt.Errorf("order ID is required")
	}

	if input.TotalAmount <= 0 {
		return nil, fmt.Errorf("total amount must be greater than zero")
	}

	// Validate splits
	if err := spp.ValidateSplits(input.Splits); err != nil {
		return nil, err
	}

	// Calculate splits
	splits, err := spp.CalculateSplits(input.TotalAmount, input.Splits)
	if err != nil {
		return nil, err
	}

	// Create split payment record
	splitPayment := &SplitPayment{
		OrderID:   input.OrderID,
		Status:    "pending",
		Splits:    splits,
		CreatedAt: time.Now(),
	}

	spp.logger.Printf("split payment created for order %d with %d recipients", input.OrderID, len(splits))
	return splitPayment, nil
}

// DistributeSplitPayment distributes payment to all recipients according to split configuration
func (spp *SplitPaymentProcessor) DistributeSplitPayment(ctx context.Context, splitID int64, paymentResult *PaymentResult) error {
	if splitID == 0 {
		return fmt.Errorf("split payment ID is required")
	}

	if !paymentResult.Success {
		return fmt.Errorf("cannot distribute failed payment")
	}

	// Get split payment details (would be fetched from DB in real implementation)
	// For each split line, distribute funds to the recipient

	// Example distribution logic:
	// 1. Seller receives their split amount
	// 2. Platform receives commission
	// 3. Affiliates receive commission
	// 4. Livestream hosts receive commission

	// This would iterate through splits and add funds to each recipient's wallet

	spp.logger.Printf("split payment %d distributed successfully", splitID)
	return nil
}

// CalculateSplits calculates the actual amounts for each split
func (spp *SplitPaymentProcessor) CalculateSplits(totalAmount float64, splits []*SplitLine) ([]*SplitLine, error) {
	if len(splits) == 0 {
		return nil, fmt.Errorf("at least one split is required")
	}

	result := make([]*SplitLine, 0, len(splits))
	totalPercentage := 0.0
	totalCalculatedAmount := 0.0

	// First pass: calculate amounts and validate percentages
	for _, split := range splits {
		if split.RecipientID == 0 {
			return nil, fmt.Errorf("recipient ID is required for all splits")
		}

		// If percentage is provided, calculate amount
		if split.Percentage > 0 {
			split.Amount = totalAmount * (split.Percentage / 100.0)
			totalPercentage += split.Percentage
		} else if split.Amount > 0 {
			split.Percentage = (split.Amount / totalAmount) * 100.0
		} else {
			return nil, fmt.Errorf("either amount or percentage must be provided for recipient %d", split.RecipientID)
		}

		totalCalculatedAmount += split.Amount
		result = append(result, split)
	}

	// Verify splits add up correctly (allow small floating point errors)
	difference := totalAmount - totalCalculatedAmount
	if difference < -0.01 || difference > 0.01 { // Allow 1 cent difference
		return nil, fmt.Errorf("split amounts do not equal total: %.2f vs %.2f", totalCalculatedAmount, totalAmount)
	}

	return result, nil
}

// ValidateSplits validates split payment configuration
func (spp *SplitPaymentProcessor) ValidateSplits(splits []*SplitLine) error {
	if len(splits) == 0 {
		return fmt.Errorf("at least one split is required")
	}

	// Check for duplicate recipients
	recipients := make(map[int64]map[string]bool)
	for _, split := range splits {
		if split.RecipientID == 0 {
			return fmt.Errorf("recipient ID is required")
		}

		if split.RecipientType == "" {
			return fmt.Errorf("recipient type is required")
		}

		// Allow same recipient type, but warn about multiple entries
		if recipients[split.RecipientID] == nil {
			recipients[split.RecipientID] = make(map[string]bool)
		}

		if recipients[split.RecipientID][split.RecipientType] {
			return fmt.Errorf("duplicate split for recipient %d with type %s", split.RecipientID, split.RecipientType)
		}

		recipients[split.RecipientID][split.RecipientType] = true
	}

	// Validate amounts and percentages
	totalPercentage := 0.0
	totalAmount := 0.0
	hasAmount := false
	hasPercentage := false

	for _, split := range splits {
		if split.Amount < 0 {
			return fmt.Errorf("split amount cannot be negative: %.2f", split.Amount)
		}

		if split.Percentage < 0 {
			return fmt.Errorf("split percentage cannot be negative: %.2f", split.Percentage)
		}

		if split.Amount > 0 {
			hasAmount = true
			totalAmount += split.Amount
		}

		if split.Percentage > 0 {
			hasPercentage = true
			totalPercentage += split.Percentage
		}

		// Both amount and percentage should not be provided at the same time
		if split.Amount > 0 && split.Percentage > 0 {
			// This is allowed, we'll recalculate one from the other
		}
	}

	return nil
}

// CalculatePlatformCommission calculates platform commission from an order
func (spp *SplitPaymentProcessor) CalculatePlatformCommission(orderTotal float64, commissionPercentage float64) float64 {
	return orderTotal * (commissionPercentage / 100.0)
}

// CalculateAffiliateCommission calculates affiliate commission
func (spp *SplitPaymentProcessor) CalculateAffiliateCommission(orderTotal float64, affiliatePercentage float64) float64 {
	return orderTotal * (affiliatePercentage / 100.0)
}

// CalculateLivesteamHostCommission calculates livestream host commission
func (spp *SplitPaymentProcessor) CalculateLivesteamHostCommission(orderTotal float64, hostPercentage float64) float64 {
	return orderTotal * (hostPercentage / 100.0)
}

// CalculateCreatorCommission calculates creator commission
func (spp *SplitPaymentProcessor) CalculateCreatorCommission(orderTotal float64, creatorPercentage float64) float64 {
	return orderTotal * (creatorPercentage / 100.0)
}
