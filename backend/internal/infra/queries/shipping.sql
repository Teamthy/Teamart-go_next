-- queries/shipping.sql
-- Shipping and tracking queries

-- ===== SHIPPING PROFILES =====

-- name: CreateShippingProfile :one
INSERT INTO shipping_profiles (
    merchant_id, name, description, is_default, zones, rates, carriers,
    default_carrier, free_shipping_threshold_cents
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetShippingProfile :one
SELECT * FROM shipping_profiles WHERE id = $1 AND deleted_at IS NULL;

-- name: GetDefaultShippingProfile :one
SELECT * FROM shipping_profiles 
WHERE merchant_id = $1 AND is_default = TRUE AND is_active = TRUE AND deleted_at IS NULL;

-- name: ListShippingProfiles :many
SELECT * FROM shipping_profiles 
WHERE merchant_id = $1 AND is_active = TRUE AND deleted_at IS NULL
ORDER BY is_default DESC, created_at;

-- name: UpdateShippingProfile :one
UPDATE shipping_profiles 
SET name = $2, zones = $3, rates = $4, carriers = $5, 
    default_carrier = $6, is_default = $7, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SetDefaultShippingProfile :exec
UPDATE shipping_profiles 
SET is_default = FALSE
WHERE merchant_id = $1 AND is_default = TRUE;

-- name: SetDefaultShippingProfileID :one
UPDATE shipping_profiles 
SET is_default = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND merchant_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteShippingProfile :one
UPDATE shipping_profiles 
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== SHIPMENTS =====

-- name: CreateShipment :one
INSERT INTO shipments (
    order_id, shipment_number, carrier, carrier_service, tracking_number,
    shipping_profile_id, weight_grams, length_cm, width_cm, height_cm,
    items_count, shipping_address, shipping_cost_cents, insurance_cost_cents,
    total_cost_cents, expected_delivery_date, notes, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
)
RETURNING *;

-- name: GetShipment :one
SELECT * FROM shipments WHERE id = $1 AND deleted_at IS NULL;

-- name: GetShipmentByTrackingNumber :one
SELECT * FROM shipments WHERE tracking_number = $1 AND deleted_at IS NULL;

-- name: GetShipmentByNumber :one
SELECT * FROM shipments WHERE shipment_number = $1 AND deleted_at IS NULL;

-- name: ListShipmentsByOrder :many
SELECT * FROM shipments 
WHERE order_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListShipmentsByStatus :many
SELECT * FROM shipments 
WHERE status = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListShipmentsByCarrier :many
SELECT * FROM shipments 
WHERE carrier = $1 AND created_at BETWEEN $2 AND $3 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: UpdateShipmentStatus :one
UPDATE shipments 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateShipmentTracking :one
UPDATE shipments 
SET tracking_number = $2, label_url = $3, label_generated_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateShipmentShipped :one
UPDATE shipments 
SET status = 'shipped', shipped_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateShipmentDelivered :one
UPDATE shipments 
SET status = 'delivered', delivered_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- ===== TRACKING EVENTS =====

-- name: CreateTrackingEvent :one
INSERT INTO tracking_events (
    shipment_id, status, location, location_details, event_time,
    description, carrier_code, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetTrackingEvent :one
SELECT * FROM tracking_events WHERE id = $1;

-- name: ListTrackingEvents :many
SELECT * FROM tracking_events 
WHERE shipment_id = $1 
ORDER BY event_time DESC;

-- name: ListLatestTrackingEvent :one
SELECT * FROM tracking_events 
WHERE shipment_id = $1 
ORDER BY event_time DESC
LIMIT 1;

-- name: ListTrackingEventsByStatus :many
SELECT * FROM tracking_events 
WHERE shipment_id = $1 AND status = $2 
ORDER BY event_time DESC;

-- name: ListRecentTrackingEvents :many
SELECT * FROM tracking_events 
WHERE created_at BETWEEN $1 AND $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
