package admin

import (
	"context"
	"time"

	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

// PostgresAdminService implements Service using an AdminRepository and DB pool
type PostgresAdminService struct {
	repo AdminRepository
	pool *database.Pool
	log  *logger.Logger
}

func NewPostgresAdminService(repo AdminRepository, pool *database.Pool, log *logger.Logger) *PostgresAdminService {
	return &PostgresAdminService{repo: repo, pool: pool, log: log}
}

func (s *PostgresAdminService) GetDashboard(ctx context.Context) (DashboardSummary, error) {
	disputes, _ := s.repo.ListDisputes(ctx, s.pool)
	payouts, _ := s.repo.ListPayouts(ctx, s.pool)
	alerts, _ := s.repo.ListFraudAlerts(ctx, s.pool)
	open := 0
	pending := 0
	for _, d := range disputes {
		if d.Status == DisputeOpen {
			open++
		}
	}
	for _, p := range payouts {
		if p.Status == PayoutPending {
			pending++
		}
	}
	return DashboardSummary{OpenDisputes: open, PendingPayouts: pending, FraudAlerts: len(alerts)}, nil
}

func (s *PostgresAdminService) ListDisputes(ctx context.Context) ([]Dispute, error) {
	return s.repo.ListDisputes(ctx, s.pool)
}

func (s *PostgresAdminService) CreateDispute(ctx context.Context, d Dispute) error {
	return s.repo.SaveDispute(ctx, s.pool, d)
}
func (s *PostgresAdminService) ListPayouts(ctx context.Context) ([]Payout, error) {
	return s.repo.ListPayouts(ctx, s.pool)
}

func (s *PostgresAdminService) ListFraudAlerts(ctx context.Context) ([]FraudAlert, error) {
	return s.repo.ListFraudAlerts(ctx, s.pool)
}

func (s *PostgresAdminService) VerifyCreator(ctx context.Context, v CreatorVerification) error {
	// Not implemented in Postgres repository for now
	return nil
}

func (s *PostgresAdminService) Refund(ctx context.Context, disputeID string) error {
	// Update dispute
	if err := s.repo.UpdateDisputeStatus(ctx, s.pool, disputeID, string(DisputeClosed)); err != nil {
		return err
	}
	// Audit
	_ = s.repo.SaveAuditLog(ctx, s.pool, AuditLog{ID: "audit-refund-" + disputeID, ActorID: 0, Action: "refund", ResourceType: "dispute", ResourceID: disputeID, Details: map[string]interface{}{}, CreatedAt: time.Now()})
	// Notify support (placeholder recipient)
	_ = s.repo.SaveNotification(ctx, s.pool, Notification{ID: "note-refund-" + disputeID, RecipientID: "support", Channel: "email", Payload: map[string]interface{}{"dispute_id": disputeID}, SentAt: nil})
	return nil
}

func (s *PostgresAdminService) RequestPayoutApproval(ctx context.Context, payoutID string, requestedBy int64, notes string) error {
	pa := PayoutApproval{ID: "pa-" + payoutID, PayoutID: payoutID, RequestedBy: requestedBy, Status: "requested", Notes: notes, CreatedAt: time.Now()}
	if err := s.repo.CreatePayoutApproval(ctx, s.pool, pa); err != nil {
		return err
	}
	// Audit
	_ = s.repo.SaveAuditLog(ctx, s.pool, AuditLog{ID: "audit-request-" + pa.ID, ActorID: requestedBy, Action: "request_payout_approval", ResourceType: "payout_approval", ResourceID: pa.ID, Details: map[string]interface{}{"notes": notes}, CreatedAt: time.Now()})
	return nil
}

func (s *PostgresAdminService) ApprovePayout(ctx context.Context, payoutID string, approverID int64) error {
	// Find approval id (simple derivation used earlier)
	approvalID := "pa-" + payoutID
	if err := s.repo.ApprovePayoutApproval(ctx, s.pool, approvalID, approverID); err != nil {
		return err
	}
	if err := s.repo.UpdatePayoutStatus(ctx, s.pool, payoutID, string(PayoutPaid)); err != nil {
		return err
	}
	// Audit & notify
	_ = s.repo.SaveAuditLog(ctx, s.pool, AuditLog{ID: "audit-approve-" + approvalID, ActorID: approverID, Action: "approve_payout", ResourceType: "payout", ResourceID: payoutID, Details: map[string]interface{}{}, CreatedAt: time.Now()})
	_ = s.repo.SaveNotification(ctx, s.pool, Notification{ID: "note-approve-" + payoutID, RecipientID: "finance", Channel: "email", Payload: map[string]interface{}{"payout_id": payoutID}, SentAt: nil})
	return nil
}

func (s *PostgresAdminService) ListAuditLogs(ctx context.Context) ([]AuditLog, error) {
	return s.repo.ListAuditLogs(ctx, s.pool)
}

func (s *PostgresAdminService) ListNotifications(ctx context.Context) ([]Notification, error) {
	return s.repo.ListNotifications(ctx, s.pool)
}

func (s *PostgresAdminService) SuspendAccount(ctx context.Context, userID string) error {
	// Account suspension is out of scope for repository; no-op here
	return nil
}
