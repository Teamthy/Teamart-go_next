-- queries/settlement.sql
-- Seller settlement and payout queries

-- ===== SELLER SETTLEMENTS =====

-- name: CreateSellerSettlement :one
INSERT INTO seller_settlements (
    merchant_id, settlement_number, period_start, period_end,
    gross_sales_cents, order_count, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetSellerSettlement :one
SELECT * FROM seller_settlements WHERE id = $1;

-- name: GetSettlementByNumber :one
SELECT * FROM seller_settlements WHERE settlement_number = $1;

-- name: ListMerchantSettlements :many
SELECT * FROM seller_settlements 
WHERE merchant_id = $1
ORDER BY period_end DESC
LIMIT $2 OFFSET $3;

-- name: ListSettlementsByStatus :many
SELECT * FROM seller_settlements 
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPendingSettlements :many
SELECT * FROM seller_settlements 
WHERE status IN ('pending', 'reviewed', 'approved', 'processing')
ORDER BY period_end ASC;

-- name: GetLatestSettlement :one
SELECT * FROM seller_settlements 
WHERE merchant_id = $1
ORDER BY period_end DESC
LIMIT 1;

-- name: UpdateSettlementStatus :one
UPDATE seller_settlements 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateSettlementAmounts :one
UPDATE seller_settlements 
SET platform_fees_cents = $2, return_refunds_cents = $3,
    chargebacks_cents = $4, reversals_cents = $5, 
    net_payout_cents = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ApproveSettlement :one
UPDATE seller_settlements 
SET status = 'approved', updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ScheduleSettlement :one
UPDATE seller_settlements 
SET status = 'processing', scheduled_payout_date = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CompleteSettlement :one
UPDATE seller_settlements 
SET status = 'paid', actual_payout_date = CURRENT_DATE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== PAYOUTS =====

-- name: CreatePayout :one
INSERT INTO payouts (
    merchant_id, settlement_id, amount_cents, currency, status,
    payment_method, payment_details, provider
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPayout :one
SELECT * FROM payouts WHERE id = $1;

-- name: ListMerchantPayouts :many
SELECT * FROM payouts 
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPayoutsByStatus :many
SELECT * FROM payouts 
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPendingPayouts :many
SELECT * FROM payouts 
WHERE status IN ('pending', 'scheduled')
ORDER BY scheduled_at ASC;

-- name: ListSettlementPayouts :many
SELECT * FROM payouts 
WHERE settlement_id = $1
ORDER BY created_at DESC;

-- name: UpdatePayoutStatus :one
UPDATE payouts 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: MarkPayoutSent :one
UPDATE payouts 
SET status = 'in_transit', sent_at = CURRENT_TIMESTAMP,
    provider_reference = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: MarkPayoutCompleted :one
UPDATE payouts 
SET status = 'completed', completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: MarkPayoutFailed :one
UPDATE payouts 
SET status = 'failed', failed_at = CURRENT_TIMESTAMP, failure_reason = $2,
    retry_count = retry_count + 1, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== ANALYTICS =====

-- name: GetSettlementStats :one
SELECT 
    COUNT(*) as total_settlements,
    SUM(gross_sales_cents) as total_gross_sales,
    SUM(net_payout_cents) as total_net_payouts,
    AVG(net_payout_cents) as avg_payout_amount,
    COUNT(CASE WHEN status = 'paid' THEN 1 END) as completed_settlements
FROM seller_settlements 
WHERE created_at BETWEEN $1 AND $2;

-- name: GetPayoutStats :one
SELECT 
    COUNT(*) as total_payouts,
    SUM(amount_cents) as total_payout_amount,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_payouts,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_payouts,
    AVG(amount_cents) as avg_payout_amount
FROM payouts 
WHERE created_at BETWEEN $1 AND $2;

-- name: GetMerchantSettlementHistory :many
SELECT 
    period_start,
    period_end,
    gross_sales_cents,
    platform_fees_cents,
    return_refunds_cents,
    net_payout_cents,
    status
FROM seller_settlements 
WHERE merchant_id = $1
ORDER BY period_end DESC
LIMIT $2;
