FROM golang:1.23 AS builder

WORKDIR /src

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-gateway ./cmd/api-gateway

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl


# Copy the built binary from builder
COPY --from=builder /app/api-gateway /app/
# Copy Swagger UI static files and OpenAPI spec
COPY --from=builder /src/swagger-ui /app/swagger-ui
# Ensure the binary is executable
RUN chmod +x /app/api-gateway

WORKDIR /app

EXPOSE 8080
EXPOSE 9000
VOLUME /data/conf

CMD ["./api-gateway", "-conf", "/data/conf"]
