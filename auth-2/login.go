package main
import (
    "github.com/golang-jwt/jwt/v5"
    "time"
    "os"
)

var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-256-bit-secret-change-in-production"))

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func generateToken(user *User) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)

    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{expirationTime},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "go-auth-service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString(jwtSecret)
}

