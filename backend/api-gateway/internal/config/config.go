package config

import (
	"os"
)

// Config holds all configuration for the API Gateway
type Config struct {
	Port            string
	AuthServiceURL  string
	UserServiceURL  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		AuthServiceURL:  getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		UserServiceURL:  getEnv("USER_SERVICE_URL", "http://localhost:8082"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 