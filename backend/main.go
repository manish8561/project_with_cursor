package main

import (
	"net/http"

	"backend/internal/handlers"
	"backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
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

	// Login endpoint
	r.POST("/api/login", userHandler.Login)

	// Start the server
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
