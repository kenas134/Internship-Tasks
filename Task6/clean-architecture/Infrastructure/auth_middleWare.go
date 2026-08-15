package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware(jwtService JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {



		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorizationnheader is required",
			})
			ctx.Abort()
			return
		}

		// 2. Expected format:
		// Authorization: Bearer <token>
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		userID,role, err := jwtService.ValidateAccessToken(tokenString)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired access token",
			})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("role",role)
		ctx.Next()
	}
}

func RequiredRole(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		role, exists := ctx.Get("role")

		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "user role not found",
			})
			ctx.Abort()
			return
		}

		userRole, ok := role.(string)

		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user role",
			})
			ctx.Abort()
			return
		}

		if userRole != requiredRole {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "you do not have permission",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
