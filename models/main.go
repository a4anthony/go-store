package models

import (
	"database/sql"
	"github.com/a4anthony/go-store/internal/database"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID              uuid.UUID  `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Password        string     `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func DbUserToUser(dbUser database.User) User {
	return User{
		ID:              dbUser.ID,
		FirstName:       dbUser.FirstName,
		LastName:        dbUser.LastName,
		Phone:           dbUser.Phone,
		Email:           dbUser.Email,
		Password:        dbUser.Password,
		EmailVerifiedAt: nullTimeToTimePtr(dbUser.EmailVerifiedAt),
		CreatedAt:       dbUser.CreatedAt,
		UpdatedAt:       dbUser.UpdatedAt,
	}

}

func nullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
