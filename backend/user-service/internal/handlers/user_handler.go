package handlers

import (
	"user-service/internal/logger"
	"user-service/internal/models"
	"user-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// GetUserByID handles requests to get a user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	id := c.Param("id")
	if id == "" {
		zapLogger.Error("Missing user ID in request", 
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	zapLogger.Info("Getting user by ID", 
		zap.String("user_id", id),
		zap.String("client_ip", c.ClientIP()),
	)

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		zapLogger.Warn("User not found", 
			zap.String("user_id", id),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("User retrieved successfully", 
		zap.String("user_id", id),
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