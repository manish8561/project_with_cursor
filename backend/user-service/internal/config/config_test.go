package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigHelpers(t *testing.T) {
	assert.Equal(t, "fallback", getEnv("USER_SERVICE_TEST_UNSET", "fallback"))
	t.Setenv("USER_SERVICE_TEST_VALUE", "configured")
	assert.Equal(t, "configured", getEnv("USER_SERVICE_TEST_VALUE", "fallback"))
	assert.Equal(t, []string{"one", "two"}, splitCSV(" one, ,two "))
}

func TestLoadConfigUsesDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("MONGO_URI", "")
	t.Setenv("MONGO_DB", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("KAFKA_BROKERS", "")
	t.Setenv("KAFKA_CLIENT_ID", "")
	t.Setenv("KAFKA_GROUP_ID", "")
	t.Setenv("KAFKA_TOPIC_USER_CREATED", "")
	t.Setenv("KAFKA_TOPIC_USER_UPDATED", "")
	t.Setenv("KAFKA_TOPIC_USER_DELETED", "")
	t.Setenv("ALLOWED_ORIGINS", "")

	cfg := LoadConfig()

	assert.Equal(t, "8082", cfg.Port)
	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoURI)
	assert.Equal(t, "user_db", cfg.MongoDB)
	assert.Equal(t, "your-secret-key", cfg.JWTSecret)
	assert.Equal(t, "", cfg.KafkaBrokers)
	assert.Equal(t, "user-service", cfg.KafkaClientID)
	assert.Equal(t, "user-service-group", cfg.KafkaGroupID)
	assert.Equal(t, "user.created.v1", cfg.KafkaTopicUserCreated)
	assert.Equal(t, "user.updated.v1", cfg.KafkaTopicUserUpdated)
	assert.Equal(t, "user.deleted.v1", cfg.KafkaTopicUserDeleted)
	assert.Equal(t, []string{"http://localhost:4200", "http://localhost:8085"}, cfg.AllowedOrigins)
}

func TestLoadConfigReadsEnvironmentOverrides(t *testing.T) {
	t.Setenv("PORT", "9092")
	t.Setenv("MONGO_URI", "mongodb://mongo.example:27017")
	t.Setenv("MONGO_DB", "user_prod")
	t.Setenv("JWT_SECRET", "user-secret")
	t.Setenv("KAFKA_BROKERS", "broker-1:9092,broker-2:9092")
	t.Setenv("KAFKA_CLIENT_ID", "user-worker")
	t.Setenv("KAFKA_GROUP_ID", "user-worker-group")
	t.Setenv("KAFKA_TOPIC_USER_CREATED", "custom.user.created")
	t.Setenv("KAFKA_TOPIC_USER_UPDATED", "custom.user.updated")
	t.Setenv("KAFKA_TOPIC_USER_DELETED", "custom.user.deleted")
	t.Setenv("ALLOWED_ORIGINS", " https://app.example.com, , https://admin.example.com ")

	cfg := LoadConfig()

	assert.Equal(t, "9092", cfg.Port)
	assert.Equal(t, "mongodb://mongo.example:27017", cfg.MongoURI)
	assert.Equal(t, "user_prod", cfg.MongoDB)
	assert.Equal(t, "user-secret", cfg.JWTSecret)
	assert.Equal(t, "broker-1:9092,broker-2:9092", cfg.KafkaBrokers)
	assert.Equal(t, "user-worker", cfg.KafkaClientID)
	assert.Equal(t, "user-worker-group", cfg.KafkaGroupID)
	assert.Equal(t, "custom.user.created", cfg.KafkaTopicUserCreated)
	assert.Equal(t, "custom.user.updated", cfg.KafkaTopicUserUpdated)
	assert.Equal(t, "custom.user.deleted", cfg.KafkaTopicUserDeleted)
	assert.Equal(t, []string{"https://app.example.com", "https://admin.example.com"}, cfg.AllowedOrigins)
}
