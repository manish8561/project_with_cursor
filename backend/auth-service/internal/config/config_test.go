package config

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCookieConfig(t *testing.T) {
	tests := []struct {
		name     string
		sameSite string
		expected http.SameSite
	}{
		{name: "strict", sameSite: "strict", expected: http.SameSiteStrictMode},
		{name: "none", sameSite: "NONE", expected: http.SameSiteNoneMode},
		{name: "default", sameSite: "unknown", expected: http.SameSiteLaxMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := NewCookieConfig(true, 0, tt.sameSite, "example.test")
			assert.Equal(t, AccessTokenCookieName, cookie.Name)
			assert.Equal(t, 24*60*60, cookie.MaxAge)
			assert.True(t, cookie.Secure)
			assert.True(t, cookie.HTTPOnly)
			assert.Equal(t, "/", cookie.Path)
			assert.Equal(t, "example.test", cookie.Domain)
			assert.Equal(t, tt.expected, cookie.SameSite)
		})
	}
}

func TestConfigHelpers(t *testing.T) {
	assert.True(t, parseBoolEnv("true", false))
	assert.False(t, parseBoolEnv("invalid", false))
	assert.Equal(t, 42, parseIntEnv("42", 1))
	assert.Equal(t, 1, parseIntEnv("invalid", 1))
	assert.Equal(t, []string{"https://one.test", "https://two.test"}, splitCSV(" https://one.test, ,https://two.test "))
}
