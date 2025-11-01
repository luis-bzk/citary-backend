package main

import (
	"citary-backend/internal/infrastructure/di"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize dependency injection container
	container := di.NewContainer()
	defer container.Cleanup()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	go func() {
		if err := container.Server.Start(); err != nil {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("\n👋 Shutting down gracefully...")

	// Perform graceful shutdown
	container.Shutdown()

	log.Println("Server stopped")
}
