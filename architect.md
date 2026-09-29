# Project Architecture

This document describes the current high-level architecture and runtime flow for the application.

## High-level architecture

```mermaid
flowchart LR
    subgraph Client[Client Layer]
        FE[Angular Frontend\nlocalhost:8085]
        BROWSER[Browser Cookies\naccess_token JWT]
    end

    subgraph Backend[Backend Services]
        GW[API Gateway\nlocalhost:8080\nRoutes + CORS + Swagger]
        AS[Auth Service\nlocalhost:8081\nLogin / Register / Validate]
        US[User Service\nlocalhost:8082\nProfile CRUD]
        NS[Notification Service\nlocalhost:8083\nPreferences + Email Delivery]
        KAFKA[Kafka\nUser lifecycle events]
    end

    subgraph Data[Data Stores]
        AUTHDB[(MongoDB\nauth_db)]
        USERDB[(MongoDB\nuser_db)]
        NOTIFYDB[(MongoDB\nnotification_db)]
    end

    FE -->|HTTP + credentials| GW
    FE -->|Session cookie| BROWSER
    BROWSER -->|JWT access_token| GW
    BROWSER -->|JWT access_token| AS
    BROWSER -->|JWT access_token| US

    GW -->|/api/auth/*| AS
    GW -->|/api/users/*| US
    GW -->|/api/notifications/*| NS

    AS -->|Create / validate / refresh user session| AUTHDB
    AS -->|Publish user.created.v1| KAFKA
    AS -->|Publish user.updated.v1| KAFKA
    AS -->|Publish user.deleted.v1| KAFKA

    KAFKA -->|Consume user lifecycle events| US
    US -->|Read / write user profiles| USERDB
    KAFKA -->|user.created.v1| NS
    NS -->|Persist preferences and delivery history| NOTIFYDB
    NS -->|SMTP welcome email| EMAIL[Email provider]

    US -->|Return profile data| FE
    AS -->|Return auth result / cookie| FE
    NS -->|Return preference and history| FE
```

## Request and authentication flow

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Gateway
    participant AuthService
    participant UserService
    participant MongoDB

    User->>Frontend: Enter credentials
    Frontend->>Gateway: POST /api/auth/login
    Gateway->>AuthService: Forward login request
    AuthService->>MongoDB: Validate email and password
    MongoDB-->>AuthService: User record
    AuthService-->>Gateway: Success + JWT
    AuthService-->>Frontend: Set HttpOnly access_token cookie
    Frontend-->>User: Authenticated session

    User->>Frontend: Request protected page
    Frontend->>Gateway: GET /api/users/me (cookie included)
    Gateway->>AuthService: Validate token
    AuthService-->>Gateway: Valid session
    Gateway->>UserService: Forward request with auth context
    UserService-->>Frontend: User profile data
```

## Event-driven synchronization flow

```mermaid
flowchart TD
    AUTH[Auth Service]
    KAFKA[Kafka Topic: user.events]
    USER[User Service]
    NOTIFY[Notification Service]
    MONGO[(MongoDB user_db.user_profiles)]
    NOTE[(MongoDB notification_db)]

    AUTH -->|user.created.v1| KAFKA
    AUTH -->|user.updated.v1| KAFKA
    AUTH -->|user.deleted.v1| KAFKA

    KAFKA -->|Consume event| USER
    USER -->|Upsert / Delete profile| MONGO

    KAFKA -->|user.created.v1| NOTIFY
    NOTIFY -->|Persist preference + delivery history| NOTE
```

## Components

### API Gateway

- Routes requests to the auth, user, and notification services
- Centralizes Swagger and OpenAPI docs
- Handles CORS and credentials-aware browser traffic
- Forwards cookie-based authentication context to downstream services

### Auth Service

- Handles registration, login, logout, and refresh
- Validates credentials against MongoDB
- Creates and validates JWTs
- Emits user lifecycle events to Kafka

### User Service

- Owns user profile state in `user_db`
- Consumes Kafka events to keep profiles synchronized
- Enforces that users only access their own profile data

### Notification Service

- Stores notification preferences and delivery history
- Consumes `user.created.v1` and sends a welcome email when enabled
- Records failed or successful delivery attempts
- Uses optional SMTP configuration and degrades gracefully when email is not configured

## Data model

- `auth_db.auth_users`: authentication records
- `user_db.user_profiles`: synced profile data
- `notification_db.notification_preferences`: per-user email preference
- `notification_db.notification_records`: delivery history and outcomes

## Operational notes

- Browser sessions use HttpOnly cookies; no token is stored in localStorage
- Kafka keeps the services loosely coupled and resilient to async propagation
- The gateway acts as the single public interface for client traffic
- The backend supports both local containerized runs and production deployment via Docker Compose and the `deploy.sh` helper
