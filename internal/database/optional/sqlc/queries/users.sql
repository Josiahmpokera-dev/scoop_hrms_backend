-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListUsersByType :many
SELECT * FROM users
WHERE user_type = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash,
    first_name,
    last_name,
    user_type,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE($2, username),
    email = COALESCE($3, email),
    first_name = COALESCE($4, first_name),
    last_name = COALESCE($5, last_name),
    user_type = COALESCE($6, user_type),
    is_active = COALESCE($7, is_active),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1;

-- name: CheckUserExistsByEmail :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE email = $1 AND deleted_at IS NULL
) AS exists;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE deleted_at IS NULL;
