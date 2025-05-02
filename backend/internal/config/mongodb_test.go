package config

import (
	"os"
	"testing"
)

func TestNewMongoDBConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("MONGODB_USERNAME", "admin")
	os.Setenv("MONGODB_PASSWORD", "password123")
	os.Setenv("MONGODB_HOST", "localhost")
	os.Setenv("MONGODB_PORT", "27017")
	os.Setenv("MONGODB_DATABASE", "testdb")

	// Create config
	config, err := NewMongoDBConfig("mongodb://admin:password123@localhost:27017", "testdb")
	if err != nil {
		t.Fatalf("Failed to create MongoDB config: %v", err)
	}

	// Verify URI construction
	expectedURI := "mongodb://admin:password123@localhost:27017"
	if config.URI != expectedURI {
		t.Errorf("Expected URI %s, got %s", expectedURI, config.URI)
	}

	// Verify database name
	if config.Database != "testdb" {
		t.Errorf("Expected database testdb, got %s", config.Database)
	}

	// Clean up
	os.Unsetenv("MONGODB_USERNAME")
	os.Unsetenv("MONGODB_PASSWORD")
	os.Unsetenv("MONGODB_HOST")
	os.Unsetenv("MONGODB_PORT")
	os.Unsetenv("MONGODB_DATABASE")
}
