package models

type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Email    string `json:"email" binding:"required,email" bson:"email"`
	Password string `json:"password" binding:"required,min=6" bson:"password"`
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
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
