# Microservices Architecture

This directory contains the microservices implementation of the backend, breaking down the monolithic application into smaller, focused services.

## Architecture Overview

The application has been decomposed into the following microservices:

### 1. Auth Service (`auth-service`)
- **Port**: 8081
- **Purpose**: Handles user authentication, registration, and session management via HttpOnly cookies
- **Endpoints**:
  - `POST /api/auth/login` - User login (sets `access_token` cookie)
  - `POST /api/auth/register` - User registration (sets `access_token` cookie)
  - `GET /api/auth/me` - Get current session user ID from cookie
  - `POST /api/auth/logout` - Clear session cookie
  - `POST /api/auth/validate` - Token validation (cookie, body, or Bearer)
  - `POST /api/auth/refresh` - Refresh session (sets new cookie)

### 2. User Service (`user-service`)
- **Port**: 8082
- **Purpose**: Manages user data and profiles
- **Endpoints** (all require `access_token` cookie or Bearer token):
  - `GET /api/users/me` - Get current user's profile
  - `GET /api/users/profile/:id` - Get user profile by ID (must match session user)
  - `GET /api/users/list` - List users (paginated)
  - `PUT /api/users/profile/:id` - Update user
  - `DELETE /api/users/profile/:id` - Delete user

### 3. API Gateway (`api-gateway`)
- **Port**: 8080
- **Purpose**: Single entry point for all client requests, handles routing and authentication
- **Features**:
  - Request routing to appropriate services
  - Authentication middleware
  - CORS handling
  - Load balancing (future enhancement)

### 4. Shared Database
- **MongoDB**: Database-per-service model
- **Collections**: `auth_db.auth_users` (auth-service), `user_db.user_profiles` (user-service, synced via Kafka)

## Authentication

Sessions use an **HttpOnly cookie** named `access_token` containing a JWT. The browser sends this cookie automatically on subsequent requests; the frontend does not store tokens in `localStorage`.

| Setting | Env var | Default |
|---------|---------|---------|
| Cookie name | — | `access_token` |
| HttpOnly | — | `true` |
| SameSite | `COOKIE_SAME_SITE` | `Lax` |
| Secure | `COOKIE_SECURE` | `false` (set `true` in HTTPS production) |
| Max age | `COOKIE_MAX_AGE` | `86400` (24 hours) |
| Allowed origins | `ALLOWED_ORIGINS` | `http://localhost:4200,http://localhost:8085` |

Both auth-service and user-service must share the same `JWT_SECRET`. CORS is configured with explicit origins and `Access-Control-Allow-Credentials: true` (wildcard `*` is not used with cookies).

## Directory Structure

```
backend/
├── auth-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── config/
│       ├── handlers/
│       ├── logger/          # Zap structured logging
│       ├── models/
│       └── services/
├── user-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── config/
│       ├── handlers/
│       ├── logger/          # Zap structured logging
│       ├── models/
│       └── services/
├── api-gateway/
│   ├── cmd/
│   │   └── api-gateway/
│   │       └── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── config/
│       ├── handlers/
│       ├── logger/          # Zap structured logging
│       └── middleware/
└── README.md
```

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.21 or later (for local development)

### Running the Services

1. **Start all services with Docker Compose**:
   ```bash
   cd deploy
   docker compose up --build
   ```

2. **Access the services**:
   - API Gateway: http://localhost:8080
   - Auth Service: http://localhost:8081
   - User Service: http://localhost:8082
   - MongoDB: localhost:27017

### Testing Environment

1. **Start test environment**:
   ```bash
   cd deploy
   docker compose -f docker-compose.test.yml up --build
   ```

2. **Access test services**:
   - API Gateway: http://localhost:8080
   - Auth Service: http://localhost:8081
   - User Service: http://localhost:8082
   - MongoDB: localhost:27017

### Development

1. **Build and run individual services**:
   ```bash
   # Auth Service
   cd auth-service
   go mod tidy
   go run main.go

   # User Service
   cd user-service
   go mod tidy
   go run main.go

   # API Gateway
   cd api-gateway
   go mod tidy
   go run main.go
   ```

