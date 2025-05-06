# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install swag and make
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    apk add --no-cache make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate Swagger documentation
RUN swag init -g main.go -o docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy Swagger documentation
COPY --from=builder /app/docs ./docs

COPY .env . 

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 