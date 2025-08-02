package config

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
}

// NewJWTConfig creates a new JWT configuration
func NewJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey: secretKey,
	}
}
