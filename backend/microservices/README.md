# Microservices Architecture

This directory contains the microservices implementation of the backend, breaking down the monolithic application into smaller, focused services.

## Architecture Overview

The application has been decomposed into the following microservices:

### 1. Auth Service (`auth-service`)
- **Port**: 8081
- **Purpose**: Handles user authentication, registration, and JWT token management
- **Endpoints**:
  - `POST /api/auth/login` - User login
  - `POST /api/auth/register` - User registration
  - `POST /api/auth/validate` - Token validation
  - `POST /api/auth/refresh` - Token refresh

### 2. User Service (`user-service`)
- **Port**: 8082
- **Purpose**: Manages user data and profiles
- **Endpoints**:
  - `GET /api/users/profile/:id` - Get user by ID
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
- **MongoDB**: Single database shared across services
- **Collections**: `users` (shared between auth and user services)

## Directory Structure

```
microservices/
├── auth-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── config/
│       ├── handlers/
│       ├── models/
│       └── services/
├── user-service/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── config/
│       ├── handlers/
│       ├── models/
│       └── services/
├── api-gateway/
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── config/
│       ├── handlers/
│       └── middleware/
├── docker-compose.yml
└── README.md
```

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.21 or later (for local development)

### Running the Services

1. **Start all services with Docker Compose**:
   ```bash
   cd backend/microservices
   docker compose up -d
   ```

2. **Access the services**:
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
   - `PORT`: Service port (default: 8080, 8081, 8082)
   - `MONGO_URI`: MongoDB connection string
   - `MONGO_DB`: Database name
   - `JWT_SECRET`: Secret key for JWT tokens
   - `AUTH_SERVICE_URL`: Auth service URL (for API Gateway)
   - `USER_SERVICE_URL`: User service URL (for API Gateway)

## API Endpoints

### Authentication (via API Gateway)
```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"user@example.com","password":"password","confirmPassword":"password"}'

# Validate Token
curl -X POST http://localhost:8080/api/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"your-jwt-token"}'
```

### User Management (via API Gateway)
```bash
# Get user profile (requires authentication)
curl -X GET http://localhost:8080/api/users/profile/user-id \
  -H "Authorization: Bearer your-jwt-token"

# List users (requires authentication)
curl -X GET http://localhost:8080/api/users/list \
  -H "Authorization: Bearer your-jwt-token"
```

## Benefits of This Architecture

1. **Scalability**: Each service can be scaled independently
2. **Maintainability**: Smaller, focused codebases
3. **Technology Flexibility**: Each service can use different technologies
4. **Fault Isolation**: Failure in one service doesn't affect others
5. **Team Organization**: Different teams can work on different services

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

### Logs
```bash
# View all service logs
docker-compose logs

# View specific service logs
docker-compose logs auth-service
docker-compose logs user-service
docker-compose logs api-gateway
```

### Health Checks
```bash
# Check service health
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # User Service
``` 