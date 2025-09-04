package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/logger"
	"auth-service/internal/services"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	zapLogger := logger.GetLogger()

	// Load environment configuration
	cfg := config.LoadConfig()
	zapLogger.Info("Configuration loaded successfully")

	// Initialize MongoDB
	mongoConfig, err := config.NewMongoDBConfig(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		zapLogger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoConfig.Close()
	zapLogger.Info("MongoDB connection established")

	// Initialize JWT service
	jwtConfig := config.NewJWTConfig(cfg.JWTSecret)
	jwtService := services.NewJWTService(jwtConfig)
	zapLogger.Info("JWT service initialized")

	// Initialize services
	authService := services.NewAuthService(mongoConfig, jwtService)
	authHandler := handlers.NewAuthHandler(authService)
	zapLogger.Info("Auth service and handlers initialized")

	// Setup routes using the router
	r := SetupRoutes(authHandler)

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		zapLogger.Info("Auth service starting", zap.String("port", cfg.Port))
		if err := r.Run(serverAddr); err != nil {
			zapLogger.Fatal("Failed to start auth service", zap.Error(err))
		}
	}()

	<-quit
	zapLogger.Info("Shutting down auth service...")
}
