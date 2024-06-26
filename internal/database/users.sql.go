// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, first_name, last_name, email, password, phone, email_verified_at)
VALUES ($1, NOW(), NOW(), $2, $3, $4, $5, $6, $7) RETURNING id, first_name, last_name, phone, email, password, email_verified_at, created_at, updated_at
`

type CreateUserParams struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	Password        string
	Phone           string
	EmailVerifiedAt sql.NullTime
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
		arg.Phone,
		arg.EmailVerifiedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
		&i.Password,
		&i.EmailVerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, first_name, last_name, phone, email, password, email_verified_at, created_at, updated_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
		&i.Password,
		&i.EmailVerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, last_name, phone, email, password, email_verified_at, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
		&i.Password,
		&i.EmailVerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users SET updated_at = NOW(), first_name = $2, last_name = $3, phone = $4 WHERE id = $1 RETURNING id, first_name, last_name, phone, email, password, email_verified_at, created_at, updated_at
`

type UpdateUserParams struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Phone     string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Phone,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Phone,
		&i.Email,
		&i.Password,
		&i.EmailVerifiedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users SET updated_at = NOW(), password = $2 WHERE id = $1
`

type UpdateUserPasswordParams struct {
	ID       uuid.UUID
	Password string
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.ID, arg.Password)
	return err
}
