package services

import (
	"auth-service/internal/config"
	"auth-service/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthService handles authentication-related business logic
type AuthService struct {
	mongoConfig *config.MongoDBConfig
	jwtService  *JWTService
}

// NewAuthService creates a new AuthService with the provided dependencies
func NewAuthService(mongoConfig *config.MongoDBConfig, jwtService *JWTService) *AuthService {
	return &AuthService{
		mongoConfig: mongoConfig,
		jwtService:  jwtService,
	}
}

// Login authenticates a user with the provided email and password.
// Returns a JWT token upon successful authentication or an error if credentials are invalid.
func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	collection := s.mongoConfig.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &models.LoginResponse{
		Status: "success",
		Token:  token,
	}, nil
}

// Register creates a new user account with the provided registration details.
// Returns a success response if registration is successful, or an error if the user already exists.
func (s *AuthService) Register(req models.RegisterRequest) (*models.RegisterResponse, error) {
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
		ID:        primitive.NewObjectID().Hex(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password, // In real app, this would be hashed
		Status:    "active",
		Role:      "customer",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user into database
	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &models.RegisterResponse{
		Status:  "success",
		Message: "User registered successfully",
	}, nil
}

// ValidateToken validates a JWT token and returns user information
func (s *AuthService) ValidateToken(tokenString string) (*models.TokenValidationResponse, error) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return &models.TokenValidationResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	return &models.TokenValidationResponse{
		Valid:  true,
		UserID: claims.UserID,
	}, nil
}

// RefreshToken generates a new token for the user
func (s *AuthService) RefreshToken(tokenString string) (*models.RefreshTokenResponse, error) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	newToken, err := s.jwtService.GenerateToken(claims.UserID)
	if err != nil {
		return nil, errors.New("failed to generate new token")
	}

	return &models.RefreshTokenResponse{
		Status: "success",
		Token:  newToken,
	}, nil
}
