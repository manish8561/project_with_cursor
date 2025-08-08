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

// @title API Gateway
// @version 1.0
// @description API Gateway for microservices architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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

	// API documentation endpoint
	r.GET("/swagger", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Gateway Documentation",
			"endpoints": gin.H{
				"health": "/health",
				"auth": gin.H{
					"login":    "POST /api/auth/login",
					"register": "POST /api/auth/register",
					"validate": "POST /api/auth/validate",
				},
				"users": gin.H{
					"profile": "GET /api/users/profile/:id",
					"list":    "GET /api/users/list",
					"update":  "PUT /api/users/profile/:id",
				},
			},
		})
	})

	// API documentation index
	r.GET("/docs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docs.html", gin.H{
			"title": "API Documentation",
			"services": []map[string]string{
				{
					"name":        "Authentication Service",
					"description": "User authentication and registration",
					"endpoints":   "POST /api/auth/login, POST /api/auth/register, POST /api/auth/validate",
				},
				{
					"name":        "User Service",
					"description": "User profile management",
					"endpoints":   "GET /api/users/profile/:id, GET /api/users/list, PUT /api/users/profile/:id",
				},
			},
		})
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
		}

		// User routes (protected)
		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware(cfg.AuthServiceURL))
		{
			userRoutes.GET("/profile/:id", userHandler.GetUserProfile)
			userRoutes.GET("/list", userHandler.ListUsers)
			userRoutes.PUT("/profile/:id", userHandler.UpdateUserProfile)
		}
	}

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("API Gateway starting on port %s", cfg.Port)
		log.Printf("API documentation available at: http://localhost:%s/swagger", cfg.Port)
		log.Printf("API documentation available at: http://localhost:%s/docs", cfg.Port)
		if err := r.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start API Gateway: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down API Gateway...")
}
