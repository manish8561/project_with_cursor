package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/internal/config"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"

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

func main() {
	// Load environment configuration
	cfg := config.LoadConfig()

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(EnableCORS)

	// Initialize services
	authHandler := handlers.NewAuthHandler(cfg.AuthServiceURL)
	userHandler := handlers.NewUserHandler(cfg.UserServiceURL)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "api-gateway"})
	})

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/validate", authHandler.ValidateToken)
			authRoutes.POST("/refresh", authHandler.RefreshToken)
		}

		// User routes (protected)
		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware(cfg.AuthServiceURL))
		{
			userRoutes.GET("/profile/:id", userHandler.GetUserByID)
			userRoutes.GET("/list", userHandler.ListUsers)
			userRoutes.PUT("/profile/:id", userHandler.UpdateUser)
			userRoutes.DELETE("/profile/:id", userHandler.DeleteUser)
		}
	}

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API Gateway starting on port %s", cfg.Port)
		if err := r.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start API Gateway: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down API Gateway...")
}
