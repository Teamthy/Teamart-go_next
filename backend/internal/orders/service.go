package orders

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
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
	ID            int64
	OrderNumber   string
	UserID        int64
	CustomerName  string
	CustomerEmail string
	TotalAmount   float64
	Status        string
	PaymentMethod string
	ItemsCount    int64
	CreatedAt     string
	UpdatedAt     string
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

	totalAmount := strconv.FormatFloat(input.TotalAmount, 'f', -1, 64)
	order, err := s.queries.CreateOrder(ctx, int32(input.UserID), totalAmount, input.Status)
	if err != nil {
		s.logger.Errorf("failed to create order: %v", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	s.logger.Infof("order created successfully with ID: %d", order.ID)

	orderData, err := s.buildOrderData(ctx, order)
	if err != nil {
		return nil, err
	}

	return &CreateOrderOutput{
		ID:            orderData.ID,
		OrderNumber:   orderData.OrderNumber,
		UserID:        orderData.UserID,
		CustomerName:  orderData.CustomerName,
		CustomerEmail: orderData.CustomerEmail,
		TotalAmount:   orderData.TotalAmount,
		Status:        orderData.Status,
		PaymentMethod: orderData.PaymentMethod,
		ItemsCount:    orderData.ItemsCount,
		CreatedAt:     orderData.CreatedAt,
		UpdatedAt:     orderData.UpdatedAt,
	}, nil
}

// GetOrderByIDInput represents the input
type GetOrderByIDInput struct {
	OrderID int64
}

// GetOrderByIDOutput represents the output
type GetOrderByIDOutput struct {
	ID            int64
	OrderNumber   string
	UserID        int64
	CustomerName  string
	CustomerEmail string
	TotalAmount   float64
	Status        string
	PaymentMethod string
	ItemsCount    int64
	CreatedAt     string
	UpdatedAt     string
}

// GetOrderByID retrieves an order by its ID
func (s *Service) GetOrderByID(ctx context.Context, input *GetOrderByIDInput) (*GetOrderByIDOutput, error) {
	if input.OrderID == 0 {
		return nil, fmt.Errorf("order ID is required")
	}

	s.logger.Debugf("fetching order with ID: %d", input.OrderID)

	order, err := s.queries.GetOrderByID(ctx, int32(input.OrderID))
	if err != nil {
		s.logger.Errorf("failed to fetch order: %v", err)
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	orderData, err := s.buildOrderData(ctx, order)
	if err != nil {
		return nil, err
	}

	return &GetOrderByIDOutput{
		ID:            orderData.ID,
		OrderNumber:   orderData.OrderNumber,
		UserID:        orderData.UserID,
		CustomerName:  orderData.CustomerName,
		CustomerEmail: orderData.CustomerEmail,
		TotalAmount:   orderData.TotalAmount,
		Status:        orderData.Status,
		PaymentMethod: orderData.PaymentMethod,
		ItemsCount:    orderData.ItemsCount,
		CreatedAt:     orderData.CreatedAt,
		UpdatedAt:     orderData.UpdatedAt,
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
	ID            int64
	OrderNumber   string
	UserID        int64
	CustomerName  string
	CustomerEmail string
	TotalAmount   float64
	Status        string
	PaymentMethod string
	ItemsCount    int64
	CreatedAt     string
	UpdatedAt     string
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

	orders, err := s.queries.ListOrdersByUserID(ctx, int32(input.UserID), input.Limit, input.Offset)
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
		orderData, err := s.buildOrderData(ctx, order)
		if err != nil {
			return nil, err
		}
		output.Orders[i] = orderData
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

	orders, err := s.queries.ListOrdersByStatus(ctx, input.Status, input.Limit, input.Offset)
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
		orderData, err := s.buildOrderData(ctx, order)
		if err != nil {
			return nil, err
		}
		output.Orders[i] = orderData
	}

	s.logger.Infof("fetched %d orders with status: %s", len(orders), input.Status)

	return output, nil
}

// ListAllOrdersInput represents the input
type ListAllOrdersInput struct {
	Limit  int32
	Offset int32
}

// UpdateOrderStatusInput represents the input for updating an order status
type UpdateOrderStatusInput struct {
	OrderID int64
	Status  string
}

// UpdateOrderStatusOutput represents the output after updating order status
type UpdateOrderStatusOutput struct {
	ID            int64
	OrderNumber   string
	UserID        int64
	CustomerName  string
	CustomerEmail string
	TotalAmount   float64
	Status        string
	PaymentMethod string
	ItemsCount    int64
	CreatedAt     string
	UpdatedAt     string
}

// UpdateOrderStatus updates the status of an order
func (s *Service) UpdateOrderStatus(ctx context.Context, input *UpdateOrderStatusInput) (*UpdateOrderStatusOutput, error) {
	if input.OrderID == 0 {
		return nil, fmt.Errorf("order ID is required")
	}
	if input.Status == "" {
		return nil, fmt.Errorf("status is required")
	}

	order, err := s.queries.UpdateOrderStatus(ctx, int32(input.OrderID), input.Status)
	if err != nil {
		s.logger.Errorf("failed to update order status: %v", err)
		return nil, fmt.Errorf("failed to update order status: %w", err)
	}

	orderData, err := s.buildOrderData(ctx, order)
	if err != nil {
		return nil, err
	}

	return &UpdateOrderStatusOutput{
		ID:            orderData.ID,
		OrderNumber:   orderData.OrderNumber,
		UserID:        orderData.UserID,
		CustomerName:  orderData.CustomerName,
		CustomerEmail: orderData.CustomerEmail,
		TotalAmount:   orderData.TotalAmount,
		Status:        orderData.Status,
		PaymentMethod: orderData.PaymentMethod,
		ItemsCount:    orderData.ItemsCount,
		CreatedAt:     orderData.CreatedAt,
		UpdatedAt:     orderData.UpdatedAt,
	}, nil
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

	orders, err := s.queries.ListAllOrders(ctx, input.Limit, input.Offset)
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
		orderData, err := s.buildOrderData(ctx, order)
		if err != nil {
			return nil, err
		}
		output.Orders[i] = orderData
	}

	s.logger.Infof("fetched %d orders", len(orders))

	return output, nil
}

func (s *Service) buildOrderData(ctx context.Context, order queries.Order) (OrderData, error) {
	totalAmountValue, err := numericToFloat64(order.TotalAmount)
	if err != nil {
		s.logger.Errorf("failed to parse order total amount: %v", err)
		return OrderData{}, fmt.Errorf("failed to parse order total amount: %w", err)
	}

	orderNumber := buildOrderNumber(int64(order.ID))
	customerName := "Guest"
	customerEmail := ""
	if order.UserID != 0 {
		user, err := s.queries.GetUserByID(ctx, order.UserID)
		if err == nil {
			if user.Name != "" {
				customerName = user.Name
			} else if user.Email != "" {
				customerName = user.Email
			}
			customerEmail = user.Email
		} else {
			s.logger.Debugf("could not load customer info for user %d: %v", order.UserID, err)
		}
	}

	itemsCount, err := s.queries.CountOrderItems(ctx, order.ID)
	if err != nil {
		s.logger.Debugf("could not count order items for order %d: %v", order.ID, err)
		itemsCount = 0
	}

	paymentMethod := "card"
	if order.Status == "pending" {
		paymentMethod = "card"
	}

	return OrderData{
		ID:            int64(order.ID),
		OrderNumber:   orderNumber,
		UserID:        int64(order.UserID),
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		TotalAmount:   totalAmountValue,
		Status:        order.Status,
		PaymentMethod: paymentMethod,
		ItemsCount:    itemsCount,
		CreatedAt:     order.CreatedAt.String(),
		UpdatedAt:     order.UpdatedAt.String(),
	}, nil
}

func buildOrderNumber(orderID int64) string {
	return fmt.Sprintf("ORD-%06d", orderID)
}

func numericToFloat64(value pgtype.Numeric) (float64, error) {
	floatValue, err := value.Float64Value()
	if err != nil {
		return 0, err
	}
	return floatValue.Float64, nil
}
