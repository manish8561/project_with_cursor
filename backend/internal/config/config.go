package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig holds configuration and client for MongoDB connection.
type MongoDBConfig struct {
	URI      string        // MongoDB connection URI
	Database string        // MongoDB database name
	Client   *mongo.Client // MongoDB client instance
}

// GetCollection returns a handle to the named collection in the configured database.
func (c *MongoDBConfig) GetCollection(name string) *mongo.Collection {
	return c.Client.Database(c.Database).Collection(name)
}

// Close disconnects the MongoDB client.
func (c *MongoDBConfig) Close() error {
	return c.Client.Disconnect(context.Background())
}

// JWTConfig holds configuration for JWT authentication.
type JWTConfig struct {
	SecretKey string // Secret key for signing JWT tokens
}

// getEnv returns the value of the environment variable for the given key,
// or the provided default value if the environment variable is not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// NewMongoDBConfig creates and returns a new MongoDBConfig, connecting to the specified URI and database.
// It pings the database to verify the connection.
func NewMongoDBConfig(uri, database string) (*MongoDBConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return &MongoDBConfig{
		URI:      uri,
		Database: database,
		Client:   client,
	}, nil
}

// NewJWTConfig creates and returns a new JWTConfig with the provided secret key.
func NewJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey: secretKey,
	}
}
