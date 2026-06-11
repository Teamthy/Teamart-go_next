-- queries/fraud.sql
-- Fraud detection and compliance queries

-- name: CreateFraudFlag :one
INSERT INTO fraud_flags (
    order_id, risk_score, risk_level, flags, status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetFraudFlag :one
SELECT * FROM fraud_flags WHERE id = $1;

-- name: GetOrderFraudFlag :one
SELECT * FROM fraud_flags WHERE order_id = $1;

-- name: ListFlaggedOrders :many
SELECT ff.* FROM fraud_flags ff
WHERE status = 'flagged'
ORDER BY ff.risk_level DESC, ff.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListHighRiskOrders :many
SELECT ff.* FROM fraud_flags ff
WHERE risk_level IN ('high', 'critical')
ORDER BY ff.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListReviewedFlags :many
SELECT * FROM fraud_flags 
WHERE status IN ('reviewed', 'approved')
ORDER BY reviewed_at DESC
LIMIT $1 OFFSET $2;

-- name: ListUnreviewedFlags :many
SELECT * FROM fraud_flags 
WHERE status = 'flagged'
ORDER BY created_at ASC;

-- name: UpdateFraudFlagStatus :one
UPDATE fraud_flags 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ReviewFraudFlag :one
UPDATE fraud_flags 
SET status = 'reviewed', reviewed_by = $2, reviewed_at = CURRENT_TIMESTAMP,
    review_notes = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ApproveFraudFlag :one
UPDATE fraud_flags 
SET status = 'approved', reviewed_by = $2, reviewed_at = CURRENT_TIMESTAMP,
    review_notes = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: RejectFraudFlag :one
UPDATE fraud_flags 
SET status = 'rejected', reviewed_by = $2, reviewed_at = CURRENT_TIMESTAMP,
    review_notes = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListOrdersByRiskScore :many
SELECT ff.*, o.order_number, o.total
FROM fraud_flags ff
JOIN orders o ON ff.order_id = o.id
WHERE ff.risk_score >= $1 AND ff.created_at BETWEEN $2 AND $3
ORDER BY ff.risk_score DESC
LIMIT $4 OFFSET $5;

-- ===== ANALYTICS =====

-- name: GetFraudStats :one
SELECT 
    COUNT(*) as total_flagged_orders,
    COUNT(CASE WHEN risk_level = 'low' THEN 1 END) as low_risk_count,
    COUNT(CASE WHEN risk_level = 'medium' THEN 1 END) as medium_risk_count,
    COUNT(CASE WHEN risk_level = 'high' THEN 1 END) as high_risk_count,
    COUNT(CASE WHEN risk_level = 'critical' THEN 1 END) as critical_risk_count,
    COUNT(CASE WHEN status = 'approved' THEN 1 END) as approved_as_legitimate,
    COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected_as_fraud
FROM fraud_flags 
WHERE created_at BETWEEN $1 AND $2;

-- name: GetFraudReviewMetrics :one
SELECT 
    COUNT(CASE WHEN status = 'flagged' THEN 1 END) as pending_review,
    COUNT(CASE WHEN status = 'reviewed' THEN 1 END) as under_investigation,
    COUNT(CASE WHEN status IN ('approved', 'rejected') THEN 1 END) as resolved,
    ROUND(AVG(EXTRACT(EPOCH FROM (reviewed_at - created_at)))::NUMERIC / 3600, 2) as avg_review_time_hours
FROM fraud_flags 
WHERE created_at BETWEEN $1 AND $2;

-- name: GetRiskScoreDistribution :many
SELECT 
    CASE 
        WHEN risk_score < 25 THEN 'low'
        WHEN risk_score < 50 THEN 'medium'
        WHEN risk_score < 75 THEN 'high'
        ELSE 'critical'
    END as risk_category,
    COUNT(*) as count,
    ROUND(100.0 * COUNT(*) / SUM(COUNT(*)) OVER (), 2) as percentage
FROM fraud_flags 
WHERE created_at BETWEEN $1 AND $2
GROUP BY risk_category
ORDER BY risk_score ASC;
