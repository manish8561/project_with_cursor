package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Client, func()) {
	// Connect to test MongoDB with authentication
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:password123@localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Clean up any existing data
	err = client.Database("testdb").Drop(context.Background())
	if err != nil {
		t.Fatalf("Failed to clean up test database: %v", err)
	}

	// Cleanup function
	cleanup := func() {
		client.Database("testdb").Drop(context.Background())
		client.Disconnect(context.Background())
	}

	return client, cleanup
}

func TestUserService_Login(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	// Setup test data
	collection := client.Database("testdb").Collection("users")
	testUser := models.User{
		ID:       "1",
		Email:    "test@example.com",
		Password: "password123", // In real app, this would be hashed
	}
	_, err := collection.InsertOne(context.Background(), testUser)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Create service
	service := NewUserService(&config.MongoDBConfig{
		Client:   client,
		Database: "testdb",
	})

	// Test cases
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid credentials",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Invalid password",
			email:    "test@example.com",
			password: "wrongpass",
			wantErr:  true,
		},
		{
			name:     "Non-existent user",
			email:    "nonexistent@example.com",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_Register(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	// Create service
	service := NewUserService(&config.MongoDBConfig{
		Client:   client,
		Database: "testdb",
	})

	// Test cases
	tests := []struct {
		name    string
		req     models.RegisterRequest
		wantErr bool
	}{
		{
			name: "Valid registration",
			req: models.RegisterRequest{
				Email:           "new@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			wantErr: false,
		},
		{
			name: "Duplicate email",
			req: models.RegisterRequest{
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			wantErr: true,
		},
	}

	// Insert a test user for duplicate email test
	collection := client.Database("testdb").Collection("users")
	_, err := collection.InsertOne(context.Background(), models.User{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
