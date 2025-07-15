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

	// Create services
	mongoConfig := &config.MongoDBConfig{
		Client:   client,
		Database: "testdb",
	}
	jwtService := NewJWTService(config.NewJWTConfig("test-secret"))
	service := NewUserService(mongoConfig, jwtService)

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

	// Create services
	mongoConfig := &config.MongoDBConfig{
		Client:   client,
		Database: "testdb",
	}
	jwtService := NewJWTService(config.NewJWTConfig("test-secret"))
	service := NewUserService(mongoConfig, jwtService)

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

func TestUserService_ListUsers(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	// Create services
	mongoConfig := &config.MongoDBConfig{
		Client:   client,
		Database: "testdb",
	}
	jwtService := NewJWTService(config.NewJWTConfig("test-secret"))
	service := NewUserService(mongoConfig, jwtService)

	collection := client.Database("testdb").Collection("users")
	// Insert users with different roles
	users := []models.User{
		{ID: "1", Name: "Alice", Email: "alice@example.com", Password: "pass", Role: "customer"},
		{ID: "2", Name: "Bob", Email: "bob@example.com", Password: "pass", Role: "admin"},
		{ID: "3", Name: "Carol", Email: "carol@example.com", Password: "pass", Role: "customer"},
	}
	for _, u := range users {
		_, err := collection.InsertOne(context.Background(), u)
		if err != nil {
			t.Fatalf("Failed to insert user: %v", err)
		}
	}

	// Test pagination and filtering
	result, total, err := service.ListUsers(1, 10)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if total != 2 {
		t.Errorf("Expected total 2 customers, got %d", total)
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 customers in result, got %d", len(result))
	}
	for _, u := range result {
		if u.Role != "customer" {
			t.Errorf("Non-customer user returned: %+v", u)
		}
	}
}
