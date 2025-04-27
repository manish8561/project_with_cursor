package main

import (
	"fmt"
	"net/http"

	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment configuration
	cfg := config.LoadConfig()

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Initialize services
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	r.POST("/api/login", userHandler.Login)
	r.POST("/api/register", userHandler.Register)

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	r.Run(serverAddr) // listen and serve on configured port
}
