package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigHelpers(t *testing.T) {
	assert.Equal(t, "fallback", getEnv("USER_SERVICE_TEST_UNSET", "fallback"))
	t.Setenv("USER_SERVICE_TEST_VALUE", "configured")
	assert.Equal(t, "configured", getEnv("USER_SERVICE_TEST_VALUE", "fallback"))
	assert.Equal(t, []string{"one", "two"}, splitCSV(" one, ,two "))
}