2. **Environment Variables**:
   ```bash
   # Auth Service
   PORT=8081
   MONGO_URI=mongodb://localhost:27017
   MONGO_DB=auth_db
   JWT_SECRET=your-secret-key
   COOKIE_SECURE=false
   COOKIE_SAME_SITE=Lax
   ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8085
   LOG_LEVEL=info

   # User Service
   PORT=8082
   MONGO_URI=mongodb://localhost:27017
   MONGO_DB=user_db
   JWT_SECRET=your-secret-key
   ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8085
   LOG_LEVEL=info

   # API Gateway
   PORT=8080
   AUTH_SERVICE_URL=http://localhost:8081
   USER_SERVICE_URL=http://localhost:8082
   LOG_LEVEL=info
   ```

## API Endpoints

### Authentication (via API Gateway)

Login and register set an `access_token` HttpOnly cookie. Use `-c` / `-b` with curl to persist and send cookies.

```bash
# Login (saves cookie to cookies.txt)
curl -c cookies.txt -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Register
curl -c cookies.txt -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"user@example.com","password":"password","confirmPassword":"password"}'

# Current session
curl -b cookies.txt http://localhost:8080/api/auth/me

# Logout
curl -b cookies.txt -X POST http://localhost:8080/api/auth/logout

# Validate session (cookie sent automatically with -b)
curl -b cookies.txt -X POST http://localhost:8080/api/auth/validate

# Refresh session (sets a new cookie)
curl -b cookies.txt -c cookies.txt -X POST http://localhost:8080/api/auth/refresh
```

Legacy Bearer/body token usage is still supported for validate and refresh:

```bash
curl -X POST http://localhost:8080/api/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"your-jwt-token"}'
```

### User Management (via API Gateway)

```bash
# Get current user profile (requires session cookie)
curl -b cookies.txt http://localhost:8080/api/users/me

# Get user profile by ID (must match session user)
curl -b cookies.txt http://localhost:8080/api/users/profile/USER_ID

# List users (requires authentication)
curl -b cookies.txt http://localhost:8080/api/users/list
```

## Rate Limiting (API Gateway)

The API Gateway enforces a simple per-IP token bucket rate limit on proxied routes (e.g., `/api/auth/*`, `/api/users/*`). Defaults can be tuned via environment variables:

- `RATE_LIMIT_RPS`: Requests per second refill rate (float, default: 10)
- `RATE_LIMIT_BURST`: Burst capacity in requests (int, default: 20)

When the limit is exceeded, requests receive HTTP 429 Too Many Requests with body `rate limit exceeded`.

### Examples

```bash
# Allow a higher throughput during load testing
RATE_LIMIT_RPS=50 RATE_LIMIT_BURST=100 \
  docker compose -f deploy/docker-compose.yml up -d --build
```

## Docker Commands

### Production
```bash
# Start all services
cd deploy
docker compose up --build

# Stop services
docker compose down

# View logs
docker compose logs

# View specific service logs
docker compose logs auth-service
docker compose logs user-service
docker compose logs api-gateway
```

### Testing
```bash
# Start test environment
cd deploy
docker compose -f docker-compose.test.yml up --build

# Stop test services
docker compose -f docker-compose.test.yml down

# View test logs
docker compose -f docker-compose.test.yml logs
```

### Individual Services
```bash
# Build individual services
docker build -t auth-service auth-service/
docker build -t user-service user-service/
docker build -t api-gateway api-gateway/

# Run individual services
docker run -p 8081:8081 auth-service
docker run -p 8082:8082 user-service
docker run -p 8080:8080 api-gateway
```

## Logging

All services implement structured logging using **Zap** for high-performance, structured JSON logging.

### Features
- **Structured JSON Output**: All logs are in JSON format for easy parsing
- **Environment-based Log Levels**: Configure via `LOG_LEVEL` environment variable
- **Service Identification**: Each log entry includes the service name
- **Performance Optimized**: Uses Uber's Zap logger for minimal overhead

