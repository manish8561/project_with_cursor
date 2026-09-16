package services

import "user-service/internal/models"

// UserServiceInterface defines the contract for user service operations
type UserServiceInterface interface {
	GetUserByID(id string) (*models.User, error)
	ListUsers(page, pageSize int) (*models.UserListResponse, error)
	UpdateUser(id string, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id string) error
	CreateUser(req models.CreateUserRequest) (*models.User, error)
	UpsertUserProfileFromEvent(event models.UserEvent) error
	DeleteUserProfileFromEvent(event models.UserEvent) error
}
