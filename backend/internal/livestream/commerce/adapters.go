package commerce

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/teamart/commerce-api/internal/orders"
	"github.com/teamart/commerce-api/internal/products"
	"github.com/teamart/commerce-api/pkg/logger"
)

type ProductServiceAdapter struct {
	service *products.Service
}

func NewProductServiceAdapter(service *products.Service) ProductFetcher {
	return &ProductServiceAdapter{service: service}
}

func (a *ProductServiceAdapter) FetchProduct(productID int64) (*ProductDetails, error) {
	if productID == 0 {
		return nil, fmt.Errorf("product ID is required")
	}

	product, err := a.service.GetProductByID(context.Background(), &products.GetProductByIDInput{ProductID: productID})
	if err != nil {
		return nil, err
	}

	return &ProductDetails{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

type OrderServiceAdapter struct {
	service *orders.Service
}

func NewOrderServiceAdapter(service *orders.Service) OrderCreator {
	return &OrderServiceAdapter{service: service}
}

func (a *OrderServiceAdapter) CreateOrder(input *OrderInput) (*OrderResult, error) {
	if input == nil {
		return nil, fmt.Errorf("order input is required")
	}

	order, err := a.service.CreateOrder(context.Background(), &orders.CreateOrderInput{
		UserID:      input.UserID,
		TotalAmount: input.TotalAmount,
		Status:      input.Status,
	})
	if err != nil {
		return nil, err
	}

	return &OrderResult{
		OrderID:   order.ID,
		Reference: fmt.Sprintf("order-%d", order.ID),
	}, nil
}

type SimplePaymentProcessor struct {
	logger *logger.Logger
}

func NewSimplePaymentProcessor(logger *logger.Logger) PaymentProcessor {
	return &SimplePaymentProcessor{logger: logger}
}

func (p *SimplePaymentProcessor) Charge(userID int64, amount float64, method string) (*PaymentResult, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("payment amount must be greater than zero")
	}
	if method == "" {
		method = "card"
	}

	paymentID := fmt.Sprintf("payment-%d", time.Now().UnixNano())
	if p.logger != nil {
		p.logger.Infof("processing payment user=%d amount=%.2f method=%s payment_id=%s", userID, amount, method, paymentID)
	}

	return &PaymentResult{PaymentID: paymentID, Method: method}, nil
}

type WalletEntry struct {
	Amount     float64
	Reason     string
	RecordedAt time.Time
}

type InMemoryWalletManager struct {
	mu     sync.Mutex
	ledger map[int64][]WalletEntry
	logger *logger.Logger
}

func NewInMemoryWalletManager(logger *logger.Logger) WalletManager {
	return &InMemoryWalletManager{
		ledger: make(map[int64][]WalletEntry),
		logger: logger,
	}
}

func (m *InMemoryWalletManager) RecordPayout(userID int64, amount float64, reason string) error {
	if amount == 0 {
		return nil
	}

	entry := WalletEntry{
		Amount:     amount,
		Reason:     reason,
		RecordedAt: time.Now(),
	}

	m.mu.Lock()
	m.ledger[userID] = append(m.ledger[userID], entry)
	m.mu.Unlock()

	if m.logger != nil {
		if userID == 0 {
			m.logger.Infof("platform payout recorded amount=%.2f reason=%s", amount, reason)
		} else {
			m.logger.Infof("wallet payout recorded user=%d amount=%.2f reason=%s", userID, amount, reason)
		}
	}

	return nil
}

type LoggerEventDispatcher struct {
	logger *logger.Logger
}

func NewLoggerEventDispatcher(logger *logger.Logger) EventDispatcher {
	return &LoggerEventDispatcher{logger: logger}
}

func (d *LoggerEventDispatcher) PublishEvent(eventType string, payload map[string]interface{}) error {
	if d.logger != nil {
		d.logger.Infof("publish event %s payload=%+v", eventType, payload)
	}
	return nil
}
