# Full Stack Project

A full-stack application with Angular frontend and Go backend.

## Project Structure

```
.
├── backend/           # Go backend service
├── frontend/          # Angular frontend
└── deploy/           # Deployment configurations
    └── docker-compose.yml
```

## Prerequisites

- Docker
- Docker Compose
- Node.js (for frontend development)
- Go (for backend development)

## Development Setup

### Backend
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend
```bash
cd frontend
npm install
ng serve
```

## Docker Deployment

### Using Docker Compose

1. Build and start all services:
```bash
cd deploy
docker-compose up --build
```

2. Access the services:
- Frontend: http://localhost
- Backend API: http://localhost:8080

### Individual Services

#### Backend
```bash
cd backend
docker build -t backend .
docker run -p 8080:8080 backend
```

#### Frontend
```bash
cd frontend
docker build -t frontend .
docker run -p 80:80 frontend
```

## API Endpoints

### Backend
- Health Check: `GET /health`
- Login: `POST /api/login`
- Register: `POST /api/register`

### Frontend
- Login Page: `/login`
- Dashboard: `/dashboard`

## Environment Variables

### Backend
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release)

### Frontend
- Configured through nginx.conf

## Development Notes

- Backend runs in debug mode by default
- Frontend uses Angular's development server
- Docker setup includes health checks and automatic restarts
- Services are connected through a Docker network

## Production Considerations

- Set `GIN_MODE=release` for production
- Use proper environment variables
- Implement proper security measures
- Configure proper logging
- Set up monitoring
