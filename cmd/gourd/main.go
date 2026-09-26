package main

import (
	"gourd/internal/gateway"

	"log"

	"github.com/gofiber/fiber/v3"
)

// @title Gourd API
// @version 0.1
// @description API Gateway for the Gourd

// @BasePath /
// @schemes http https

func main() {
	app := fiber.New()

	app.Get("/health", gateway.Health)
	app.Get("/rules", gateway.ListRules)

	docs := app.Group("/documents")
	docs.Post("/file", gateway.ProcessFile)
	docs.Post("/raw", gateway.ProcessRaw)

	log.Fatal(app.Listen(":3000"))
}
