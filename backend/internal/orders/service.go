package orders

import (
	"context"
	"fmt"

	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Service provides business logic for order operations
type Service struct {
	queries *queries.Queries
	logger  *logger.Logger
}

// NewService creates a new order service
func NewService(queries *queries.Queries, logger *logger.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger,
	}
}

// CreateOrderInput represents the input for creating an order
type CreateOrderInput struct {
	UserID      int64
	TotalAmount float64
	Status      string
}

// CreateOrderOutput represents the output after creating an order
type CreateOrderOutput struct {
	ID          int64
	UserID      int64
	TotalAmount float64
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

// CreateOrder creates a new order with validation
func (s *Service) CreateOrder(ctx context.Context, input *CreateOrderInput) (*CreateOrderOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.TotalAmount <= 0 {
		return nil, fmt.Errorf("total amount must be greater than zero")
	}
	if input.Status == "" {
		input.Status = "pending"
	}

	s.logger.Debugf("creating order for user: %d with amount: %.2f", input.UserID, input.TotalAmount)

	order, err := s.queries.CreateOrder(ctx, queries.CreateOrderParams{
		UserID:      input.UserID,
		TotalAmount: input.TotalAmount,
		Status:      input.Status,
	})
	if err != nil {
		s.logger.Errorf("failed to create order: %v", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	s.logger.Infof("order created successfully with ID: %d", order.ID)

	return &CreateOrderOutput{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.String(),
		UpdatedAt:   order.UpdatedAt.String(),
	}, nil
}

// GetOrderByIDInput represents the input
type GetOrderByIDInput struct {
	OrderID int64
}

// GetOrderByIDOutput represents the output
type GetOrderByIDOutput struct {
	ID          int64
	UserID      int64
	TotalAmount float64
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

// GetOrderByID retrieves an order by its ID
func (s *Service) GetOrderByID(ctx context.Context, input *GetOrderByIDInput) (*GetOrderByIDOutput, error) {
	if input.OrderID == 0 {
		return nil, fmt.Errorf("order ID is required")
	}

	s.logger.Debugf("fetching order with ID: %d", input.OrderID)

	order, err := s.queries.GetOrderByID(ctx, input.OrderID)
	if err != nil {
		s.logger.Errorf("failed to fetch order: %v", err)
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	return &GetOrderByIDOutput{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.String(),
		UpdatedAt:   order.UpdatedAt.String(),
	}, nil
}

// ListOrdersByUserIDInput represents the input
type ListOrdersByUserIDInput struct {
	UserID int64
	Limit  int32
	Offset int32
}

// ListOrdersOutput represents the output
type ListOrdersOutput struct {
	Orders []OrderData
	Limit  int32
	Offset int32
}

type OrderData struct {
	ID          int64
	UserID      int64
	TotalAmount float64
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

// ListOrdersByUserID retrieves orders by user ID
func (s *Service) ListOrdersByUserID(ctx context.Context, input *ListOrdersByUserIDInput) (*ListOrdersOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	s.logger.Debugf("listing orders for user: %d with limit: %d, offset: %d", input.UserID, input.Limit, input.Offset)

	orders, err := s.queries.ListOrdersByUserID(ctx, queries.ListOrdersByUserIDParams{
		UserID: input.UserID,
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		s.logger.Errorf("failed to list orders: %v", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	output := &ListOrdersOutput{
		Orders: make([]OrderData, len(orders)),
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	for i, order := range orders {
		output.Orders[i] = OrderData{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt.String(),
			UpdatedAt:   order.UpdatedAt.String(),
		}
	}

	s.logger.Infof("fetched %d orders for user: %d", len(orders), input.UserID)

	return output, nil
}

// ListOrdersByStatusInput represents the input
type ListOrdersByStatusInput struct {
	Status string
	Limit  int32
	Offset int32
}

// ListOrdersByStatus retrieves orders by status
func (s *Service) ListOrdersByStatus(ctx context.Context, input *ListOrdersByStatusInput) (*ListOrdersOutput, error) {
	if input.Status == "" {
		return nil, fmt.Errorf("status is required")
	}
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	s.logger.Debugf("listing orders with status: %s", input.Status)

	orders, err := s.queries.ListOrdersByStatus(ctx, queries.ListOrdersByStatusParams{
		Status: input.Status,
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		s.logger.Errorf("failed to list orders: %v", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	output := &ListOrdersOutput{
		Orders: make([]OrderData, len(orders)),
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	for i, order := range orders {
		output.Orders[i] = OrderData{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt.String(),
			UpdatedAt:   order.UpdatedAt.String(),
		}
	}

	s.logger.Infof("fetched %d orders with status: %s", len(orders), input.Status)

	return output, nil
}

// ListAllOrdersInput represents the input
type ListAllOrdersInput struct {
	Limit  int32
	Offset int32
}

// ListAllOrders retrieves all orders with pagination
func (s *Service) ListAllOrders(ctx context.Context, input *ListAllOrdersInput) (*ListOrdersOutput, error) {
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	s.logger.Debugf("listing all orders with limit: %d, offset: %d", input.Limit, input.Offset)

	orders, err := s.queries.ListAllOrders(ctx, queries.ListAllOrdersParams{
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		s.logger.Errorf("failed to list orders: %v", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	output := &ListOrdersOutput{
		Orders: make([]OrderData, len(orders)),
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	for i, order := range orders {
		output.Orders[i] = OrderData{
			ID:          order.ID,
			UserID:      order.UserID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   order.CreatedAt.String(),
			UpdatedAt:   order.UpdatedAt.String(),
		}
	}

	s.logger.Infof("fetched %d orders", len(orders))

	return output, nil
}
