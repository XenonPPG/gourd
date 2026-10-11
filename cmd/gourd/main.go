package main

import (
	"context"
	"errors"
	"fmt"
	"gourd/internal/config"
	"gourd/internal/counter"
	"gourd/internal/document"
	"gourd/internal/gateway"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// config
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		return fmt.Errorf("init config: %w", err)
	}

	// counter service
	counterService, err := counter.New()
	if err != nil {
		return fmt.Errorf("init counter: %w", err)
	}
	log.Println("counter initialized")

	// services
	documentService := document.New(cfg.DocumentService.WorkerCount, counterService)
	gatewayService := gateway.New(documentService, counterService)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", cfg.Server.Address)
		serverErr <- gatewayService.Start(cfg.Server.Address)
	}()

	var runErr error

	// graceful shutdown
	select {
	case err := <-serverErr:
		if err != nil {
			log.Printf("server stopped unexpectedly: %v", err)
			runErr = fmt.Errorf("server: %w", err)
		}

	case <-ctx.Done():
		stop()
		log.Printf("shutdown signal received, waiting up to %s for active requests", cfg.Server.ShutdownTimeout)

		if err := gatewayService.Shutdown(cfg.Server.ShutdownTimeout); err != nil {
			log.Printf("server shutdown error: %v", err)
			runErr = fmt.Errorf("server shutdown: %w", err)
		} else {
			log.Println("server stopped gracefully")
		}
	}

	log.Println("stopping document service...")
	if err := documentService.Shutdown(cfg.Server.ShutdownTimeout); err != nil {
		log.Println("failed to stop document service: ", err)
	}

	log.Println("saving counter...")
	if err := counterService.Save(); err != nil {
		log.Printf("failed to save counter: %v", err)
		return errors.Join(runErr, fmt.Errorf("save counter: %w", err))
	}
	log.Println("counter saved")

	return runErr
}
