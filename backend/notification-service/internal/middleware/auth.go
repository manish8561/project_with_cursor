package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "user_id"

func CookieAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenValue, _ := c.Cookie("access_token")
		if tokenValue == "" {
			tokenValue = bearerToken(c.GetHeader("Authorization"))
		}
		if tokenValue == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization token is required"})
			return
		}
		userID, err := parseUserID(tokenValue, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func bearerToken(header string) string {
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}

func parseUserID(tokenValue, secret string) (string, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}
	userID, _ := claims["user_id"].(string)
	if userID == "" {
		return "", fmt.Errorf("invalid user ID in token")
	}
	return userID, nil
}
