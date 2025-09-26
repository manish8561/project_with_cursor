package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/logger"
	"auth-service/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EnableCORS is a middleware function that enables CORS for all routes
func EnableCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}

// SetupRoutes configures all routes for the auth service
func SetupRoutes(authHandler *handlers.AuthHandler, log logger.Logger) *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(EnableCORS)
	
	// User middleware
	r.Use(middleware.ZapMiddleware(log))

	// API routes
	api := r.Group("/api/auth")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)
		api.POST("/validate", authHandler.ValidateToken)
		api.POST("/refresh", authHandler.RefreshToken)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "auth-service"})
	})

	return r
}
