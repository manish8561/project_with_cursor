# Backend Server

A Go backend server using the Gin framework with user authentication functionality.

## Project Structure

```
backend/
├── internal/
│   ├── handlers/          # HTTP request handlers
│   │   └── user_handler.go
│   ├── models/           # Data structures
│   │   └── user.go
│   └── services/         # Business logic
│       └── user_service.go
├── docs/                 # Swagger API documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── main.go              # Server entry point
├── go.mod              # Go module definition
├── generate_swagger.sh # Swagger documentation generator
└── .gitignore         # Git ignore rules
```

## Features

- RESTful API endpoints
- User authentication
- CORS middleware
- Input validation
- Health check endpoint
- Swagger API documentation

## Prerequisites

- Go 1.21 or higher
- Git

## Setup

1. Clone the repository
2. Navigate to the backend directory:
   ```bash
   cd backend
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Server

Start the server:
```bash
go run main.go
```

The server will start on port 8080.

## API Documentation (Swagger)

### Generate Swagger Documentation

To regenerate the Swagger API documentation after making changes to the API:

```bash
# Make the script executable (first time only)
chmod +x generate_swagger.sh

# Generate the documentation
./generate_swagger.sh
```

This will:
- Install the `swag` tool if not already installed
- Generate updated `docs.go`, `swagger.json`, and `swagger.yaml` files
- Update the API documentation based on the current code

### View API Documentation

Once the server is running, you can view the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:
- Interactive API testing
- Request/response examples
- Schema definitions
- Authentication information

### Manual Swagger Generation

If you prefer to generate Swagger docs manually:

```bash
# Install swag tool
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g main.go -o docs
```

## API Endpoints

### Health Check
- **GET** `/api/health`
  - Returns server status
  - Response:
    ```json
    {
      "status": "ok"
    }
    ```

### User Authentication
- **POST** `/api/user/login`
  - Authenticates a user
  - Request body:
    ```json
    {
      "email": "test@example.com",
      "password": "password123"
    }
    ```
  - Success response:
    ```json
    {
      "status": "success",
      "token": "jwt-token-here"
    }
    ```

- **POST** `/api/user/register`
  - Registers a new user
  - Request body:
    ```json
    {
      "name": "John Doe",
      "email": "john@example.com",
      "password": "password123",
      "confirmPassword": "password123"
    }
    ```
  - Success response:
    ```json
    {
      "status": "success",
      "message": "User registered successfully"
    }
    ```

- **GET** `/api/user/profile`
  - Get user profile (requires authentication)
  - Headers: `Authorization: Bearer <token>`
  - Success response:
    ```json
    {
      "id": "user-id",
      "name": "John Doe",
      "email": "john@example.com"
    }
    ```

## Testing the API

You can test the endpoints using curl:

1. Health check:
   ```bash
   curl http://localhost:8080/api/health
   ```

2. Login:
   ```bash
   curl -X POST http://localhost:8080/api/user/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

3. Register:
   ```bash
   curl -X POST http://localhost:8080/api/user/register \
     -H "Content-Type: application/json" \
     -d '{"name":"John Doe","email":"john@example.com","password":"password123","confirmPassword":"password123"}'
   ```

## Development

### Current Implementation
- MongoDB database integration
- JWT token authentication
- Password validation
- User registration and login
- Protected routes
- Swagger API documentation

### Future Improvements
1. Password hashing (bcrypt)
2. Email verification
3. Password reset functionality
4. Role-based access control
5. Rate limiting
6. Comprehensive unit tests
7. API versioning

## Security Notes

⚠️ **Important**: This is a development version. For production use:
- Implement proper password hashing with bcrypt
- Use secure JWT token generation with proper expiration
- Add rate limiting and request throttling
- Implement proper database connection pooling
- Add input sanitization and validation
- Use environment variables for sensitive data
- Enable HTTPS
- Add request logging and monitoring

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License. 