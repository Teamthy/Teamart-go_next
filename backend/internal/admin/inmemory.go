package admin

import (
	"context"
	"sync"
	"time"
)

type InMemoryService struct {
	mu       sync.Mutex
	disputes map[string]Dispute
	payouts  map[string]Payout
	alerts   map[string]FraudAlert
	audits   map[string]AuditLog
	notes    map[string]Notification
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		disputes: make(map[string]Dispute),
		payouts:  make(map[string]Payout),
		alerts:   make(map[string]FraudAlert),
		audits:   make(map[string]AuditLog),
		notes:    make(map[string]Notification),
	}
}

func (s *InMemoryService) GetDashboard(ctx context.Context) (DashboardSummary, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	open := 0
	pending := 0
	for _, d := range s.disputes {
		if d.Status == DisputeOpen {
			open++
		}
	}
	for _, p := range s.payouts {
		if p.Status == PayoutPending {
			pending++
		}
	}
	return DashboardSummary{OpenDisputes: open, PendingPayouts: pending, FraudAlerts: len(s.alerts)}, nil
}

func (s *InMemoryService) ListDisputes(ctx context.Context) ([]Dispute, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Dispute, 0, len(s.disputes))
	for _, d := range s.disputes {
		out = append(out, d)
	}
	return out, nil
}

func (s *InMemoryService) CreateDispute(ctx context.Context, d Dispute) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d.ID == "" {
		d.ID = time.Now().Format("20060102150405")
	}
	d.CreatedAt = time.Now()
	s.disputes[d.ID] = d
	return nil
}

func (s *InMemoryService) ApprovePayout(ctx context.Context, payoutID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.payouts[payoutID]
	if !ok {
		return nil
	}
	p.Status = PayoutPaid
	s.payouts[payoutID] = p
	return nil
}

func (s *InMemoryService) RequestPayoutApproval(ctx context.Context, payoutID string, requestedBy int64, notes string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := "pa-" + payoutID
	s.audits[id] = AuditLog{ID: id, ActorID: requestedBy, Action: "request_payout_approval", ResourceType: "payout", ResourceID: payoutID, Details: map[string]interface{}{"notes": notes}, CreatedAt: time.Now()}
	return nil
}

func (s *InMemoryService) ListAuditLogs(ctx context.Context) ([]AuditLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]AuditLog, 0, len(s.audits))
	for _, a := range s.audits {
		out = append(out, a)
	}
	return out, nil
}

func (s *InMemoryService) ListNotifications(ctx context.Context) ([]Notification, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Notification, 0, len(s.notes))
	for _, n := range s.notes {
		out = append(out, n)
	}
	return out, nil
}

func (s *InMemoryService) ListPayouts(ctx context.Context) ([]Payout, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Payout, 0, len(s.payouts))
	for _, p := range s.payouts {
		out = append(out, p)
	}
	return out, nil
}

func (s *InMemoryService) ListFraudAlerts(ctx context.Context) ([]FraudAlert, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]FraudAlert, 0, len(s.alerts))
	for _, a := range s.alerts {
		out = append(out, a)
	}
	return out, nil
}

func (s *InMemoryService) VerifyCreator(ctx context.Context, v CreatorVerification) error {
	// In-memory is a no-op that accepts verification
	return nil
}

func (s *InMemoryService) Refund(ctx context.Context, disputeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	d, ok := s.disputes[disputeID]
	if !ok {
		return nil
	}
	d.Status = DisputeClosed
	s.disputes[disputeID] = d
	return nil
}

func (s *InMemoryService) SuspendAccount(ctx context.Context, userID string) error {
	// no-op for in-memory
	return nil
}
