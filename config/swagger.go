package config

import (
	_ "github.com/a4anthony/go-store/docs" // load swagger docs
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func AddSwaggerRoutes(app *fiber.App) {
	// setup swagger
	app.Get("/swagger/*", swagger.HandlerDefault)
}
