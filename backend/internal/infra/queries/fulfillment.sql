-- queries/fulfillment.sql
-- Fulfillment and inventory management queries

-- ===== FULFILLMENT BATCHES =====

-- name: CreateFulfillmentBatch :one
INSERT INTO fulfillment_batches (
    merchant_id, batch_number, status, orders_count, items_count, notes
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetFulfillmentBatch :one
SELECT * FROM fulfillment_batches WHERE id = $1;

-- name: GetBatchByNumber :one
SELECT * FROM fulfillment_batches WHERE batch_number = $1;

-- name: ListMerchantBatches :many
SELECT * FROM fulfillment_batches 
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListBatchesByStatus :many
SELECT * FROM fulfillment_batches 
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActiveBatches :many
SELECT * FROM fulfillment_batches 
WHERE status IN ('pending', 'picking', 'packing', 'ready_to_ship')
ORDER BY created_at ASC;

-- name: UpdateBatchStatus :one
UPDATE fulfillment_batches 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: StartBatch :one
UPDATE fulfillment_batches 
SET status = 'picking', started_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CompleteBatch :one
UPDATE fulfillment_batches 
SET status = 'completed', completed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== INVENTORY RESERVATIONS =====

-- name: CreateInventoryReservation :one
INSERT INTO inventory_reservations (
    order_id, product_id, variant_id, quantity_reserved, expires_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetInventoryReservation :one
SELECT * FROM inventory_reservations WHERE id = $1;

-- name: ListOrderReservations :many
SELECT * FROM inventory_reservations 
WHERE order_id = $1
ORDER BY created_at DESC;

-- name: ListProductReservations :many
SELECT * FROM inventory_reservations 
WHERE product_id = $1 AND status IN ('reserved', 'partially_released')
ORDER BY created_at DESC;

-- name: ListActiveReservations :many
SELECT * FROM inventory_reservations 
WHERE status = 'reserved' AND expires_at > CURRENT_TIMESTAMP
ORDER BY expires_at ASC;

-- name: ListExpiredReservations :many
SELECT * FROM inventory_reservations 
WHERE status = 'reserved' AND expires_at <= CURRENT_TIMESTAMP
ORDER BY expires_at ASC;

-- name: UpdateReservationStatus :one
UPDATE inventory_reservations 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ReleaseReservation :one
UPDATE inventory_reservations 
SET status = 'fully_released', quantity_released = quantity_reserved,
    released_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: PartiallyReleaseReservation :one
UPDATE inventory_reservations 
SET status = 'partially_released', quantity_released = $2,
    released_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== ANALYTICS =====

-- name: GetInventoryStatus :one
SELECT 
    COUNT(*) as total_reservations,
    COUNT(CASE WHEN status = 'reserved' THEN 1 END) as active_reservations,
    COUNT(CASE WHEN status = 'fully_released' THEN 1 END) as released_reservations,
    SUM(quantity_reserved) as total_reserved_quantity,
    SUM(quantity_released) as total_released_quantity
FROM inventory_reservations;

-- name: GetReservationsByProduct :many
SELECT 
    product_id,
    COUNT(*) as total_reservations,
    SUM(quantity_reserved) as total_reserved_quantity,
    SUM(CASE WHEN status = 'reserved' THEN quantity_reserved ELSE 0 END) as currently_reserved
FROM inventory_reservations 
WHERE product_id = $1
GROUP BY product_id;
