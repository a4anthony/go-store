package router

import (
	"github.com/a4anthony/go-store/handlers"
	"github.com/a4anthony/go-store/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	v1ApiRoutes := app.Group("/v1")
	apiRoutes := v1ApiRoutes.Group("/api")

	app.Get("/", handlers.HandleHealthCheck)
	app.Get("/health", handlers.HandleHealthCheck)

	usersGroup := apiRoutes.Group("/users")
	usersGroup.Post("/register", handlers.Register)
	usersGroup.Post("/login", handlers.Login)
	usersGroup.Get("/me", middlewares.JwtAuthMiddleware(), handlers.Me)
	usersGroup.Delete("", middlewares.JwtAuthMiddleware(), handlers.DeleteUser)
}
