package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates authentication middleware that validates tokens via auth service
func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		// Validate token with auth service
		valid, userID, err := validateTokenWithAuthService(authServiceURL, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
			return
		}

		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user ID in context
		c.Set("userID", userID)
		c.Next()
	}
}

// validateTokenWithAuthService validates a token by calling the auth service
func validateTokenWithAuthService(authServiceURL, token string) (bool, string, error) {
	// Create request body
	requestBody := map[string]string{"token": token}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, "", err
	}

	// Create request to auth service
	req, err := http.NewRequest("POST", authServiceURL+"/api/auth/validate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	// Parse response
	var response struct {
		Valid  bool   `json:"valid"`
		UserID string `json:"user_id"`
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return false, "", err
	}

	return response.Valid, response.UserID, nil
} 