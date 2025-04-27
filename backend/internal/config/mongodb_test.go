package config

import (
	"os"
	"testing"
)

func TestNewMongoDBConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("MONGODB_USERNAME", "testuser")
	os.Setenv("MONGODB_PASSWORD", "testpass")
	os.Setenv("MONGODB_HOST", "testhost")
	os.Setenv("MONGODB_PORT", "27018")
	os.Setenv("MONGODB_DATABASE", "testdb")

	// Create config
	config := NewMongoDBConfig()

	// Verify URI construction
	expectedURI := "mongodb://testuser:testpass@testhost:27018"
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
