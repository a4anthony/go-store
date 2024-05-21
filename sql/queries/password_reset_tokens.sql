-- create a new password reset token for the given email or update the existing one
-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (email, created_at, token)
VALUES ($1, NOW(), $2)
ON CONFLICT (email) DO UPDATE
SET created_at = NOW(), token = $2;



-- name: GetPasswordResetTokenByEmail :one
SELECT * FROM password_reset_tokens WHERE email = $1;


-- name: DeletePasswordResetToken :exec
DELETE FROM password_reset_tokens WHERE email = $1;
