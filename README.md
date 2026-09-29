# Full Stack Microservices Project

A full-stack application built with an Angular frontend and a Go-based microservice backend. The project uses an API gateway, cookie-based JWT authentication, MongoDB per service, and Kafka-based event synchronization.

## Architecture at a glance

- Frontend: Angular app served on port 8085
- API Gateway: single entry point on port 8080
- Auth Service: login/register/session validation on port 8081
- User Service: profile CRUD and authorization checks on port 8082
- Notification Service: welcome email + preference/history APIs on port 8083
- Kafka: user lifecycle event bus
- MongoDB: service-specific databases

## Services and responsibilities

| Service | Port | Responsibility |
| --- | --- | --- |
| Frontend | 8085 | Angular UI, protected routes, credentialed browser requests |
| API Gateway | 8080 | Routing, authentication, CORS, Swagger docs |
| Auth Service | 8081 | Registration, login, logout, session validation, JWT cookie issuance |
| User Service | 8082 | User profile CRUD and profile synchronization |
| Notification Service | 8083 | Email preference management, delivery history, welcome email processing |
| MongoDB | 27017 | Persistent storage for auth, user, and notification data |
| Kafka | 9092 | User-created/updated/deleted events |

## Repository layout

```text
.
├── backend/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── notification-service/
│   ├── user-service/
│   ├── README.md
│   ├── architect.md
│   ├── Makefile
│   └── .trivyignore
├── deploy/
│   ├── docker-compose.yml
│   ├── docker-compose.prod.yml
│   └── docker-compose.test.yml
├── frontend/
├── .githooks/
├── deploy.sh
├── README.md
├── architect.md
├── test.sh
└── .env
```

## Quick start

### Prerequisites

- Docker + Docker Compose
- Node.js for frontend work
- Go 1.26+ for backend work

### Local development stack

Start the full environment:

```bash
docker compose -f deploy/docker-compose.yml up -d --build
```

Access points:

- Frontend: http://localhost:8085
- API Gateway: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/
- Auth Service: http://localhost:8081
- User Service: http://localhost:8082
- Notification Service: http://localhost:8083
- MongoDB: localhost:27017
- Kafka: localhost:9092

Stop everything:

```bash
docker compose -f deploy/docker-compose.yml down
```

### Production deployment

Use the provided script to generate secrets and launch the production stack:

```bash
./deploy.sh
```

You can also run the production compose file directly if needed:

```bash
docker compose -f deploy/docker-compose.prod.yml --env-file .env up -d --build
```

## Authentication model

The application uses HttpOnly cookie-based sessions for browser clients.

- Cookie name: `access_token`
- JWT is stored in the cookie and sent automatically by the browser
- Frontend does not store auth tokens in localStorage
- Auth validation supports cookie, request body, or Authorization bearer tokens depending on the request flow
- `JWT_SECRET` must match across auth-service and user-service
- CORS is configured with explicit allowed origins and credentials enabled

## API docs and gateway behavior

The API Gateway is the single public entry point for the frontend and exposes the central Swagger/OpenAPI documentation.

OpenAPI entry points:

- http://localhost:8080/swagger/
- http://localhost:8080/swagger/index.html
- http://localhost:8080/swagger/doc.json
- http://localhost:8080/swagger/doc.yaml

Key routes:

- `POST /api/auth/login`
- `POST /api/auth/register`
- `GET /api/auth/me`
- `POST /api/auth/logout`
- `POST /api/auth/validate`
- `POST /api/auth/refresh`
- `GET /api/users/me`
- `GET /api/users/list`
- `GET /api/users/profile/:id`
- `PUT /api/users/profile/:id`
- `DELETE /api/users/profile/:id`
- `GET /api/notifications/preferences`
- `PUT /api/notifications/preferences`
- `GET /api/notifications/history`

## Event-driven synchronization

The platform uses Kafka for asynchronous user lifecycle tracking.

- `user.created.v1`
- `user.updated.v1`
- `user.deleted.v1`

The auth service emits lifecycle events after user changes, and the user service consumes them to keep profile data consistent. The notification service consumes `user.created.v1` to send the welcome email when the user has email notifications enabled.

## Development workflow

Install repository hooks:

```bash
./.githooks/install.sh
```

Run backend formatting and tests for all Go services:

```bash
make -C backend lint-format
```

Run vulnerability checks before pushing:

```bash
make -C backend security-check
```

## Additional documentation

- [backend/README.md](backend/README.md)
- [backend/architect.md](backend/architect.md)
- [architect.md](architect.md)

## Notes

- MongoDB is configured as a database-per-service model.
- Notification records expire after 90 days.
- SMTP settings are optional; failed deliveries are logged as failed notifications instead of breaking the app.
- The API Gateway centralizes documentation and request routing for the frontend.
