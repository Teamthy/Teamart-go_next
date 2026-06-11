-- queries/payments.sql
-- Payment intent queries for checkout and payment processing

-- name: CreatePaymentIntent :one
INSERT INTO payment_intents (
    order_id, external_id, provider, method, amount_cents, 
    currency, status, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetPaymentIntent :one
SELECT * FROM payment_intents WHERE id = $1;

-- name: GetPaymentIntentByOrderID :one
SELECT * FROM payment_intents WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1;

-- name: GetPaymentIntentByExternalID :one
SELECT * FROM payment_intents WHERE external_id = $1;

-- name: UpdatePaymentIntentStatus :one
UPDATE payment_intents 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdatePaymentIntentAuthorized :one
UPDATE payment_intents 
SET status = 'authorized', authorized_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdatePaymentIntentCaptured :one
UPDATE payment_intents 
SET status = 'captured', captured_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdatePaymentIntentFailed :one
UPDATE payment_intents 
SET status = 'failed', failed_at = CURRENT_TIMESTAMP, failure_reason = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListPaymentIntentsByStatus :many
SELECT * FROM payment_intents 
WHERE status = $1 
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListPaymentIntentsByProvider :many
SELECT * FROM payment_intents 
WHERE provider = $1 AND created_at > $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- ===== PAYMENT TRANSACTIONS =====

-- name: CreatePaymentTransaction :one
INSERT INTO payment_transactions (
    payment_intent_id, transaction_type, provider_transaction_id, 
    amount_cents, status, provider_response
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetPaymentTransaction :one
SELECT * FROM payment_transactions WHERE id = $1;

-- name: ListPaymentTransactions :many
SELECT * FROM payment_transactions 
WHERE payment_intent_id = $1 
ORDER BY created_at DESC;

-- name: ListPaymentTransactionsByType :many
SELECT * FROM payment_transactions 
WHERE payment_intent_id = $1 AND transaction_type = $2 
ORDER BY created_at DESC;

-- ===== PAYMENT AUDITS =====

-- name: CreatePaymentAudit :one
INSERT INTO payment_audits (
    payment_intent_id, action, actor_id, actor_type, message, 
    payload, ip_address
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ListPaymentAudits :many
SELECT * FROM payment_audits 
WHERE payment_intent_id = $1 
ORDER BY created_at DESC;

-- name: ListPaymentAuditsByAction :many
SELECT * FROM payment_audits 
WHERE payment_intent_id = $1 AND action = $2 
ORDER BY created_at DESC;

-- name: ListPaymentAuditsByDate :many
SELECT * FROM payment_audits 
WHERE payment_intent_id = $1 AND created_at BETWEEN $2 AND $3
ORDER BY created_at DESC;

-- ===== ANALYTICS QUERIES =====

-- name: GetPaymentStats :one
SELECT 
    COUNT(*) as total_payments,
    SUM(amount_cents) as total_amount_cents,
    COUNT(CASE WHEN status = 'captured' THEN 1 END) as successful_count,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_count
FROM payment_intents 
WHERE created_at BETWEEN $1 AND $2;

-- name: GetPaymentStatsByProvider :many
SELECT 
    provider,
    COUNT(*) as payment_count,
    SUM(amount_cents) as total_amount_cents,
    COUNT(CASE WHEN status = 'captured' THEN 1 END) as successful_count
FROM payment_intents 
WHERE created_at BETWEEN $1 AND $2
GROUP BY provider
ORDER BY total_amount_cents DESC;
