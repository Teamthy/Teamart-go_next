-- queries/discounts.sql
-- Discount codes and promo queries

-- ===== DISCOUNT CODES =====

-- name: CreateDiscountCode :one
INSERT INTO discount_codes (
    code, campaign_name, discount_type, discount_value, discount_value_cents,
    min_order_value_cents, max_discount_cents, max_uses_total, max_uses_per_customer,
    applicable_to, applicable_items, valid_from, valid_until, is_active, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: GetDiscountCode :one
SELECT * FROM discount_codes WHERE id = $1;

-- name: GetDiscountCodeByCode :one
SELECT * FROM discount_codes WHERE code = $1 AND is_active = TRUE;

-- name: ListActiveDiscountCodes :many
SELECT * FROM discount_codes 
WHERE is_active = TRUE AND valid_from <= CURRENT_TIMESTAMP AND valid_until >= CURRENT_TIMESTAMP
ORDER BY valid_until ASC
LIMIT $1 OFFSET $2;

-- name: ListDiscountCodesByCreator :many
SELECT * FROM discount_codes 
WHERE created_by = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ValidateDiscountCode :one
SELECT * FROM discount_codes 
WHERE code = $1 
AND is_active = TRUE 
AND valid_from <= CURRENT_TIMESTAMP 
AND valid_until >= CURRENT_TIMESTAMP 
AND (max_uses_total IS NULL OR times_used < max_uses_total);

-- name: UpdateDiscountCodeUsage :one
UPDATE discount_codes 
SET times_used = times_used + 1,
    total_discount_given_cents = total_discount_given_cents + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DisableDiscountCode :one
UPDATE discount_codes 
SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== CREATOR PROMOS =====

-- name: CreateCreatorPromo :one
INSERT INTO creator_promos (
    creator_id, name, description, promo_type, discount_percentage,
    discount_amount_cents, product_ids, category_ids, start_at, end_at,
    is_active, livestream_id, livestream_only, min_order_value_cents, max_uses
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: GetCreatorPromo :one
SELECT * FROM creator_promos WHERE id = $1;

-- name: ListCreatorActivePromos :many
SELECT * FROM creator_promos 
WHERE creator_id = $1 AND is_active = TRUE 
AND start_at <= CURRENT_TIMESTAMP AND end_at >= CURRENT_TIMESTAMP
ORDER BY start_at DESC;

-- name: ListCreatorPromos :many
SELECT * FROM creator_promos 
WHERE creator_id = $1
ORDER BY start_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLivestreamPromos :many
SELECT * FROM creator_promos 
WHERE livestream_id = $1 AND is_active = TRUE
ORDER BY start_at DESC;

-- name: ValidateCreatorPromo :one
SELECT * FROM creator_promos 
WHERE id = $1 
AND is_active = TRUE 
AND start_at <= CURRENT_TIMESTAMP 
AND end_at >= CURRENT_TIMESTAMP
AND (max_uses IS NULL OR times_used < max_uses);

-- name: UpdateCreatorPromoUsage :one
UPDATE creator_promos 
SET times_used = times_used + 1, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DisableCreatorPromo :one
UPDATE creator_promos 
SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== ANALYTICS =====

-- name: GetDiscountCodeStats :one
SELECT 
    COUNT(*) as total_codes,
    COUNT(CASE WHEN is_active = TRUE THEN 1 END) as active_codes,
    SUM(times_used) as total_uses,
    SUM(total_discount_given_cents) as total_discount_amount_cents
FROM discount_codes 
WHERE created_by = $1;

-- name: GetPromoStats :one
SELECT 
    COUNT(*) as total_promos,
    COUNT(CASE WHEN is_active = TRUE THEN 1 END) as active_promos,
    SUM(times_used) as total_uses,
    COUNT(CASE WHEN timestamptz(NOW()) BETWEEN start_at AND end_at THEN 1 END) as currently_running
FROM creator_promos 
WHERE creator_id = $1;

-- name: GetMostUsedDiscountCodes :many
SELECT 
    code,
    times_used,
    total_discount_given_cents,
    ROUND(total_discount_given_cents::NUMERIC / NULLIF(times_used, 0), 2) as avg_discount_per_use
FROM discount_codes 
WHERE created_by = $1
ORDER BY times_used DESC
LIMIT $2;
