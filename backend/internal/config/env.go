package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds environment configuration for the application.
type Config struct {
	Port      string // Port on which the server runs
	JWTSecret string // Secret key for JWT authentication
	MongoURI  string // MongoDB connection URI
	MongoDB   string // MongoDB database name
}

// LoadConfig loads environment variables from a .env file (if present) and returns a Config struct.
// If the .env file is not found, it uses default values for each configuration field.
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using default values")
	}

	return &Config{
		Port:      getEnvWithDefault("PORT", "8080"),
		JWTSecret: getEnvWithDefault("JWT_SECRET", "your-secret-key"),
		MongoURI:  getEnvWithDefault("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnvWithDefault("MONGO_DB", "testdb"),
	}
}

// getEnvWithDefault returns the value of the environment variable for the given key,
// or the provided default value if the environment variable is not set.
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
