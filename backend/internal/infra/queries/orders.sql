-- queries/orders.sql
-- Comprehensive order management queries

-- name: CreateOrder :one
INSERT INTO orders (
    order_number, user_id, source, creator_id, customer_email, customer_phone,
    status, subtotal, shipping_cost, tax_amount, discount_amount, total,
    payment_method, payment_status, billing_address, shipping_address, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders WHERE id = $1 AND deleted_at IS NULL;

-- name: GetOrderByNumber :one
SELECT * FROM orders WHERE order_number = $1 AND deleted_at IS NULL;

-- name: GetOrdersByUser :many
SELECT * FROM orders 
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrdersByMerchant :many
SELECT * FROM orders 
WHERE creator_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrdersByStatus :many
SELECT * FROM orders 
WHERE status = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrdersByStatusAndMerchant :many
SELECT * FROM orders 
WHERE creator_id = $1 AND status = $2 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateOrderStatus :one
UPDATE orders 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateOrderPaymentStatus :one
UPDATE orders 
SET payment_status = $2, paid_amount = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateOrderTracking :one
UPDATE orders 
SET tracking_number = $2, shipped_at = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CancelOrder :one
UPDATE orders 
SET status = 'cancelled', cancelled_at = CURRENT_TIMESTAMP, 
    cancellation_reason = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteOrder :one
UPDATE orders 
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== ANALYTICS =====

-- name: CountOrdersByStatus :many
SELECT 
    status,
    COUNT(*) as count
FROM orders 
WHERE deleted_at IS NULL
GROUP BY status;

-- name: GetOrderMetrics :one
SELECT 
    COUNT(*) as total_orders,
    COUNT(DISTINCT user_id) as unique_customers,
    SUM(total) as total_revenue,
    AVG(total) as avg_order_value,
    MIN(total) as min_order_value,
    MAX(total) as max_order_value
FROM orders 
WHERE created_at BETWEEN $1 AND $2 AND deleted_at IS NULL;

-- name: GetMerchantOrderMetrics :one
SELECT 
    COUNT(*) as total_orders,
    SUM(total) as total_revenue,
    AVG(total) as avg_order_value,
    COUNT(CASE WHEN status IN ('delivered', 'completed') THEN 1 END) as completed_orders
FROM orders 
WHERE creator_id = $1 AND created_at BETWEEN $2 AND $3 AND deleted_at IS NULL;

-- name: CountOrdersByStatus :one
SELECT COUNT(*) FROM orders WHERE status = $1;
