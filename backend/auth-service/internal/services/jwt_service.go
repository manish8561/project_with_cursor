package services

import (
	"auth-service/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token generation and validation using a secret key.
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWTService with the provided JWT configuration.
func NewJWTService(jwtConfig *config.JWTConfig) *JWTService {
	return &JWTService{
		secretKey: []byte(jwtConfig.SecretKey),
	}
}

// Claims defines the custom and registered claims for JWT tokens.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a signed JWT token string for the given user ID.
// The token is valid for 24 hours from the time of issuance.
func (s *JWTService) GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken parses and validates the provided JWT token string.
// It returns the Claims if the token is valid, or an error otherwise.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
