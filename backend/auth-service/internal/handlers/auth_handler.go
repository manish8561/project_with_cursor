package handlers

import (
	"auth-service/internal/logger"
	"auth-service/internal/models"
	"auth-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler with the provided auth service
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zapLogger.Error("Failed to bind login request", 
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("Login attempt", 
		zap.String("email", req.Email),
		zap.String("client_ip", c.ClientIP()),
	)

	response, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		zapLogger.Warn("Login failed", 
			zap.String("email", req.Email),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("Login successful", 
		zap.String("email", req.Email),
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusOK, response)
}

// Register handles user registration requests
func (h *AuthHandler) Register(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zapLogger.Error("Failed to bind register request", 
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("Registration attempt", 
		zap.String("email", req.Email),
		zap.String("client_ip", c.ClientIP()),
	)

	response, err := h.authService.Register(req)
	if err != nil {
		zapLogger.Warn("Registration failed", 
			zap.String("email", req.Email),
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("Registration successful", 
		zap.String("email", req.Email),
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusCreated, response)
}

// ValidateToken handles token validation requests
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	var req models.TokenValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zapLogger.Error("Failed to bind token validation request", 
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Debug("Token validation attempt", 
		zap.String("client_ip", c.ClientIP()),
	)

	response, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		zapLogger.Warn("Token validation failed", 
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Debug("Token validation successful", 
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	zapLogger := logger.GetLogger()
	
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zapLogger.Error("Failed to bind refresh token request", 
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Debug("Token refresh attempt", 
		zap.String("client_ip", c.ClientIP()),
	)

	response, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		zapLogger.Warn("Token refresh failed", 
			zap.Error(err),
			zap.String("client_ip", c.ClientIP()),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	zapLogger.Info("Token refresh successful", 
		zap.String("client_ip", c.ClientIP()),
	)
	c.JSON(http.StatusOK, response)
}
