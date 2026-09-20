# Backend Architecture Diagram

This document defines the complete backend flow for the application, including the client, API gateway, auth service, user service, MongoDB storage, and Kafka event sync.

## High-level architecture

```mermaid
flowchart LR
    subgraph Client[Client Layer]
        FE[Angular Frontend\nlocalhost:8085]
        BROWSER[Browser Cookies\naccess_token JWT]
    end

    subgraph Backend[Backend Services]
        GW[API Gateway\nlocalhost:8080\nRoutes + CORS]
        AS[Auth Service\nlocalhost:8081\nLogin / Register / Validate / Refresh]
        US[User Service\nlocalhost:8082\nProfile CRUD + Auth Middleware]
        KAFKA[Kafka\nUser lifecycle events]
    end

    subgraph Data[Data Stores]
        AUTHDB[(MongoDB\nauth_db)]
        USERDB[(MongoDB\nuser_db)]
    end

    FE -->|HTTP + credentials| GW
    FE -->|Session cookie| BROWSER
    BROWSER -->|JWT access_token| GW
    BROWSER -->|JWT access_token| AS
    BROWSER -->|JWT access_token| US

    GW -->|/api/auth/*| AS
    GW -->|/api/users/*| US

    AS -->|Create / validate / refresh user session| AUTHDB
    AS -->|Publish user.created.v1| KAFKA
    AS -->|Publish user.updated.v1| KAFKA
    AS -->|Publish user.deleted.v1| KAFKA

    KAFKA -->|Consume user lifecycle events| US
    US -->|Read / write user profiles| USERDB

    US -->|Return profile data| FE
    AS -->|Return auth result / cookie| FE
```

## Authentication flow

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Gateway
    participant AuthService
    participant MongoDB

    User->>Frontend: Enter credentials
    Frontend->>Gateway: POST /api/auth/login
    Gateway->>AuthService: Forward login request
    AuthService->>MongoDB: Validate email/password
    MongoDB-->>AuthService: User record
    AuthService-->>Gateway: JWT token + success response
    AuthService-->>Frontend: Set HttpOnly access_token cookie
    Frontend-->>User: Authenticated session

    User->>Frontend: Request protected page
    Frontend->>Gateway: GET /api/users/me (cookie included)
    Gateway->>AuthService: Validate token
    AuthService-->>Gateway: Valid user session
    Gateway->>UserService: Forward request with auth context
    UserService-->>Frontend: User profile data
```

## User profile synchronization flow

```mermaid
flowchart TD
    AUTH[Auth Service]
    KAFKA[Kafka Topic: user.events]
    USER[User Service]
    MONGO[(MongoDB user_db.user_profiles)]

    AUTH -->|user.created.v1| KAFKA
    AUTH -->|user.updated.v1| KAFKA
    AUTH -->|user.deleted.v1| KAFKA

    KAFKA -->|Consume event| USER
    USER -->|Upsert / Delete profile| MONGO

    USER -->|Read profile| MONGO
```

## Request lifecycle summary

1. The frontend sends requests to the API gateway.
2. The gateway routes the request to the relevant microservice.
3. Auth service handles login, registration, token validation, refresh, and logout.
4. User service handles profile retrieval, updates, listing, and deletion.
5. JWT is stored in an HttpOnly cookie named `access_token`.
6. Events are published to Kafka when user lifecycle changes occur.
7. User service consumes those events to keep the profile database synchronized.

## Notes

- `auth-service` owns `auth_db` authentication data.
- `user-service` owns `user_db` profile data.
- The API gateway acts as the single entry point for frontend traffic.
- Cookie-based auth is used for browser sessions, while JWT validation also supports Bearer token fallback in some flows.
