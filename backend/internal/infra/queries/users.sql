-- Users Queries

-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, name, password_hash, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET name = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, email, name, password_hash, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
