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
├── main.go              # Server entry point
├── go.mod              # Go module definition
└── .gitignore         # Git ignore rules
```

## Features

- RESTful API endpoints
- User authentication
- CORS middleware
- Input validation
- Health check endpoint

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

## API Endpoints

### Health Check
- **GET** `/health`
  - Returns server status
  - Response:
    ```json
    {
      "status": "ok"
    }
    ```

### User Authentication
- **POST** `/api/login`
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
      "token": "dummy-token",
      "user": {
        "id": 1,
        "email": "test@example.com"
      }
    }
    ```
  - Error response:
    ```json
    {
      "error": "Invalid credentials"
    }
    ```

## Testing the API

You can test the endpoints using curl:

1. Health check:
   ```bash
   curl http://localhost:8080/health
   ```

2. Login:
   ```bash
   curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

## Development

### Current Implementation
- Uses in-memory user storage
- Basic authentication without password hashing
- Dummy token generation

### Future Improvements
1. JWT token generation
2. Password hashing
3. Database integration
4. User registration
5. Input validation middleware
6. Unit tests
7. API documentation

## Security Notes

⚠️ **Important**: This is a development version. For production use:
- Implement proper password hashing
- Use secure JWT token generation
- Add rate limiting
- Implement proper database storage
- Add input sanitization
- Use environment variables for sensitive data

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License. 