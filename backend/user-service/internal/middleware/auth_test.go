package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func signedToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	assert.NoError(t, err)
	return token
}

func authRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", CookieAuth(secret), func(c *gin.Context) {
		userID, _ := c.Get(ContextUserIDKey)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})
	return router
}

func TestCookieAuthAcceptsBearerToken(t *testing.T) {
	const secret = "test-secret"
	router := authRouter(secret)
	token := signedToken(t, secret, jwt.MapClaims{"user_id": "user-123", "exp": time.Now().Add(time.Hour).Unix()})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `{"user_id":"user-123"}`, response.Body.String())
}

func TestCookieAuthCookieTakesPrecedence(t *testing.T) {
	const secret = "test-secret"
	router := authRouter(secret)
	cookieToken := signedToken(t, secret, jwt.MapClaims{"user_id": "cookie-user", "exp": time.Now().Add(time.Hour).Unix()})
	headerToken := signedToken(t, secret, jwt.MapClaims{"user_id": "header-user", "exp": time.Now().Add(time.Hour).Unix()})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: AccessTokenCookieName, Value: cookieToken})
	req.Header.Set("Authorization", "Bearer "+headerToken)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `{"user_id":"cookie-user"}`, response.Body.String())
}

func TestCookieAuthRejectsInvalidRequests(t *testing.T) {
	const secret = "test-secret"
	router := authRouter(secret)

	tests := []struct {
		name  string
		token string
		body  string
	}{
		{name: "missing token", body: `{"error":"authorization token is required"}`},
		{name: "invalid signature", token: signedToken(t, "other-secret", jwt.MapClaims{"user_id": "user-123", "exp": time.Now().Add(time.Hour).Unix()}), body: `{"error":"invalid token"}`},
		{name: "missing user id", token: signedToken(t, secret, jwt.MapClaims{"exp": time.Now().Add(time.Hour).Unix()}), body: `{"error":"invalid user ID in token"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)

			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.JSONEq(t, tt.body, response.Body.String())
		})
	}
}
