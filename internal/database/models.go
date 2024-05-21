// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description sql.NullString
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PasswordResetToken struct {
	Email     string
	Token     string
	CreatedAt time.Time
}

type SubCategory struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description sql.NullString
	IsActive    bool
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Phone           string
	Email           string
	Password        string
	EmailVerifiedAt sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
