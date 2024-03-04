package utils

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

type ErrorMsgNew struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func getErrorMsg(fe validator.FieldError, fieldName string) string {
	switch fe.Tag() {
	case "required":
		return "The " + fieldName + " field is required."
	case "email":
		return "The " + fieldName + " field is not valid."
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "min":
		return "The " + fieldName + " field should be greater than " + fe.Param() + " characters."
	case "max":
		return "The " + fieldName + " field should be less than " + fe.Param() + " characters."
	case "password_confirmation":
		return "The " + fieldName + " field confirmation does not match."
	case "datetime_future":
		return "The " + fieldName + " field must be date in the future."
	case "datetime":
		return "The " + fieldName + " field must be a valid datetime. (Format: 2006-01-02T15:04:05)"
	case "username":
		return "The " + fieldName + " field must only contain letters, numbers and underscores."
	case "unique":
		return "The " + fieldName + " field already exists."
	}
	return "Unknown error " + fe.Tag()
}

func SetError(err error) ErrorMsg {
	errors := make(map[string]string)
	var data ErrorMsg
	data.Message = "Validation error"

	for _, err := range err.(validator.ValidationErrors) {
		msg := getErrorMsg(err, ToSnakeCase(err.Field()))
		errors[ToSnakeCase(err.Field())] = msg
	}
	data.Errors = errors
	return data
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func CheckParams[T any](c *fiber.Ctx, params T) error {
	if err := c.BodyParser(params); err != nil {
		if err.Error() != "Unprocessable Entity" {
			return err
		}
	}
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
