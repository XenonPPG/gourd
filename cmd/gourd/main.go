package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gourd/internal/counter"
	"gourd/internal/document"
	"gourd/internal/gateway"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

const (
	listenAddr      = ":3000"
	shutdownTimeout = 10 * time.Second
)

// @title Gourd API
// @version 0.1
// @description API Gateway for the Gourd

// @BasePath /
// @schemes http https

func main() {
	if err := run(); err != nil {
		log.Fatalf("gourd: %v", err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	counterService, err := counter.New()
	if err != nil {
		return fmt.Errorf("init counter: %w", err)
	}
	log.Println("counter initialized")

	documentService := document.New(counterService)
	gatewayService := gateway.New(documentService, counterService)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	app.Get("/health", gatewayService.Health)
	app.Get("/rules/:type", gatewayService.ListRules)
	app.Get("/counter", gatewayService.GetCounterValue)

	docs := app.Group("/documents")
	docs.Post("/file", gatewayService.ProcessFile)
	docs.Post("/raw", gatewayService.ProcessRaw)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", listenAddr)
		serverErr <- app.Listen(listenAddr)
	}()

	var runErr error

	select {
	case err := <-serverErr:
		if err != nil {
			log.Printf("server stopped unexpectedly: %v", err)
			runErr = fmt.Errorf("server: %w", err)
		}

	case <-ctx.Done():
		stop()
		log.Printf("shutdown signal received, waiting up to %s for active requests", shutdownTimeout)

		if err := app.ShutdownWithTimeout(shutdownTimeout); err != nil {
			log.Printf("server shutdown error: %v", err)
			runErr = fmt.Errorf("server shutdown: %w", err)
		} else {
			log.Println("server stopped gracefully")
		}
	}

	log.Println("saving counter...")
	if err := counterService.Save(); err != nil {
		log.Printf("failed to save counter: %v", err)
		return errors.Join(runErr, fmt.Errorf("save counter: %w", err))
	}
	log.Println("counter saved")

	return runErr
}
