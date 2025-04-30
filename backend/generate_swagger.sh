#!/bin/bash

# Install swag if not already installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g main.go -o docs

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo "Swagger documentation generated successfully!"
    echo "You can view the documentation at: http://localhost:8080/swagger/index.html"
else
    echo "Failed to generate Swagger documentation"
    exit 1
fi 