package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`

	jwt.RegisteredClaims
}

func JWTSecret() []byte {

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "development-secret-change-this"
	}

	return []byte(secret)
}

func AuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "authorization header is required",
				},
			)

			ctx.Abort()
			return
		}

		parts := strings.SplitN(
			authHeader,
			" ",
			2,
		)

		if len(parts) != 2 ||
			parts[0] != "Bearer" {

			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid authorization header format",
				},
			)

			ctx.Abort()
			return
		}

		tokenString := parts[1]

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return JWTSecret(), nil
			},
		)

		if err != nil || !token.Valid {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid or expired token",
				},
			)

			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		role, exists := ctx.Get("role")

		if !exists {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{
					"error": "user role not found",
				},
			)

			ctx.Abort()
			return
		}

		if role != requiredRole {
			ctx.JSON(
				http.StatusForbidden,
				gin.H{
					"error": "you do not have permission",
				},
			)

			ctx.Abort()
			return
		}

		ctx.Next()
	}
}