#!/bin/bash
set -euo pipefail

compose_file="deploy/docker-compose.test.yml"

cleanup() {
  docker compose -f "$compose_file" down
}

trap cleanup EXIT

# Stop any existing containers
cleanup

# Start MongoDB for testing
docker compose -f "$compose_file" up -d mongodb

# Wait for MongoDB to be ready
echo "Waiting for MongoDB to be ready..."
sleep 10

# Run all backend service tests
make -C backend test
