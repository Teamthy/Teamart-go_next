package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

type AdminRepository interface {
	SaveDispute(ctx context.Context, pool *database.Pool, d Dispute) error
	UpdateDisputeStatus(ctx context.Context, pool *database.Pool, id string, status string) error
	ListDisputes(ctx context.Context, pool *database.Pool) ([]Dispute, error)

	SavePayout(ctx context.Context, pool *database.Pool, p Payout) error
	UpdatePayoutStatus(ctx context.Context, pool *database.Pool, id string, status string) error
	ListPayouts(ctx context.Context, pool *database.Pool) ([]Payout, error)

	SaveFraudAlert(ctx context.Context, pool *database.Pool, a FraudAlert) error
	ListFraudAlerts(ctx context.Context, pool *database.Pool) ([]FraudAlert, error)
	// audit & notifications
	SaveAuditLog(ctx context.Context, pool *database.Pool, a AuditLog) error
	ListAuditLogs(ctx context.Context, pool *database.Pool) ([]AuditLog, error)
	SaveNotification(ctx context.Context, pool *database.Pool, n Notification) error
	ListNotifications(ctx context.Context, pool *database.Pool) ([]Notification, error)
	// payout approvals
	CreatePayoutApproval(ctx context.Context, pool *database.Pool, pa PayoutApproval) error
	ApprovePayoutApproval(ctx context.Context, pool *database.Pool, approvalID string, approverID int64) error
}

type PostgresAdminRepository struct {
	log *logger.Logger
}

func NewPostgresAdminRepository(log *logger.Logger) *PostgresAdminRepository {
	return &PostgresAdminRepository{log: log}
}

