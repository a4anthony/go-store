package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(firstName string, lastName string, phone string, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = config.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:              uuid.New(),
		FirstName:       firstName,
		LastName:        lastName,
		Phone:           phone,
		Email:           email,
		Password:        string(hashedPassword),
		EmailVerifiedAt: sql.NullTime{},
	})
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
