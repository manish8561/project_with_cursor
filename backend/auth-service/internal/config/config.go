package config

import (
	"os"
)

// Config holds all configuration for the auth service
type Config struct {
	Port      string
	MongoURI  string
	MongoDB   string
	JWTSecret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:      getEnv("PORT", "8081"),
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("MONGO_DB", "auth_db"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
