package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userServiceURL string
}

// NewUserHandler creates a new user handler
func NewUserHandler(userServiceURL string) *UserHandler {
	return &UserHandler{
		userServiceURL: userServiceURL,
	}
}

// GetUserProfileRequest represents the get user profile request
// @Summary Get user profile
// @Description Get user profile by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User profile"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/users/profile/{id} [get]
type GetUserProfileRequest struct {
	ID string `uri:"id" binding:"required"`
}

// ListUsersRequest represents the list users request
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of users per page" default(10)
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of users"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/users/list [get]
type ListUsersRequest struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// UpdateUserRequest represents the update user request
// @Summary Update user profile
// @Description Update user profile information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "User update data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/users/profile/{id} [put]
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// GetUserProfile handles getting user profile
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	var req GetUserProfileRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward request to user service
	resp, err := h.forwardRequest(c.Request.Method, fmt.Sprintf("/api/users/profile/%s", req.ID), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// ListUsers handles listing users
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	// Forward request to user service
	resp, err := h.forwardRequest(c.Request.Method, "/api/users/list", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// UpdateUserProfile handles updating user profile
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Forward request to user service
	resp, err := h.forwardRequest(c.Request.Method, fmt.Sprintf("/api/users/profile/%s", userID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	c.Data(http.StatusOK, "application/json", resp)
}

// forwardRequest forwards the request to the user service
func (h *UserHandler) forwardRequest(method, path string, body interface{}) ([]byte, error) {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s%s", h.userServiceURL, path)
	var req *http.Request

	if jsonBody != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

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