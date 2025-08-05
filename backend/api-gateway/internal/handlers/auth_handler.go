package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authServiceURL string
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authServiceURL string) *AuthHandler {
	return &AuthHandler{
		authServiceURL: authServiceURL,
	}
}

// LoginRequest represents the login request
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/login [post]
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the registration request
// @Summary User registration
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid data"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/auth/register [post]
type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// ValidateRequest represents the token validation request
// @Summary Validate JWT token
// @Description Validate JWT token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ValidateRequest true "Token to validate"
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 400 {object} map[string]interface{} "Invalid token"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/auth/validate [post]
type ValidateRequest struct {
	Token string `json:"token" binding:"required"`
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward request to auth service
	resp, err := h.forwardRequest(c.Request.Method, "/api/auth/login", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward request to auth service
	resp, err := h.forwardRequest(c.Request.Method, "/api/auth/register", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	c.Data(http.StatusCreated, "application/json", resp)
}

// ValidateToken handles token validation
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward request to auth service
	resp, err := h.forwardRequest(c.Request.Method, "/api/auth/validate", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// forwardRequest forwards the request to the auth service
func (h *AuthHandler) forwardRequest(method, path string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", h.authServiceURL, path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	responseBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
