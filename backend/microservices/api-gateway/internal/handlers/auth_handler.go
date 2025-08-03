package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication operations via proxy
type AuthHandler struct {
	authServiceURL string
}

// NewAuthHandler creates a new AuthHandler with the provided auth service URL
func NewAuthHandler(authServiceURL string) *AuthHandler {
	return &AuthHandler{
		authServiceURL: authServiceURL,
	}
}

// Login handles user login requests by proxying to auth service
func (h *AuthHandler) Login(c *gin.Context) {
	h.proxyRequest(c, "/api/auth/login", "POST")
}

// Register handles user registration requests by proxying to auth service
func (h *AuthHandler) Register(c *gin.Context) {
	h.proxyRequest(c, "/api/auth/register", "POST")
}

// ValidateToken handles token validation requests by proxying to auth service
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	h.proxyRequest(c, "/api/auth/validate", "POST")
}

// RefreshToken handles token refresh requests by proxying to auth service
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	h.proxyRequest(c, "/api/auth/refresh", "POST")
}

// proxyRequest forwards the request to the auth service
func (h *AuthHandler) proxyRequest(c *gin.Context, path, method string) {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create the request to the auth service
	req, err := http.NewRequest(method, h.authServiceURL+path, bytes.NewBuffer(body))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with auth service"})
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