### Log Levels
Set the `LOG_LEVEL` environment variable to control logging verbosity:
- `debug`: Most verbose, includes debug information
- `info`: General information (default)
- `warn`: Warning messages
- `error`: Error messages only

### Example Log Output
```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "caller": "main.go:25",
  "message": "Server starting",
  "service": "auth-service",
  "port": 8081
}
```

### Usage in Code
```go
import "your-service/internal/logger"

// Initialize logger (done in main.go)
logger.InitLogger()

// Use logger
logger.GetLogger().Info("User authenticated", 
    zap.String("user_id", userID),
    zap.String("email", email),
)

logger.GetLogger().Error("Database connection failed",
    zap.Error(err),
    zap.String("database", "mongodb"),
)
```

## Benefits of This Architecture

1. **Scalability**: Each service can be scaled independently
2. **Maintainability**: Smaller, focused codebases
3. **Technology Flexibility**: Each service can use different technologies
4. **Fault Isolation**: Failure in one service doesn't affect others
5. **Team Organization**: Different teams can work on different services
6. **Observability**: Structured logging across all services for better monitoring

## Future Enhancements

1. **Service Discovery**: Implement service discovery (Consul, etcd)
2. **Load Balancing**: Add load balancers for each service
3. **Circuit Breakers**: Implement circuit breakers for service communication
4. **Distributed Tracing**: Add tracing (Jaeger, Zipkin)
5. **Monitoring**: Implement metrics and monitoring (Prometheus, Grafana)
6. **Message Queues**: Add async communication (RabbitMQ, Kafka)
7. **API Documentation**: Add Swagger/OpenAPI documentation
8. **Testing**: Add comprehensive test suites for each service

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 8080, 8081, 8082, and 27017 are available
2. **MongoDB connection**: Check if MongoDB is running and accessible
3. **Service communication**: Verify service URLs in API Gateway configuration
4. **Authentication errors**: Check `JWT_SECRET` matches across auth-service and user-service; verify `ALLOWED_ORIGINS` includes your frontend origin when using credentialed CORS
5. **Cookie not sent**: Ensure requests use credentials (`withCredentials: true` in browsers); in Docker, use the frontend at `http://localhost:8085` so `/api` is same-origin

### Logs
```bash
# View all service logs
docker compose logs

# View specific service logs
docker compose logs auth-service
docker compose logs user-service
docker compose logs api-gateway

# Follow logs in real-time
docker compose logs -f
```

### Health Checks
```bash
# Check service health
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # User Service
```

## Environment Variables

### Production Environment
```bash
# MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password
MONGO_INITDB_DATABASE=auth_db

# Auth Service
PORT=8081
MONGO_URI=mongodb://admin:password@mongodb:27017
MONGO_DB=auth_db
JWT_SECRET=your-super-secret-jwt-key
COOKIE_SECURE=false
COOKIE_SAME_SITE=Lax
ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8085
LOG_LEVEL=info

# User Service
PORT=8082
MONGO_URI=mongodb://admin:password@mongodb:27017
MONGO_DB=user_db
JWT_SECRET=your-super-secret-jwt-key
ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8085
LOG_LEVEL=info

# API Gateway
PORT=8080
AUTH_SERVICE_URL=http://auth-service:8081
USER_SERVICE_URL=http://user-service:8082
LOG_LEVEL=info
```

### Test Environment
```bash
# MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password123
MONGO_INITDB_DATABASE=testdb

# Auth Service
PORT=8081
MONGO_URI=mongodb://admin:password123@mongodb:27017
MONGO_DB=testdb
JWT_SECRET=test-jwt-secret-key-for-testing
COOKIE_SECURE=false
COOKIE_SAME_SITE=Lax
ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8085
LOG_LEVEL=debug

# User Service
PORT=8082
MONGO_URI=mongodb://admin:password123@mongodb:27017
MONGO_DB=testdb
JWT_SECRET=test-jwt-secret-key-for-testing
ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8085
LOG_LEVEL=debug

# API Gateway
PORT=8080
AUTH_SERVICE_URL=http://auth-service:8081
USER_SERVICE_URL=http://user-service:8082
LOG_LEVEL=debug
``` 