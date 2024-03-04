-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, first_name, last_name, email, password, phone, email_verified_at)
VALUES ($1, NOW(), NOW(), $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;