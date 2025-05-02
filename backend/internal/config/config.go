package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	URI      string
	Database string
	Client   *mongo.Client
}

func (c *MongoDBConfig) GetCollection(name string) *mongo.Collection {
	return c.Client.Database(c.Database).Collection(name)
}

func (c *MongoDBConfig) Close() error {
	return c.Client.Disconnect(context.Background())
}

type JWTConfig struct {
	SecretKey string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

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

func NewJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey: secretKey,
	}
}

// ... existing code ...
