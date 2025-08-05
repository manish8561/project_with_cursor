# Full Stack Microservices Project

A full-stack application with Angular frontend and Go microservices backend.

## Project Structure

```
.
├── backend/           # Go microservices backend
│   ├── auth-service/  # Authentication service
│   ├── user-service/  # User management service
│   ├── api-gateway/   # API Gateway service
│   └── README.md      # Backend documentation
├── frontend/          # Angular frontend
└── deploy/           # Deployment configurations
    ├── docker-compose.yml      # Production deployment
    └── docker-compose.test.yml # Testing deployment
```

## Prerequisites

- Docker and Docker Compose
- Node.js (for frontend development)
- Go 1.21+ (for backend development)

## Quick Start

### Using Docker Compose (Recommended)

1. **Start all services**:
   ```bash
   cd deploy
   docker compose up --build
   ```

2. **Access the services**:
   - Frontend: http://localhost:4200
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

## API Documentation

### Swagger Documentation

The API Gateway provides centralized Swagger documentation for all microservices:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Documentation**: http://localhost:8080/docs
- **OpenAPI JSON**: http://localhost:8080/swagger/doc.json
- **OpenAPI YAML**: http://localhost:8080/swagger/doc.yaml

### Documentation Features

✅ **Centralized Documentation**: Single entry point for all API documentation
✅ **Interactive Testing**: Test APIs directly from the Swagger UI
✅ **Authentication Support**: JWT Bearer token authentication
✅ **Request/Response Examples**: Detailed examples for all endpoints
✅ **Error Codes**: Comprehensive error response documentation
✅ **API Versioning**: Version control for API changes

### Documentation Structure

The API documentation is organized by service:

#### Authentication Service (`/api/auth`)
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration  
- `POST /api/auth/validate` - Token validation

#### User Service (`/api/users`)
- `GET /api/users/profile/:id` - Get user profile
- `GET /api/users/list` - List users (paginated)
- `PUT /api/users/profile/:id` - Update user profile

### Why Centralized Documentation?

1. **Single Source of Truth**: All API documentation in one place
2. **Consistent Experience**: Uniform documentation across all services
3. **Easier Maintenance**: One documentation to update
4. **Better Developer Experience**: No need to navigate multiple Swagger UIs
5. **API Gateway Integration**: Documentation matches the actual API Gateway routes

## Microservices Architecture

### Services Overview

1. **Auth Service** (Port 8081)
   - User authentication and registration
   - JWT token management
   - Password hashing and validation

2. **User Service** (Port 8082)
   - User profile management
   - User data CRUD operations
   - User search and listing

3. **API Gateway** (Port 8080)
   - Single entry point for all requests
   - Request routing to appropriate services
   - Authentication middleware
   - CORS handling
   - **Centralized API Documentation**

4. **Frontend** (Port 4200)
   - Angular application
   - User interface for all operations
   - Authentication and protected routes

### Database
- **MongoDB**: Shared database across all services
- **Collections**: `users` (shared between auth and user services)

## Development Setup

### Backend Development

1. **Individual Service Development**:
   ```bash
   # Auth Service
   cd backend/auth-service
   go mod tidy
   go run main.go

   # User Service
   cd backend/user-service
   go mod tidy
   go run main.go

   # API Gateway
   cd backend/api-gateway
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

   # User Service
   PORT=8082
   MONGO_URI=mongodb://localhost:27017
   MONGO_DB=auth_db

   # API Gateway
   PORT=8080
   AUTH_SERVICE_URL=http://localhost:8081
   USER_SERVICE_URL=http://localhost:8082
   ```

### Frontend Development

```bash
cd frontend
npm install
ng serve
```

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
docker build -t auth-service backend/auth-service/
docker build -t user-service backend/user-service/
docker build -t api-gateway backend/api-gateway/

# Run individual services
docker run -p 8081:8081 auth-service
docker run -p 8082:8082 user-service
docker run -p 8080:8080 api-gateway
```

## Health Checks

```bash
# Check service health
curl http://localhost:8080/health  # API Gateway
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # User Service
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 8080, 8081, 8082, 4200, and 27017 are available
2. **MongoDB connection**: Check if MongoDB is running and accessible
3. **Service communication**: Verify service URLs in API Gateway configuration
4. **Authentication errors**: Check JWT secret configuration

### Logs and Debugging
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

# User Service
PORT=8082
MONGO_URI=mongodb://admin:password@mongodb:27017
MONGO_DB=auth_db

# API Gateway
PORT=8080
AUTH_SERVICE_URL=http://auth-service:8081
USER_SERVICE_URL=http://user-service:8082
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

# User Service
PORT=8082
MONGO_URI=mongodb://admin:password123@mongodb:27017
MONGO_DB=testdb

# API Gateway
PORT=8080
AUTH_SERVICE_URL=http://auth-service:8081
USER_SERVICE_URL=http://user-service:8082
```

## Features

- **Microservices Architecture**: Scalable and maintainable service decomposition
- **User Authentication**: JWT-based authentication system
- **User Management**: Complete CRUD operations for user profiles
- **API Gateway**: Single entry point with routing and middleware
- **MongoDB Integration**: Shared database across services
- **Docker Support**: Complete containerization for all services
- **Testing Environment**: Separate test configuration
- **Health Checks**: Service health monitoring
- **CORS Support**: Cross-origin resource sharing
- **Protected Routes**: Authentication-based route protection
- **Centralized API Documentation**: Swagger/OpenAPI documentation

## Future Enhancements

1. **Service Discovery**: Implement service discovery (Consul, etcd)
2. **Load Balancing**: Add load balancers for each service
3. **Circuit Breakers**: Implement circuit breakers for service communication
4. **Distributed Tracing**: Add tracing (Jaeger, Zipkin)
5. **Monitoring**: Implement metrics and monitoring (Prometheus, Grafana)
6. **Message Queues**: Add async communication (RabbitMQ, Kafka)
7. **API Documentation**: ✅ **COMPLETED** - Swagger/OpenAPI documentation
8. **Testing**: Add comprehensive test suites for each service

## Documentation

- [Backend README](backend/README.md) - Detailed backend documentation
- [Frontend README](frontend/README.md) - Frontend-specific instructions
