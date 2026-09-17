package services

import (
	"auth-service/internal/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testJWTSecret = "test-secret-that-is-not-used-in-production"

func newTestJWTService() *JWTService {
	return NewJWTService(config.NewJWTConfig(testJWTSecret))
}

func TestJWTServiceGenerateAndValidateToken(t *testing.T) {
	service := newTestJWTService()

	token, err := service.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned an empty token")
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("ValidateToken() user ID = %q, want %q", claims.UserID, "user-123")
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(time.Now()) {
		t.Error("ValidateToken() returned a token that is not expired in the future")
	}
}

func TestJWTServiceValidateTokenRejectsInvalidTokens(t *testing.T) {
	service := newTestJWTService()

	expiredToken := signedTestToken(t, testJWTSecret, "user-123", time.Now().Add(-time.Hour))
	wrongSecretToken := signedTestToken(t, "another-secret", "user-123", time.Now().Add(time.Hour))

	for _, tt := range []struct {
		name  string
		token string
	}{
		{name: "malformed", token: "not-a-jwt"},
		{name: "expired", token: expiredToken},
		{name: "wrong signature", token: wrongSecretToken},
	} {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)
			if err == nil {
				t.Fatal("ValidateToken() error = nil, want an error")
			}
			if claims != nil {
				t.Errorf("ValidateToken() claims = %#v, want nil", claims)
			}
		})
	}
}

func TestAuthServiceValidateToken(t *testing.T) {
	jwtService := newTestJWTService()
	authService := NewAuthService(nil, jwtService, nil)

	validToken, err := jwtService.GenerateToken("user-456")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	response, err := authService.ValidateToken(validToken)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if !response.Valid || response.UserID != "user-456" || response.Message != "" {
		t.Errorf("ValidateToken() response = %#v, want valid response for user-456", response)
	}

	response, err = authService.ValidateToken("invalid-token")
	if err != nil {
		t.Fatalf("ValidateToken() invalid token error = %v, want nil", err)
	}
	if response.Valid || response.Message != "Invalid token" {
		t.Errorf("ValidateToken() invalid response = %#v, want invalid token response", response)
	}
}

func TestAuthServiceRefreshToken(t *testing.T) {
	jwtService := newTestJWTService()
	authService := NewAuthService(nil, jwtService, nil)

	token, err := jwtService.GenerateToken("user-789")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	response, err := authService.RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken() error = %v", err)
	}
	if response.Status != "success" || response.Token == "" {
		t.Errorf("RefreshToken() response = %#v, want successful response with token", response)
	}

	claims, err := jwtService.ValidateToken(response.Token)
	if err != nil {
		t.Fatalf("refreshed token validation error = %v", err)
	}
	if claims.UserID != "user-789" {
		t.Errorf("refreshed token user ID = %q, want %q", claims.UserID, "user-789")
	}

	if _, err := authService.RefreshToken("invalid-token"); err == nil {
		t.Error("RefreshToken() invalid token error = nil, want an error")
	}
}

func signedTestToken(t *testing.T, secret, userID string, expiresAt time.Time) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}
	return signed
}
