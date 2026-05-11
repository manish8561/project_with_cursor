package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"user-service/internal/config"
	"user-service/internal/handlers"
	"user-service/internal/logger"
	"user-service/internal/middleware"
	"user-service/internal/services"
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
		os.Exit(1)
	}
	defer func() {
		if closeErr := mongoConfig.Close(); closeErr != nil {
			log.Error("MongoDB close error", zap.Error(closeErr))
		}
	}()
	log.Info("MongoDB connection established")

	// Initialize Kafka publisher for user lifecycle events.
	publisher, err := services.NewKafkaPublisher(
		cfg.KafkaBrokers,
		cfg.KafkaClientID,
		cfg.KafkaTopicUserCreated,
		cfg.KafkaTopicUserUpdated,
		cfg.KafkaTopicUserDeleted,
	)
	if err != nil {
		log.Error("Failed to initialize Kafka publisher", zap.Error(err))
	}
	defer func() {
		if publisher != nil {
			if closeErr := publisher.Close(); closeErr != nil {
				log.Error("Failed to close Kafka publisher", zap.Error(closeErr))
			}
		}
	}()

	// Initialize services
	userService := services.NewUserService(mongoConfig, publisher)
	userHandler := handlers.NewUserHandler(userService, log)
	log.Info("User service and handlers initialized")

	// Initialize Kafka consumer for user lifecycle events.
	consumer, err := services.NewUserEventConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaGroupID,
		cfg.KafkaClientID,
		cfg.KafkaTopicUserCreated,
		cfg.KafkaTopicUserUpdated,
		cfg.KafkaTopicUserDeleted,
		userService,
		log,
	)
	if err != nil {
		log.Error("Failed to initialize Kafka consumer", zap.Error(err))
	}
	defer func() {
		if consumer != nil {
			if closeErr := consumer.Close(); closeErr != nil {
				log.Error("Failed to close Kafka consumer", zap.Error(closeErr))
			}
		}
	}()

	consumerCtx, cancelConsumer := context.WithCancel(context.Background())
	defer cancelConsumer()
	go func() {
		if consumer != nil {
			consumer.Start(consumerCtx)
		}
	}()

	// Setup routes using the router
	r := SetupRoutes(userHandler, log)

	// Start the server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("User service starting", zap.String("port", cfg.Port))
		if err := r.Run(serverAddr); err != nil {
			log.Error("Failed to start user service", zap.Error(err))
		}
	}()

	<-quit
	cancelConsumer()
	log.Info("Shutting down user service...")
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

// SetupRoutes configures all routes for the user service
func SetupRoutes(userHandler *handlers.UserHandler, log logger.Logger) *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(EnableCORS)

	// User middleware
	r.Use(middleware.ZapMiddleware(log))

	// API routes
	api := r.Group("/api/users")
	{
		api.GET("/profile/:id", userHandler.GetUserByID)
		api.GET("/list", userHandler.ListUsers)
		api.PUT("/profile/:id", userHandler.UpdateUser)
		api.DELETE("/profile/:id", userHandler.DeleteUser)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "user-service"})
	})

	return r
}
