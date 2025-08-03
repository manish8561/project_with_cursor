package services

import (
	"user-service/internal/config"
	"user-service/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService handles user-related business logic
type UserService struct {
	mongoConfig *config.MongoDBConfig
}

// NewUserService creates a new UserService with the provided MongoDB configuration
func NewUserService(mongoConfig *config.MongoDBConfig) *UserService {
	return &UserService{
		mongoConfig: mongoConfig,
	}
}

// GetUserByID retrieves a user by their unique identifier
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	collection := s.mongoConfig.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// ListUsers returns a paginated list of users and the total count
func (s *UserService) ListUsers(page, pageSize int) (*models.UserListResponse, error) {
	collection := s.mongoConfig.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)

	filter := bson.M{"role": "customer"}
	findOptions := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.UserListResponse{
		Users: users,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest) (*models.User, error) {
	collection := s.mongoConfig.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":       req.Name,
			"email":      req.Email,
			"role":       req.Role,
			"updatedAt":  time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("user not found")
	}

	// Return updated user
	return s.GetUserByID(id)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	collection := s.mongoConfig.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
} 