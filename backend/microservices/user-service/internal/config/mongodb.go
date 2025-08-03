package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBConfig handles MongoDB connection and operations
type MongoDBConfig struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDBConfig creates a new MongoDB configuration
func NewMongoDBConfig(uri, dbName string) (*MongoDBConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	return &MongoDBConfig{
		client:   client,
		database: database,
	}, nil
}

// GetCollection returns a MongoDB collection
func (m *MongoDBConfig) GetCollection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Close closes the MongoDB connection
func (m *MongoDBConfig) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
} 