package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"user-service/internal/logger"
	"user-service/internal/models"
	"user-service/internal/services"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new UserHandler with the provided user service
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserByID handles requests to get the current user's profile using JWT token
func (h *UserHandler) GetUserByID(c *gin.Context) {
	zapLogger := logger.GetLogger()

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		zapLogger.Error("Missing Authorization header",
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token is required"})
		return
	}

	// Extract the token from "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader { // No Bearer prefix found
		zapLogger.Error("Invalid Authorization header format")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		return
	}

	// Parse the token to get user ID
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Replace with your actual secret key
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		zapLogger.Error("Invalid token", 
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		zapLogger.Error("Invalid token claims")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	// Get user ID from token claims
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		zapLogger.Error("User ID not found in token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
		return
	}

	zapLogger.Info("Getting user by ID from token", 
		zap.String("user_id", userID),
		zap.String("client_ip", c.ClientIP()),
	)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		zapLogger.Warn("User not found", 
			zap.String("user_id", userID),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("User retrieved successfully", 
		zap.String("user_id", userID),
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusOK, user)
}

// ListUsers handles requests to list users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	zapLogger.Info("Listing users", 
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("client_ip", c.ClientIP()),
	)

	response, err := h.userService.ListUsers(page, size)
	if err != nil {
		zapLogger.Error("Failed to list users", 
			zap.Error(err),
			zap.Int("page", page),
			zap.Int("size", size),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("Users listed successfully", 
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusOK, response)
}

// UpdateUser handles requests to update a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	id := c.Param("id")
	if id == "" {
		zapLogger.Error("Missing user ID in update request", 
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zapLogger.Error("Failed to bind update user request", 
			zap.Error(err),
			zap.String("user_id", id),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("Updating user", 
		zap.String("user_id", id),
		zap.String("client_ip", c.ClientIP()),
	)

	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		zapLogger.Warn("Failed to update user", 
			zap.String("user_id", id),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("User updated successfully", 
		zap.String("user_id", id),
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles requests to delete a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	id := c.Param("id")
	if id == "" {
		zapLogger.Error("Missing user ID in delete request", 
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	zapLogger.Info("Deleting user", 
		zap.String("user_id", id),
		zap.String("client_ip", c.ClientIP()),
	)

	err := h.userService.DeleteUser(id)
	if err != nil {
		zapLogger.Warn("Failed to delete user", 
			zap.String("user_id", id),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("User deleted successfully", 
		zap.String("user_id", id),
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
} 