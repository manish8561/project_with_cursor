package models

import "time"

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

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UserListResponse is the response for paginated user list
type UserListResponse struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
}
