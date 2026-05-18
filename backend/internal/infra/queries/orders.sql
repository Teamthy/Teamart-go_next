-- Orders Queries

-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, status)
VALUES ($1, $2, $3)
RETURNING id, user_id, total_amount, status, created_at, updated_at;

-- name: GetOrderByID :one
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: ListOrdersByUserID :many
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListOrdersByStatus :many
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAllOrders :many
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, total_amount, status, created_at, updated_at;

-- name: UpdateOrderTotalAmount :one
UPDATE orders
SET total_amount = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, user_id, total_amount, status, created_at, updated_at;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: CountOrdersByUserID :one
SELECT COUNT(*) FROM orders WHERE user_id = $1;

-- name: CountOrdersByStatus :one
SELECT COUNT(*) FROM orders WHERE status = $1;
