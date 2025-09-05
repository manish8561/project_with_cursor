# API Gateway Service

The API Gateway serves as the single entry point for all client requests in the microservices architecture. It handles request routing, authentication, and provides a unified interface for the backend services.

## Overview

- **Port**: 8080
- **Framework**: Kratos v2
- **Purpose**: Route requests to appropriate microservices (auth-service, user-service)
- **Features**: Request routing, authentication middleware, CORS handling

## Architecture

The API Gateway routes requests to:
- **Auth Service** (port 8081): Authentication and user management
- **User Service** (port 8082): User profile operations

## Directory Structure

```
api-gateway/
├── cmd/
│   └── api-gateway/
│       ├── main.go          # Application entry point
│       ├── wire.go          # Dependency injection
│       └── wire_gen.go      # Generated wire code
├── internal/
│   ├── biz/                 # Business logic
│   ├── conf/                # Configuration
│   ├── data/                # Data layer
│   ├── logger/              # Zap structured logging
│   ├── server/              # HTTP/gRPC servers
│   └── service/             # Service layer
├── api/                     # Protocol buffer definitions
├── configs/                 # Configuration files
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## Logging

The API Gateway implements structured logging using **Zap** for high-performance, structured JSON logging.

### Features
- **Structured JSON Output**: All logs are in JSON format for easy parsing
- **Environment-based Log Levels**: Configure via `LOG_LEVEL` environment variable
- **Service Identification**: Each log entry includes "service": "api-gateway"
- **Performance Optimized**: Uses Uber's Zap logger for minimal overhead

### Log Levels
Set the `LOG_LEVEL` environment variable:
- `debug`: Most verbose, includes debug information
- `info`: General information (default)
- `warn`: Warning messages
- `error`: Error messages only

### Example Log Output
```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "caller": "main.go:72",
  "message": "API Gateway starting",
  "service": "api-gateway",
  "service.id": "hostname",
  "service.name": "api-gateway",
  "service.version": "v1.0.0"
}
```

### Usage in Code
```go
import (
    "api-gateway/internal/logger"
    "go.uber.org/zap"
)

// Initialize logger (done in main.go)
logger.InitLogger()

// Use logger
logger.GetLogger().Info("Request received", 
    zap.String("method", "POST"),
    zap.String("path", "/api/auth/login"),
    zap.String("client_ip", clientIP),
)

logger.GetLogger().Error("Service unavailable",
    zap.Error(err),
    zap.String("service", "auth-service"),
    zap.String("url", serviceURL),
)
```

## Environment Variables

```bash
# Server Configuration
PORT=8080

# Service URLs
AUTH_SERVICE_URL=http://auth-service:8081
USER_SERVICE_URL=http://user-service:8082

# Logging
LOG_LEVEL=info  # debug, info, warn, error

# Optional: Service identification
SERVICE_NAME=api-gateway
SERVICE_VERSION=v1.0.0
```

## Development

### Prerequisites
- Go 1.21 or later
- Docker (optional)

### Running Locally
```bash
# Install dependencies
go mod tidy

# Run the service
go run cmd/api-gateway/main.go

# Or build and run
go build -o bin/api-gateway cmd/api-gateway/main.go
./bin/api-gateway
```

### Building
```bash
# Build binary
go build -o bin/api-gateway cmd/api-gateway/main.go

# Build Docker image
docker build -t api-gateway .
```

## Docker

### Build and Run
```bash
# Build image
docker build -t api-gateway .

# Run container
docker run -p 8080:8080 \
  -e AUTH_SERVICE_URL=http://auth-service:8081 \
  -e USER_SERVICE_URL=http://user-service:8082 \
  -e LOG_LEVEL=info \
  api-gateway
```

### Docker Compose
The service is configured to run with Docker Compose alongside other microservices:

```bash
# Start all services
docker compose -f ../deploy/docker-compose.yml up --build

# View logs
docker compose -f ../deploy/docker-compose.yml logs api-gateway
```

## API Endpoints

The API Gateway proxies requests to the appropriate services:

### Health Check
- `GET /health` - API Gateway health status

### Authentication (proxied to auth-service)
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration  
- `POST /api/auth/validate` - Token validation
- `POST /api/auth/refresh` - Token refresh
- `GET /api/auth/health` - Auth service health

### User Management (proxied to user-service)
- `GET /api/users/profile/{id}` - Get user profile
- `GET /api/users/list` - List users
- `GET /api/users/health` - User service health

## Monitoring and Debugging

### Health Checks
```bash
# Check API Gateway health
curl http://localhost:8080/health

# Check service health via gateway
curl http://localhost:8080/api/auth/health
curl http://localhost:8080/api/users/health
```

### Logs
```bash
# View real-time logs
docker compose logs -f api-gateway

# Filter by log level (if using structured logging tools)
docker compose logs api-gateway | jq 'select(.level=="error")'
```

## Kratos Framework

This service is built using the Kratos framework. For Kratos-specific operations:

### Install Kratos CLI
```bash
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

### Generate Code
```bash
# Generate API files from proto
make api

# Generate wire dependency injection
cd cmd/api-gateway
wire

# Generate all files
make all
```

