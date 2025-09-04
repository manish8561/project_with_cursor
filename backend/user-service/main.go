package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"user-service/internal/config"
	"user-service/internal/handlers"
	"user-service/internal/logger"
	"user-service/internal/services"
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

	// Initialize services
	userService := services.NewUserService(mongoConfig)
	userHandler := handlers.NewUserHandler(userService)
	zapLogger.Info("User service and handlers initialized")

	// Setup routes using the router
	r := SetupRoutes(userHandler)

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		zapLogger.Info("User service starting", zap.String("port", cfg.Port))
		if err := r.Run(serverAddr); err != nil {
			zapLogger.Fatal("Failed to start user service", zap.Error(err))
		}
	}()

	<-quit
	zapLogger.Info("Shutting down user service...")
}
