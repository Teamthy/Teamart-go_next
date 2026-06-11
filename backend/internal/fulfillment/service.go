package fulfillment

import (
	"context"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Service provides business logic for fulfillment and inventory
type Service struct {
	queries *queries.Queries
	logger  *logger.Logger
}

// NewService creates a new fulfillment service
func NewService(q *queries.Queries, log *logger.Logger) *Service {
	return &Service{queries: q, logger: log}
}

// FulfillmentBatch represents a batch of orders to fulfill
type FulfillmentBatch struct {
	ID          int64
	BatchNumber string
	MerchantID  int64
	OrderIDs    []int64
	Status      string
	CreatedAt   time.Time
}

// CreateBatchInput input
type CreateBatchInput struct {
	MerchantID int64
	OrderIDs   []int64
	ShipByDate string
}

// CreateBatch creates a fulfillment batch (stub)
func (s *Service) CreateBatch(ctx context.Context, in *CreateBatchInput) (*FulfillmentBatch, error) {
	if in.MerchantID == 0 {
		return nil, fmt.Errorf("merchant_id required")
	}
	if len(in.OrderIDs) == 0 {
		return nil, fmt.Errorf("order_ids required")
	}

	batch := &FulfillmentBatch{
		ID:          1,
		BatchNumber: "BATCH-2026-001",
		MerchantID:  in.MerchantID,
		OrderIDs:    in.OrderIDs,
		Status:      "created",
		CreatedAt:   time.Now(),
	}

	s.logger.Infof("created fulfillment batch %s for merchant %d", batch.BatchNumber, in.MerchantID)
	return batch, nil
}

// ReserveStock reserves stock for an order (stub)
func (s *Service) ReserveStock(ctx context.Context, orderID int64, items []struct{ProductID, Quantity int64}) error {
	if orderID == 0 {
		return fmt.Errorf("order_id required")
	}
	// TODO: implement reservation using queries
	return nil
}

// UpdateInventory adjusts product inventory (stub)
func (s *Service) UpdateInventory(ctx context.Context, productID int64, delta int64, reason string) error {
	if productID == 0 {
		return fmt.Errorf("product_id required")
	}
	// TODO: persist inventory change
	return nil
}
