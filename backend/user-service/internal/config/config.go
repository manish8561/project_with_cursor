package config

import (
	"os"
)

// Config holds all configuration for the user service
type Config struct {
	Port                  string
	MongoURI              string
	MongoDB               string
	KafkaBrokers          string
	KafkaClientID         string
	KafkaGroupID          string
	KafkaTopicUserCreated string
	KafkaTopicUserUpdated string
	KafkaTopicUserDeleted string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:                  getEnv("PORT", "8082"),
		MongoURI:              getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:               getEnv("MONGO_DB", "user_db"),
		KafkaBrokers:          getEnv("KAFKA_BROKERS", ""),
		KafkaClientID:         getEnv("KAFKA_CLIENT_ID", "user-service"),
		KafkaGroupID:          getEnv("KAFKA_GROUP_ID", "user-service-group"),
		KafkaTopicUserCreated: getEnv("KAFKA_TOPIC_USER_CREATED", "user.created.v1"),
		KafkaTopicUserUpdated: getEnv("KAFKA_TOPIC_USER_UPDATED", "user.updated.v1"),
		KafkaTopicUserDeleted: getEnv("KAFKA_TOPIC_USER_DELETED", "user.deleted.v1"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
