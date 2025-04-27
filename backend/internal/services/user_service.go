package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	mongoConfig *config.MongoDBConfig
}

func NewUserService(mongoConfig *config.MongoDBConfig) *UserService {
	return &UserService{
		mongoConfig: mongoConfig,
	}
}

// GetMongoConfig returns the MongoDB configuration
func (s *UserService) GetMongoConfig() *config.MongoDBConfig {
	return s.mongoConfig
}

func (s *UserService) Login(email, password string) (*models.LoginResponse, error) {
	collection := s.mongoConfig.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// In a real application, you would generate a JWT token here
	return &models.LoginResponse{
		Token: "dummy-token", // Replace with JWT token in production
		User:  user,
	}, nil
}

func (s *UserService) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
	collection := s.mongoConfig.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	count, err := collection.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	// Create new user
	newUser := models.User{
		ID:       primitive.NewObjectID().Hex(),
		Email:    req.Email,
		Password: req.Password, // In real app, this would be hashed
	}

	// Insert user into database
	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &models.RegisterResponse{
		Message: "User registered successfully",
		User:    newUser,
	}, nil
}
