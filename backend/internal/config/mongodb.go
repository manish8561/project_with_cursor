package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	URI      string
	Database string
	Client   *mongo.Client
}

func NewMongoDBConfig() *MongoDBConfig {
	username := getEnvWithDefault("MONGODB_USERNAME", "admin")
	password := getEnvWithDefault("MONGODB_PASSWORD", "password123")
	host := getEnvWithDefault("MONGODB_HOST", "localhost")
	port := getEnvWithDefault("MONGODB_PORT", "27017")
	database := getEnvWithDefault("MONGODB_DATABASE", "appdb")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Println(err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err)
	}

	return &MongoDBConfig{
		URI:      uri,
		Database: database,
		Client:   client,
	}
}

func (c *MongoDBConfig) GetCollection(name string) *mongo.Collection {
	return c.Client.Database(c.Database).Collection(name)
}

func (c *MongoDBConfig) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.Client.Disconnect(ctx); err != nil {
		log.Println(err)
	}
}
