package main

import (
	"fmt"
	"github.com/a4anthony/go-store/config"
	"github.com/a4anthony/go-store/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"os"
)

func main() {
	err := config.LoadENV()
	if err != nil {
		log.Fatalf("Error loading environment variables")
	} else {
		fmt.Println("Environment variables loaded!")
	}

	config.ConnectDb()

	// create app
	app := fiber.New()

	// attach middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// setup routes
	router.SetupRoutes(app)

	port := os.Getenv("PORT")
	err = app.Listen(":" + port)
	if err != nil {
		panic(err)
	}
}
