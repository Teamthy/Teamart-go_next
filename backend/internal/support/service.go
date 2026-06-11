package support

import (
	"context"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Service provides support ticket handling
type Service struct {
	queries *queries.Queries
	logger  *logger.Logger
}

// NewService creates a new support service
func NewService(q *queries.Queries, log *logger.Logger) *Service {
	return &Service{queries: q, logger: log}
}

// SupportTicket represents a ticket
type SupportTicket struct {
	ID          int64
	TicketNumber string
	MerchantID  int64
	CustomerID  int64
	Subject     string
	Description string
	Status      string
	Priority    string
	CreatedAt   time.Time
}

// CreateTicket creates a support ticket (stub)
func (s *Service) CreateTicket(ctx context.Context, merchantID, customerID int64, subject, description, category, priority string) (*SupportTicket, error) {
	if subject == "" || description == "" {
		return nil, fmt.Errorf("subject and description required")
	}
	t := &SupportTicket{
		ID:           1,
		TicketNumber: "TKT-2026-001",
		MerchantID:   merchantID,
		CustomerID:   customerID,
		Subject:      subject,
		Description:  description,
		Status:       "open",
		Priority:     priority,
		CreatedAt:    time.Now(),
	}
	return t, nil
}

// AddMessage adds a message to a ticket (stub)
func (s *Service) AddMessage(ctx context.Context, ticketID int64, message string, isInternal bool) error {
	if ticketID == 0 {
		return fmt.Errorf("ticket_id required")
	}
	if message == "" {
		return fmt.Errorf("message required")
	}
	// TODO: persist message
	return nil
}
