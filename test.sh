#!/bin/bash

# Stop any existing containers
docker compose -f deploy/docker-compose.test.yml down

# Start MongoDB for testing
docker compose -f deploy/docker-compose.test.yml up -d mongodb

# Wait for MongoDB to be ready
echo "Waiting for MongoDB to be ready..."
sleep 10

# Run the tests
cd backend
go test ./... -v

# Clean up
cd ..
docker compose -f deploy/docker-compose.test.yml down 