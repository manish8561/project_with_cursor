package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login handles user authentication and returns a JWT token
//
// The following Swagger annotations are used to automatically generate API documentation:
// - @Summary: Brief description of the endpoint
// - @Description: Detailed explanation of the endpoint's functionality
// - @Tags: Groups this endpoint under the "auth" category in Swagger UI
// - @Accept: Specifies that the endpoint accepts JSON data
// - @Produce: Specifies that the endpoint returns JSON data
// - @Param: Defines the request body structure and validation rules
// - @Success: Documents successful response with status code and response model
// - @Failure: Documents error responses with status codes and error model
// - @Router: Specifies the HTTP method and path for the endpoint
//
// These annotations are automatically parsed by the docs.go file to generate
// the complete Swagger documentation at runtime.
//
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	response, err := h.userService.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	response.Status = "success"
	c.JSON(http.StatusOK, response)
}

// Register handles new user registration
//
// The Swagger annotations below follow the same pattern as the Login handler,
// but with different parameters and response models. The docs.go file uses
// reflection to:
// 1. Find all handler methods with Swagger annotations
// 2. Parse the annotations to extract endpoint details
// 3. Generate the complete Swagger paths object
// 4. Inject the generated paths into the Swagger template
//
// This approach ensures that the API documentation stays in sync with the
// actual implementation and reduces manual documentation maintenance.
//
// @Summary Register a new user
// @Description Register a new user with email, password, and name
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User registration data"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var registerReq models.RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
		return
	}

	response, err := h.userService.Register(registerReq)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "User with this email already exists"})
		} else {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		}
		return
	}

	response.Status = "success"
	c.JSON(http.StatusCreated, response)
}

// GetProfile returns the profile of the currently authenticated user
// @Summary Get user profile
// @Description Get the profile of the currently authenticated user
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	user, err := h.userService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// Return user profile without password
	profile := gin.H{
		"name":      user.Name,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}

	c.JSON(http.StatusOK, profile)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
