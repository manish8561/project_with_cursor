package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"user-service/internal/config"
	"user-service/internal/handlers"
	"user-service/internal/services"
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

	// Initialize services
	userService := services.NewUserService(mongoConfig)
	userHandler := handlers.NewUserHandler(userService)

	// Setup routes using the router
	r := SetupRoutes(userHandler)

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("User service starting on port %s", cfg.Port)
		if err := r.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start user service: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down user service...")
}
