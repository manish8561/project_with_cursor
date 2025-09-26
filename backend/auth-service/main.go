package main

import (
	"fmt"
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
