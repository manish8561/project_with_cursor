package handlers

import (
	"backend/internal/services"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/config"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestHandler(t *testing.T) (*UserHandler, func()) {
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

	// Create service
	service := services.NewUserService(&config.MongoDBConfig{
		Client:   client,
		Database: "testdb",
	})

	// Create handler
	handler := NewUserHandler(service)

	// Cleanup function
	cleanup := func() {
		client.Database("testdb").Drop(context.Background())
		client.Disconnect(context.Background())
	}

	return handler, cleanup
}

func TestUserHandler_Register(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	tests := []struct {
		name       string
		req        models.RegisterRequest
		wantStatus int
	}{
		{
			name: "Valid registration",
			req: models.RegisterRequest{
				Email:           "test@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid email format",
			req: models.RegisterRequest{
				Email:           "invalid-email",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Password mismatch",
			req: models.RegisterRequest{
				Email:           "test2@example.com",
				Password:        "password123",
				ConfirmPassword: "differentpass",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(rr)
			c.Request = req

			// Call handler
			handler.Register(c)

			// Check status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create test user directly in the database
	collection := handler.userService.GetMongoConfig().GetCollection("users")
	testUser := models.User{
		ID:       primitive.NewObjectID().Hex(),
		Email:    "test@example.com",
		Password: "password123", // In real app, this would be hashed
	}
	_, err := collection.InsertOne(context.Background(), testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Verify the user was inserted
	var insertedUser models.User
	err = collection.FindOne(context.Background(), bson.M{"email": "test@example.com"}).Decode(&insertedUser)
	if err != nil {
		t.Fatalf("Failed to find inserted user: %v", err)
	}

	tests := []struct {
		name       string
		req        models.LoginRequest
		wantStatus int
	}{
		{
			name: "Valid login",
			req: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid password",
			req: models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Non-existent user",
			req: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(rr)
			c.Request = req

			// Call handler
			handler.Login(c)

			// Check status code
			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}
