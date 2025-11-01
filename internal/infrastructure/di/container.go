package di

import (
	"citary-backend/internal/domain/usecases/auth"
	"citary-backend/internal/infrastructure/config"
	httpServer "citary-backend/internal/infrastructure/http"
	authHandler "citary-backend/internal/infrastructure/http/handlers/auth"
	"citary-backend/internal/infrastructure/http/router"
	"citary-backend/internal/infrastructure/persistence/postgres"
	"citary-backend/internal/infrastructure/persistence/postgres/repositories"
	"context"
	"log"
	"time"
)

// Container holds all application dependencies
type Container struct {
	Server *httpServer.Server
	dbConn *postgres.Connection
}

// NewContainer creates and initializes the dependency injection container
func NewContainer() *Container {
	// Load configuration
	config.Load()
	cfg := config.AppConfig

	// Initialize database connection
	dbConn, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Initialize repositories
	userRepository := repositories.NewUserRepositoryImpl(dbConn.DB)

	// Initialize use cases
	signupUserUseCase := auth.NewSignupUserUseCase(userRepository)

	// Initialize HTTP handlers
	authHandlerInstance := authHandler.NewAuthHandler(signupUserUseCase)

	// Initialize router
	routerInstance := router.NewRouter(authHandlerInstance)

	// Initialize HTTP server
	server := httpServer.NewServer(cfg.Port, routerInstance.SetupRoutes())

	return &Container{
		Server: server,
		dbConn: dbConn,
	}
}

// Cleanup closes all resources and connections
func (c *Container) Cleanup() {
	log.Println("Closing connections...")
	if err := c.dbConn.Close(); err != nil {
		log.Printf("Error closing PostgreSQL connection: %v", err)
	}
}

// Shutdown performs graceful shutdown of the server
func (c *Container) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Server.Shutdown(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	c.Cleanup()
}
