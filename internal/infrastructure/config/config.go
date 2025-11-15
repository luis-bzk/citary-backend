package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	// Server configuration
	Port int

	// Database configuration
	DatabaseURL string

	// SMTP configuration
	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string

	// Frontend configuration
	FrontendURL string
}

// AppConfig is the global configuration instance
var AppConfig *Config

// Load loads application configuration from environment variables
func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Required variables
	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	smtpHost := getEnv("SMTP_HOST", "")
	if smtpHost == "" {
		log.Fatal("SMTP_HOST environment variable is required")
	}

	smtpPort := getEnv("SMTP_PORT", "")
	if smtpPort == "" {
		log.Fatal("SMTP_PORT environment variable is required")
	}

	smtpUsername := getEnv("SMTP_USERNAME", "")
	if smtpUsername == "" {
		log.Fatal("SMTP_USERNAME environment variable is required")
	}

	smtpPassword := getEnv("SMTP_PASSWORD", "")
	if smtpPassword == "" {
		log.Fatal("SMTP_PASSWORD environment variable is required")
	}

	smtpFromEmail := getEnv("SMTP_FROM_EMAIL", "")
	if smtpFromEmail == "" {
		log.Fatal("SMTP_FROM_EMAIL environment variable is required")
	}

	// Optional variables with defaults
	port := getEnvAsInt("PORT", 3001)
	smtpFromName := getEnv("SMTP_FROM_NAME", "Citary")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")

	AppConfig = &Config{
		Port:          port,
		DatabaseURL:   databaseURL,
		SMTPHost:      smtpHost,
		SMTPPort:      smtpPort,
		SMTPUsername:  smtpUsername,
		SMTPPassword:  smtpPassword,
		SMTPFromEmail: smtpFromEmail,
		SMTPFromName:  smtpFromName,
		FrontendURL:   frontendURL,
	}

	log.Printf("Configuration loaded: PORT=%d, SMTP_HOST=%s, FRONTEND_URL=%s",
		AppConfig.Port, AppConfig.SMTPHost, AppConfig.FrontendURL)
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
