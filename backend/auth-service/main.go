package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/services"
)

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

	// Initialize services
	authService := services.NewAuthService(mongoConfig, jwtService)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup routes using the router
	r := SetupRoutes(authHandler)

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
