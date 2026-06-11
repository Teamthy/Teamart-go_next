-- queries/analytics.sql
-- Analytics and reporting queries

-- ===== ORDER ANALYTICS =====

-- name: CreateOrderAnalytic :one
INSERT INTO order_analytics (
    merchant_id, date_bucket, orders_created, orders_completed,
    gross_revenue_cents, items_sold
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetOrderAnalytic :one
SELECT * FROM order_analytics 
WHERE merchant_id = $1 AND date_bucket = $2;

-- name: ListMerchantAnalytics :many
SELECT * FROM order_analytics 
WHERE merchant_id = $1 AND date_bucket BETWEEN $2 AND $3
ORDER BY date_bucket DESC;

-- name: UpdateOrderAnalytic :one
UPDATE order_analytics 
SET orders_created = $2, orders_completed = $3, gross_revenue_cents = $4,
    items_sold = $5, updated_at = CURRENT_TIMESTAMP
WHERE merchant_id = $1 AND date_bucket = $2
RETURNING *;

-- ===== CONSOLIDATED METRICS =====

-- name: GetDailyMerchantMetrics :one
SELECT 
    DATE(created_at) as date,
    COUNT(*) as orders_count,
    SUM(total) as total_revenue_cents,
    AVG(total) as avg_order_value_cents,
    COUNT(DISTINCT user_id) as unique_customers,
    COUNT(CASE WHEN status IN ('cancelled', 'refunded') THEN 1 END) as failed_orders
FROM orders 
WHERE creator_id = $1 AND DATE(created_at) = $2 AND deleted_at IS NULL
GROUP BY DATE(created_at);

-- name: GetWeeklyMerchantMetrics :many
SELECT 
    DATE_TRUNC('week', created_at) as week_start,
    COUNT(*) as orders_count,
    SUM(total) as total_revenue_cents,
    AVG(total) as avg_order_value_cents,
    COUNT(DISTINCT user_id) as unique_customers
FROM orders 
WHERE creator_id = $1 AND created_at BETWEEN $2 AND $3 AND deleted_at IS NULL
GROUP BY DATE_TRUNC('week', created_at)
ORDER BY week_start DESC;

-- name: GetMonthlyMerchantMetrics :many
SELECT 
    DATE_TRUNC('month', created_at) as month_start,
    COUNT(*) as orders_count,
    SUM(total) as total_revenue_cents,
    AVG(total) as avg_order_value_cents,
    COUNT(DISTINCT user_id) as unique_customers
FROM orders 
WHERE creator_id = $1 AND created_at BETWEEN $2 AND $3 AND deleted_at IS NULL
GROUP BY DATE_TRUNC('month', created_at)
ORDER BY month_start DESC;

-- ===== CONVERSION ANALYTICS =====

-- name: GetCheckoutConversion :one
SELECT 
    COUNT(*) as total_users,
    COUNT(DISTINCT user_id) as users_with_orders,
    ROUND(100.0 * COUNT(DISTINCT user_id) / NULLIF(COUNT(DISTINCT user_id), 0), 2) as conversion_rate
FROM orders 
WHERE created_at BETWEEN $1 AND $2 AND deleted_at IS NULL;

-- name: GetReturnRate :one
SELECT 
    COUNT(*) as total_orders,
    COUNT(DISTINCT r.order_id) as orders_with_returns,
    ROUND(100.0 * COUNT(DISTINCT r.order_id) / NULLIF(COUNT(*), 0), 2) as return_rate
FROM orders o
LEFT JOIN returns r ON o.id = r.order_id AND r.deleted_at IS NULL
WHERE o.created_at BETWEEN $1 AND $2 AND o.deleted_at IS NULL;

-- name: GetRefundRate :one
SELECT 
    COUNT(*) as total_orders,
    COUNT(CASE WHEN o.status IN ('cancelled', 'refunded') THEN 1 END) as refunded_orders,
    ROUND(100.0 * COUNT(CASE WHEN o.status IN ('cancelled', 'refunded') THEN 1 END) / NULLIF(COUNT(*), 0), 2) as refund_rate
FROM orders o
WHERE o.created_at BETWEEN $1 AND $2 AND o.deleted_at IS NULL;

-- ===== CREATOR ATTRIBUTION =====

-- name: GetCreatorOrderAttribution :many
SELECT 
    creator_id,
    COUNT(*) as orders_from_creator,
    SUM(total) as total_revenue_cents,
    COUNT(DISTINCT user_id) as unique_customers
FROM orders 
WHERE created_at BETWEEN $1 AND $2 AND creator_id IS NOT NULL AND deleted_at IS NULL
GROUP BY creator_id
ORDER BY total_revenue_cents DESC;

-- ===== PRODUCT PERFORMANCE =====

-- name: GetTopProducts :many
SELECT 
    oi.product_id,
    COUNT(*) as units_sold,
    SUM(oi.total_cents) as revenue_cents,
    AVG(oi.unit_price_cents) as avg_price_cents
FROM order_items oi
JOIN orders o ON oi.order_id = o.id
WHERE o.created_at BETWEEN $1 AND $2 AND o.deleted_at IS NULL AND oi.deleted_at IS NULL
GROUP BY oi.product_id
ORDER BY units_sold DESC
LIMIT $3;

-- name: GetProductMetrics :one
SELECT 
    product_id,
    COUNT(*) as total_units_sold,
    SUM(total_cents) as total_revenue_cents,
    AVG(unit_price_cents) as avg_price_cents,
    MIN(unit_price_cents) as min_price_cents,
    MAX(unit_price_cents) as max_price_cents
FROM order_items 
WHERE product_id = $1 AND deleted_at IS NULL
GROUP BY product_id;

-- ===== CUSTOMER ANALYTICS =====

-- name: GetCustomerRepeatPurchaseRate :one
SELECT 
    COUNT(DISTINCT user_id) as unique_customers,
    COUNT(DISTINCT CASE WHEN (SELECT COUNT(*) FROM orders WHERE user_id = o.user_id AND deleted_at IS NULL) > 1 THEN o.user_id END) as repeat_customers,
    ROUND(100.0 * COUNT(DISTINCT CASE WHEN (SELECT COUNT(*) FROM orders WHERE user_id = o.user_id AND deleted_at IS NULL) > 1 THEN o.user_id END) / NULLIF(COUNT(DISTINCT user_id), 0), 2) as repeat_rate
FROM orders o
WHERE o.created_at BETWEEN $1 AND $2 AND o.deleted_at IS NULL;

-- name: GetCustomerLTV :one
SELECT 
    user_id,
    COUNT(*) as total_orders,
    SUM(total) as lifetime_value_cents,
    AVG(total) as avg_order_value_cents,
    MAX(created_at) as last_order_date
FROM orders 
WHERE user_id = $1 AND deleted_at IS NULL
GROUP BY user_id;

-- ===== PAYMENT ANALYTICS =====

-- name: GetPaymentMethodDistribution :many
SELECT 
    payment_method,
    COUNT(*) as usage_count,
    SUM(total) as total_amount_cents,
    ROUND(100.0 * COUNT(*) / SUM(COUNT(*)) OVER (), 2) as percentage
FROM orders 
WHERE created_at BETWEEN $1 AND $2 AND payment_method IS NOT NULL AND deleted_at IS NULL
GROUP BY payment_method
ORDER BY usage_count DESC;

-- name: GetFailedPaymentStats :one
SELECT 
    COUNT(*) as total_failed_payments,
    COUNT(DISTINCT user_id) as users_with_failed_payments,
    COUNT(DISTINCT order_id) as orders_with_failed_payments
FROM payment_intents 
WHERE status = 'failed' AND created_at BETWEEN $1 AND $2;
