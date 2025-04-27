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
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
