package services

import (
	"backend/internal/models"
	"errors"
)

type UserService struct {
	// In a real application, you would have a database connection here
	users []models.User
}

func NewUserService() *UserService {
	return &UserService{
		users: []models.User{
			{
				ID:       1,
				Email:    "test@example.com",
				Password: "password123", // In real app, this would be hashed
			},
		},
	}
}

func (s *UserService) Login(email, password string) (*models.LoginResponse, error) {
	for _, user := range s.users {
		if user.Email == email && user.Password == password {
			// In a real application, you would generate a JWT token here
			return &models.LoginResponse{
				Token: "dummy-token", // Replace with JWT token in production
				User:  user,
			}, nil
		}
	}
	return nil, errors.New("invalid credentials")
}
