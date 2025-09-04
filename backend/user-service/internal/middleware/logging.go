package middleware

import (
	"user-service/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware creates a Gin middleware for request logging using Zap
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		zapLogger := logger.GetLogger()
		
		// Log the request details
		zapLogger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
			zap.Int("body_size", param.BodySize),
		)
		
		// Return empty string since we're using Zap for logging
		return ""
	})
}

// ZapMiddleware creates a more comprehensive Zap logging middleware
func ZapMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		
		zapLogger := logger.GetLogger()
		
		// Log incoming request
		zapLogger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
		
		// Process request
		c.Next()
		
		// Log response
		latency := time.Since(start)
		zapLogger.Info("Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("response_size", c.Writer.Size()),
		)
		
		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				zapLogger.Error("Request error",
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.Error(err),
					zap.String("client_ip", c.ClientIP()),
				)
			}
		}
	}
}
