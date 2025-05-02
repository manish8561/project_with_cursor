package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	JWTSecret string
	MongoURI  string
	MongoDB   string
}

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

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