func (r *PostgresAdminRepository) SaveDispute(ctx context.Context, pool *database.Pool, d Dispute) error {
	sqlStr := `INSERT INTO disputes (id, order_id, user_id, amount, reason, status, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6, now(), now())
ON CONFLICT (id) DO UPDATE SET order_id = EXCLUDED.order_id, user_id = EXCLUDED.user_id, amount = EXCLUDED.amount, reason = EXCLUDED.reason, status = EXCLUDED.status, updated_at = now()`
	_, err := pool.Exec(ctx, sqlStr, d.ID, d.OrderID, d.UserID, d.Amount, d.Reason, string(d.Status))
	if err != nil {
		return fmt.Errorf("save dispute: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) UpdateDisputeStatus(ctx context.Context, pool *database.Pool, id string, status string) error {
	_, err := pool.Exec(ctx, `UPDATE disputes SET status=$1, updated_at = now() WHERE id=$2`, status, id)
	if err != nil {
		return fmt.Errorf("update dispute status: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) ListDisputes(ctx context.Context, pool *database.Pool) ([]Dispute, error) {
	rows, err := pool.Query(ctx, `SELECT id, order_id, user_id, amount, reason, status, created_at FROM disputes`)
	if err != nil {
		return nil, fmt.Errorf("list disputes: %w", err)
	}
	defer rows.Close()
	var out []Dispute
	for rows.Next() {
		var d Dispute
		var status sql.NullString
		if err := rows.Scan(&d.ID, &d.OrderID, &d.UserID, &d.Amount, &d.Reason, &status, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan dispute: %w", err)
		}
		if status.Valid {
			d.Status = DisputeStatus(status.String)
		}
		out = append(out, d)
	}
	return out, nil
}

func (r *PostgresAdminRepository) SavePayout(ctx context.Context, pool *database.Pool, p Payout) error {
	_, err := pool.Exec(ctx, `INSERT INTO payouts (id, creator_id, amount, status, created_at, updated_at) VALUES ($1,$2,$3,$4, now(), now()) ON CONFLICT (id) DO UPDATE SET creator_id=EXCLUDED.creator_id, amount=EXCLUDED.amount, status=EXCLUDED.status, updated_at=now()`, p.ID, p.CreatorID, p.Amount, string(p.Status))
	if err != nil {
		return fmt.Errorf("save payout: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) UpdatePayoutStatus(ctx context.Context, pool *database.Pool, id string, status string) error {
	_, err := pool.Exec(ctx, `UPDATE payouts SET status=$1, updated_at=now() WHERE id=$2`, status, id)
	if err != nil {
		return fmt.Errorf("update payout status: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) ListPayouts(ctx context.Context, pool *database.Pool) ([]Payout, error) {
	rows, err := pool.Query(ctx, `SELECT id, creator_id, amount, status, created_at FROM payouts`)
	if err != nil {
		return nil, fmt.Errorf("list payouts: %w", err)
	}
	defer rows.Close()
	var out []Payout
	for rows.Next() {
		var p Payout
		var status sql.NullString
		if err := rows.Scan(&p.ID, &p.CreatorID, &p.Amount, &status, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan payout: %w", err)
		}
		if status.Valid {
			p.Status = PayoutStatus(status.String)
		}
		out = append(out, p)
	}
	return out, nil
}

func (r *PostgresAdminRepository) SaveFraudAlert(ctx context.Context, pool *database.Pool, a FraudAlert) error {
	dataB, err := json.Marshal(a.Data)
	if err != nil {
		return fmt.Errorf("marshal alert data: %w", err)
	}
	_, err = pool.Exec(ctx, `INSERT INTO fraud_alerts (id, subject_id, score, data, created_at) VALUES ($1,$2,$3,$4::jsonb, now()) ON CONFLICT (id) DO UPDATE SET subject_id=EXCLUDED.subject_id, score=EXCLUDED.score, data=EXCLUDED.data`, a.ID, a.SubjectID, a.Score, dataB)
	if err != nil {
		return fmt.Errorf("save alert: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) ListFraudAlerts(ctx context.Context, pool *database.Pool) ([]FraudAlert, error) {
	rows, err := pool.Query(ctx, `SELECT id, subject_id, score, data, created_at FROM fraud_alerts`)
	if err != nil {
		return nil, fmt.Errorf("list alerts: %w", err)
	}
	defer rows.Close()
	var out []FraudAlert
	for rows.Next() {
		var a FraudAlert
		var dataB []byte
		if err := rows.Scan(&a.ID, &a.SubjectID, &a.Score, &dataB, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan alert: %w", err)
		}
		var data map[string]interface{}
		_ = json.Unmarshal(dataB, &data)
		a.Data = data
		out = append(out, a)
	}
	return out, nil
}

// Audit logs
func (r *PostgresAdminRepository) SaveAuditLog(ctx context.Context, pool *database.Pool, a AuditLog) error {
	detailsB, _ := json.Marshal(a.Details)
	_, err := pool.Exec(ctx, `INSERT INTO audit_logs (id, actor_id, action, resource_type, resource_id, details, created_at) VALUES ($1,$2,$3,$4,$5,$6::jsonb, now()) ON CONFLICT (id) DO NOTHING`, a.ID, a.ActorID, a.Action, a.ResourceType, a.ResourceID, detailsB)
	if err != nil {
		return fmt.Errorf("save audit: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) ListAuditLogs(ctx context.Context, pool *database.Pool) ([]AuditLog, error) {
	rows, err := pool.Query(ctx, `SELECT id, actor_id, action, resource_type, resource_id, details, created_at FROM audit_logs`)
	if err != nil {
		return nil, fmt.Errorf("list audit: %w", err)
	}
	defer rows.Close()
	var out []AuditLog
	for rows.Next() {
		var a AuditLog
		var detailsB []byte
		if err := rows.Scan(&a.ID, &a.ActorID, &a.Action, &a.ResourceType, &a.ResourceID, &detailsB, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan audit: %w", err)
		}
		var details map[string]interface{}
		_ = json.Unmarshal(detailsB, &details)
		a.Details = details
		out = append(out, a)
	}
	return out, nil
}

// Notifications
func (r *PostgresAdminRepository) SaveNotification(ctx context.Context, pool *database.Pool, n Notification) error {
	payloadB, _ := json.Marshal(n.Payload)
	_, err := pool.Exec(ctx, `INSERT INTO notifications (id, recipient_id, channel, payload, sent_at) VALUES ($1,$2,$3,$4::jsonb,$5) ON CONFLICT (id) DO NOTHING`, n.ID, n.RecipientID, n.Channel, payloadB, n.SentAt)
	if err != nil {
		return fmt.Errorf("save notification: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) ListNotifications(ctx context.Context, pool *database.Pool) ([]Notification, error) {
	rows, err := pool.Query(ctx, `SELECT id, recipient_id, channel, payload, sent_at FROM notifications`)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()
	var out []Notification
	for rows.Next() {
		var n Notification
		var payloadB []byte
		if err := rows.Scan(&n.ID, &n.RecipientID, &n.Channel, &payloadB, &n.SentAt); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		var payload map[string]interface{}
		_ = json.Unmarshal(payloadB, &payload)
		n.Payload = payload
		out = append(out, n)
	}
	return out, nil
}

// Payout approvals
func (r *PostgresAdminRepository) CreatePayoutApproval(ctx context.Context, pool *database.Pool, pa PayoutApproval) error {
	_, err := pool.Exec(ctx, `INSERT INTO payout_approvals (id, payout_id, requested_by, status, notes, created_at) VALUES ($1,$2,$3,$4,$5, now()) ON CONFLICT (id) DO NOTHING`, pa.ID, pa.PayoutID, pa.RequestedBy, pa.Status, pa.Notes)
	if err != nil {
		return fmt.Errorf("create payout approval: %w", err)
	}
	return nil
}

func (r *PostgresAdminRepository) ApprovePayoutApproval(ctx context.Context, pool *database.Pool, approvalID string, approverID int64) error {
	_, err := pool.Exec(ctx, `UPDATE payout_approvals SET status='approved', approved_by=$1, approved_at=now() WHERE id=$2`, approverID, approvalID)
	if err != nil {
		return fmt.Errorf("approve payout approval: %w", err)
	}
	return nil
}
