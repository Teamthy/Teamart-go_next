-- queries/returns.sql
-- Returns and refunds queries

-- ===== RETURNS =====

-- name: CreateReturn :one
INSERT INTO returns (
    order_id, return_number, return_type, reason, description,
    requested_by, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetReturn :one
SELECT * FROM returns WHERE id = $1 AND deleted_at IS NULL;

-- name: GetReturnByNumber :one
SELECT * FROM returns WHERE return_number = $1 AND deleted_at IS NULL;

-- name: ListReturnsByOrder :many
SELECT * FROM returns 
WHERE order_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListReturnsByCustomer :many
SELECT * FROM returns 
WHERE requested_by = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListReturnsByStatus :many
SELECT * FROM returns 
WHERE status = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListReturnsByType :many
SELECT * FROM returns 
WHERE return_type = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateReturnStatus :one
UPDATE returns 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ApproveReturn :one
UPDATE returns 
SET status = 'approved', approved_by = $2, approved_at = CURRENT_TIMESTAMP,
    approval_notes = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RejectReturn :one
UPDATE returns 
SET status = 'rejected', rejected_by = $2, rejected_at = CURRENT_TIMESTAMP,
    rejection_reason = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RecordReturnReceipt :one
UPDATE returns 
SET status = 'received', received_at = CURRENT_TIMESTAMP,
    received_by = $2, inspection_notes = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: ResolveReturn :one
UPDATE returns 
SET status = 'resolved', resolution_type = $2, resolved_by = $3,
    resolved_at = CURRENT_TIMESTAMP, resolution_notes = $4, 
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- ===== REFUNDS =====

-- name: CreateRefund :one
INSERT INTO refunds (
    payment_intent_id, return_id, amount_cents, reason, description, status
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetRefund :one
SELECT * FROM refunds WHERE id = $1;

-- name: ListRefundsByPaymentIntent :many
SELECT * FROM refunds 
WHERE payment_intent_id = $1 
ORDER BY created_at DESC;

-- name: ListRefundsByReturn :many
SELECT * FROM refunds 
WHERE return_id = $1 
ORDER BY created_at DESC;

-- name: ListRefundsByStatus :many
SELECT * FROM refunds 
WHERE status = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateRefundStatus :one
UPDATE refunds 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateRefundProcessed :one
UPDATE refunds 
SET status = 'completed', provider_refund_id = $2, 
    processed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateRefundFailed :one
UPDATE refunds 
SET status = 'failed', failure_reason = $2, retry_count = retry_count + 1,
    last_retry_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: GetRefundsByDateRange :many
SELECT * FROM refunds 
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- ===== ANALYTICS =====

-- name: GetReturnStats :one
SELECT 
    COUNT(DISTINCT order_id) as orders_with_returns,
    COUNT(*) as total_returns,
    COUNT(CASE WHEN status = 'approved' THEN 1 END) as approved_returns,
    COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected_returns,
    COUNT(CASE WHEN status = 'resolved' THEN 1 END) as resolved_returns
FROM returns 
WHERE created_at BETWEEN $1 AND $2 AND deleted_at IS NULL;

-- name: GetRefundStats :one
SELECT 
    COUNT(*) as total_refunds,
    SUM(amount_cents) as total_refund_amount,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_refunds,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_refunds,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_refunds,
    AVG(amount_cents) as avg_refund_amount
FROM refunds 
WHERE created_at BETWEEN $1 AND $2;
