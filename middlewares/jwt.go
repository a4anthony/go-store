package middlewares

import (
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/utils"
	"github.com/gofiber/fiber/v2"
	"log"
)

func JwtAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("log middleware")
		err := utils.TokenValid(c)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		userID, err := utils.ExtractTokenID(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		_, err = config.DB.GetUser(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		return c.Next()
	}
}
