package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations via proxy
type UserHandler struct {
	userServiceURL string
}

// NewUserHandler creates a new UserHandler with the provided user service URL
func NewUserHandler(userServiceURL string) *UserHandler {
	return &UserHandler{
		userServiceURL: userServiceURL,
	}
}

// GetUserByID handles requests to get a user by ID by proxying to user service
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	h.proxyRequest(c, "/api/users/profile/"+id, "GET")
}

// ListUsers handles requests to list users by proxying to user service
func (h *UserHandler) ListUsers(c *gin.Context) {
	h.proxyRequest(c, "/api/users/list", "GET")
}

// UpdateUser handles requests to update a user by proxying to user service
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	h.proxyRequest(c, "/api/users/profile/"+id, "PUT")
}

// DeleteUser handles requests to delete a user by proxying to user service
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	h.proxyRequest(c, "/api/users/profile/"+id, "DELETE")
}

// proxyRequest forwards the request to the user service
func (h *UserHandler) proxyRequest(c *gin.Context, path, method string) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create the request to the user service
	req, err := http.NewRequest(method, h.userServiceURL+path, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with user service"})
		return
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Return the response
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
} 