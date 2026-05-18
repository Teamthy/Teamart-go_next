-- Order Items Queries

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING id, order_id, product_id, quantity, price, created_at;

-- name: GetOrderItemByID :one
SELECT id, order_id, product_id, quantity, price, created_at
FROM order_items
WHERE id = $1;

-- name: ListOrderItemsByOrderID :many
SELECT id, order_id, product_id, quantity, price, created_at
FROM order_items
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: GetOrderItemsByProductID :many
SELECT id, order_id, product_id, quantity, price, created_at
FROM order_items
WHERE product_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderItemQuantity :one
UPDATE order_items
SET quantity = $2
WHERE id = $1
RETURNING id, order_id, product_id, quantity, price, created_at;

-- name: UpdateOrderItemPrice :one
UPDATE order_items
SET price = $2
WHERE id = $1
RETURNING id, order_id, product_id, quantity, price, created_at;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

-- name: DeleteOrderItems :exec
DELETE FROM order_items
WHERE order_id = $1;

-- name: CountOrderItems :one
SELECT COUNT(*) FROM order_items WHERE order_id = $1;

-- name: GetOrderItemsWithProduct :many
SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at,
       p.sku, p.name, p.description
FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.order_id = $1
ORDER BY oi.created_at DESC;
