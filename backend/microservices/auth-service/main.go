package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/services"

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

	// Initialize MongoDB
	mongoConfig, err := config.NewMongoDBConfig(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoConfig.Close()

	// Initialize JWT service
	jwtConfig := config.NewJWTConfig(cfg.JWTSecret)
	jwtService := services.NewJWTService(jwtConfig)

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(EnableCORS)

	// Initialize services
	authService := services.NewAuthService(mongoConfig, jwtService)
	authHandler := handlers.NewAuthHandler(authService)

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

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Auth service starting on port %s", cfg.Port)
		if err := r.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start auth service: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down auth service...")
}
