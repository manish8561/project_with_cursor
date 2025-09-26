package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/logger"
	"auth-service/internal/middleware"
	"auth-service/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log, err := logger.NewZapLogger()
	if err != nil {
		log.Error("Failed to initialize logger", zap.Error(err))
	}
	defer log.Sync()

	// Load environment configuration
	cfg := config.LoadConfig()
	log.Info("Configuration loaded successfully")

	// Initialize MongoDB
	mongoConfig, err := config.NewMongoDBConfig(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Error("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoConfig.Close()
	log.Info("MongoDB connection established")

	// Initialize JWT service
	jwtConfig := config.NewJWTConfig(cfg.JWTSecret)
	jwtService := services.NewJWTService(jwtConfig)
	log.Info("JWT service initialized")

	// Initialize services
	authService := services.NewAuthService(mongoConfig, jwtService)
	authHandler := handlers.NewAuthHandler(authService, log)
	log.Info("Auth service and handlers initialized")

	// Setup routes using the router
	r := SetupRoutes(authHandler, log)

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Auth service starting", zap.String("port", cfg.Port))
		if err := r.Run(serverAddr); err != nil {
			log.Error("Failed to start auth service", zap.Error(err))
		}
	}()

	<-quit
	log.Info("Shutting down auth service...")
}

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
