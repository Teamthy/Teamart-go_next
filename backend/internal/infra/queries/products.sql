-- Products Queries

-- name: CreateProduct :one
INSERT INTO products (sku, name, description, price)
VALUES ($1, $2, $3, $4)
RETURNING id, sku, name, description, price, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, sku, name, description, price, created_at, updated_at
FROM products
WHERE id = $1;

-- name: GetProductBySKU :one
SELECT id, sku, name, description, price, created_at, updated_at
FROM products
WHERE sku = $1;

-- name: ListProducts :many
SELECT id, sku, name, description, price, created_at, updated_at
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchProducts :many
SELECT id, sku, name, description, price, created_at, updated_at
FROM products
WHERE name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, sku, name, description, price, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: CountProducts :one
SELECT COUNT(*) FROM products;
