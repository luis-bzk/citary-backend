package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	Port        int
}

// AppConfig is the global configuration instance
var AppConfig *Config

// Load loads application configuration from environment variables
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	port := getEnvAsInt("PORT", 3001)

	AppConfig = &Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}

	log.Printf("Configuration loaded: PORT=%d", AppConfig.Port)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
