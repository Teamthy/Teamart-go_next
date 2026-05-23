package admin

import "context"

// Service defines admin operations used by handlers and tooling.
type Service interface {
	GetDashboard(ctx context.Context) (DashboardSummary, error)
	ListDisputes(ctx context.Context) ([]Dispute, error)
	CreateDispute(ctx context.Context, d Dispute) error
	RequestPayoutApproval(ctx context.Context, payoutID string, requestedBy int64, notes string) error
	ApprovePayout(ctx context.Context, payoutID string, approverID int64) error
	ListPayouts(ctx context.Context) ([]Payout, error)
	ListFraudAlerts(ctx context.Context) ([]FraudAlert, error)
	VerifyCreator(ctx context.Context, v CreatorVerification) error
	Refund(ctx context.Context, disputeID string) error
	SuspendAccount(ctx context.Context, userID string) error

	// Audit and notifications
	ListAuditLogs(ctx context.Context) ([]AuditLog, error)
	ListNotifications(ctx context.Context) ([]Notification, error)
}
