package handlers

import (
	"database/sql"
	"errors"
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/internal/database"
	"github.com/a4anthony/go-store/models"
	"github.com/a4anthony/go-store/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type RegisterParams struct {
	FirstName string `validate:"required,min=2,max=32" json:"first_name"`
	LastName  string `validate:"required,min=2,max=32" json:"last_name"`
	Phone     string `validate:"required,numeric" json:"phone"`
	Email     string `validate:"required,email" json:"email"`
	Password  string `validate:"required,min=8,max=32" json:"password"`
}

type LoginParams struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8,max=32" json:"password"`
}

type UserResponse struct {
	Message     string      `json:"message"`
	User        models.User `json:"user,omitempty"`
	UserJSON    []byte      `json:"user_json,omitempty"`
	AccessToken string      `json:"access_token,omitempty"`
}

var validate = validator.New()

func Register(c *fiber.Ctx) error {
	params := new(RegisterParams)
	err := utils.CheckParams(c, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	validationErr := validate.Struct(*params)
	if validationErr != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.SetError(validationErr))
	}

	checkIfEmailExists, _ := config.DB.GetUserByEmail(c.Context(), params.Email)
	if checkIfEmailExists.Email != "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorMsg{
			Message: "Validation error",
			Errors: map[string]string{
				"email": "The email has already been taken.",
			},
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	user, err := config.DB.CreateUser(c.Context(), database.CreateUserParams{
		ID:              uuid.New(),
		FirstName:       params.FirstName,
		LastName:        params.LastName,
		Email:           params.Email,
		Password:        string(hashedPassword),
		Phone:           params.Phone,
		EmailVerifiedAt: sql.NullTime{},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(UserResponse{
		Message:     "User created successfully",
		AccessToken: token,
		User:        models.DbUserToUser(user),
	})
}

// Login godoc
// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Success 200 {object} UserResponse
// @Failure 400 {object} utils.ErrorMsg
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	params := new(LoginParams)
	err := utils.CheckParams(c, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validationErr := validate.Struct(*params)
	if validationErr != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.SetError(validationErr))
	}

	user, err := config.DB.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		log.Printf("Error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	err = utils.VerifyPassword(params.Password, user.Password)
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(
		UserResponse{
			User:        models.DbUserToUser(user),
			AccessToken: token,
			Message:     "User logged in successfully.",
		},
	)
}

func Me(c *fiber.Ctx) error {
	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	user, err := config.DB.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(UserResponse{
		User:    models.DbUserToUser(user),
		Message: "User retrieved successfully.",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	userID, err := utils.ExtractTokenID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}
	err = config.DB.DeleteUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully.",
	})
}
