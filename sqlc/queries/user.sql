-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: ListActiveUsers :many
SELECT * FROM users WHERE is_active = TRUE ORDER BY id;

-- name: SearchUsersByName :many
SELECT * FROM users
WHERE CONCAT(first_name, ' ', last_name) ILIKE '%' || @name::text || '%'
ORDER BY id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
    email,
    password_hash,
    first_name,
    last_name,
    address,
    phone,
    is_active,
    is_superuser
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    first_name = $1,
    last_name = $2,
    email = $3,
    address = $4,
    phone = $5
WHERE id = $6
RETURNING *;

-- name: PromoteUser :exec
UPDATE users SET
    is_superuser = TRUE
WHERE id = $1;

-- name: DemoteUser :exec
UPDATE users SET
    is_superuser = FALSE
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;


-- name: ActivateUser :exec
UPDATE users SET
    is_active = TRUE
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users SET
    is_active = FALSE
WHERE id = $1;