package handlers

import (
	"errors"
	"fmt"
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/internal/database"
	"github.com/a4anthony/go-store/mailer"
	"github.com/a4anthony/go-store/models"
	"github.com/a4anthony/go-store/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type ProfileParams struct {
	FirstName string `validate:"required,min=2,max=32" json:"first_name"`
	LastName  string `validate:"required,min=2,max=32" json:"last_name"`
	Phone     string `validate:"required,numeric" json:"phone"`
}

type ChangePasswordParams struct {
	OldPassword string `validate:"required,min=8,max=32" json:"old_password"`
	NewPassword string `validate:"required,min=8,max=32" json:"new_password"`
}

type ForgotPasswordParams struct {
	Email string `validate:"required,email" json:"email"`
}

type ResetPasswordParams struct {
	Token    string `validate:"required" json:"token"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,min=8,max=32" json:"password"`
}

func UpdateProfile(c *fiber.Ctx) error {
	params := new(ProfileParams)
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

	userID, _ := utils.ExtractTokenID(c)

	user, err := config.DB.UpdateUser(c.Context(), database.UpdateUserParams{
		ID:        userID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Phone:     params.Phone,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	fmt.Println(user)
	return c.Status(fiber.StatusOK).JSON(UserResponse{
		User:    models.DbUserToUser(user),
		Message: "User updated successfully.",
	})
}

func ChangePassword(c *fiber.Ctx) error {
	params := new(ChangePasswordParams)
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

	userID, _ := utils.ExtractTokenID(c)
	user, _ := config.DB.GetUser(c.Context(), userID)

	err = utils.VerifyPassword(params.OldPassword, user.Password)
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorMsg{
			Message: "Validation error",
			Errors: map[string]string{
				"old_password": "Old password is incorrect.",
			},
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = config.DB.UpdateUserPassword(c.Context(), database.UpdateUserPasswordParams{
		ID:       userID,
		Password: string(hashedPassword),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully.",
	})
}

func ForgotPassword(c *fiber.Ctx) error {
	params := new(ForgotPasswordParams)
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
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorMsg{
			Message: "Validation error",
			Errors: map[string]string{
				"email": "The email does not exist.",
			},
		})
	}

	resetToken := randstr.String(20)
	passwordResetToken := utils.Encode(resetToken)
	appUrl := os.Getenv("APP_URL")
	resetUrl := appUrl + "/reset-password?token=" + passwordResetToken + "&email=" + user.Email
	body := mailer.PrintTemplate(mailer.UserEmail{FirstName: user.FirstName, URL: resetUrl, Subject: "Forgot Password"}, "welcome.html")
	from := os.Getenv("MAIL_FROM_ADDRESS")
	to := user.Email
	subject := "Forgot Password"
	mailer.SendMail(from, to, subject, body, "")

	err = config.DB.CreatePasswordResetToken(c.Context(), database.CreatePasswordResetTokenParams{
		Email: user.Email,
		Token: resetToken,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset link sent to your email.",
	})
}

func ResetPassword(c *fiber.Ctx) error {
	params := new(ResetPasswordParams)
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

	passwordResetToken, err := config.DB.GetPasswordResetTokenByEmail(c.Context(), params.Email)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorMsg{
			Message: "Validation error",
			Errors: map[string]string{
				"token": "The token is invalid.",
			},
		})
	}

	format := "2006-01-02 15:04:05"
	createdAt := passwordResetToken.CreatedAt.Add(60 * time.Minute).Format(format)
	newCreatedAt, _ := time.Parse(format, createdAt)
	now := time.Now().Format(format)
	newNow, _ := time.Parse(format, now)

	if newCreatedAt.Before(newNow) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(utils.ErrorMsg{
			Message: "Validation error",
			Errors: map[string]string{
				"token": "The token has expired.",
			},
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	fmt.Println(string(hashedPassword))

	user, _ := config.DB.GetUserByEmail(c.Context(), params.Email)
	_ = config.DB.UpdateUserPassword(c.Context(), database.UpdateUserPasswordParams{
		ID:       user.ID,
		Password: string(hashedPassword),
	})

	_ = config.DB.DeletePasswordResetToken(c.Context(), params.Email)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password reset successfully.",
	})
}
