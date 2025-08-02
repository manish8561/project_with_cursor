package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" binding:"required,email" bson:"email"`
	Password  string    `json:"password" binding:"required,min=6" bson:"password"`
	Status    string    `json:"status" bson:"status"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

// RegisterResponse represents a registration response
type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// TokenValidationRequest represents a token validation request
type TokenValidationRequest struct {
	Token string `json:"token" binding:"required"`
}

// TokenValidationResponse represents a token validation response
type TokenValidationResponse struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id,omitempty"`
	Message string `json:"message,omitempty"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshTokenResponse represents a token refresh response
type RefreshTokenResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}
