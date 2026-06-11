-- queries/addresses.sql
-- Customer addresses management queries

-- name: CreateCustomerAddress :one
INSERT INTO customer_addresses (
    customer_id, label, recipient_name, phone_number, email,
    street_address, apartment_suite, city, state_province,
    postal_code, country_code, latitude, longitude, is_default, is_validated
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: GetCustomerAddress :one
SELECT * FROM customer_addresses WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCustomerAddresses :many
SELECT * FROM customer_addresses 
WHERE customer_id = $1 AND deleted_at IS NULL
ORDER BY is_default DESC, updated_at DESC;

-- name: GetDefaultCustomerAddress :one
SELECT * FROM customer_addresses 
WHERE customer_id = $1 AND is_default = TRUE AND deleted_at IS NULL;

-- name: GetDefaultShippingAddress :one
SELECT * FROM customer_addresses 
WHERE customer_id = $1 AND is_shipping_address = TRUE AND is_default = TRUE AND deleted_at IS NULL;

-- name: GetDefaultBillingAddress :one
SELECT * FROM customer_addresses 
WHERE customer_id = $1 AND is_billing_address = TRUE AND is_default = TRUE AND deleted_at IS NULL;

-- name: UpdateCustomerAddress :one
UPDATE customer_addresses 
SET recipient_name = $2, phone_number = $3, street_address = $4,
    apartment_suite = $5, city = $6, state_province = $7, postal_code = $8,
    latitude = $9, longitude = $10, is_validated = $11, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND customer_id = $12 AND deleted_at IS NULL
RETURNING *;

-- name: SetDefaultAddress :exec
UPDATE customer_addresses 
SET is_default = FALSE
WHERE customer_id = $1 AND is_default = TRUE;

-- name: SetDefaultAddressID :one
UPDATE customer_addresses 
SET is_default = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND customer_id = $2 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCustomerAddress :one
UPDATE customer_addresses 
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: SearchCustomerAddresses :many
SELECT * FROM customer_addresses 
WHERE customer_id = $1 AND deleted_at IS NULL AND (
    LOWER(recipient_name) ILIKE LOWER($2) OR
    LOWER(street_address) ILIKE LOWER($2) OR
    LOWER(city) ILIKE LOWER($2)
)
ORDER BY is_default DESC;
