package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	db     *mongo.Database
}

func ConnectMongo(uri, dbName string) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, err
	}
	return &Mongo{client: client, db: client.Database(dbName)}, nil
}

func (m *Mongo) Collection(name string) *mongo.Collection { return m.db.Collection(name) }

func (m *Mongo) Database() *mongo.Database { return m.db }

func (m *Mongo) Close(ctx context.Context) error { return m.client.Disconnect(ctx) }
