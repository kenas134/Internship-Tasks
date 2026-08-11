package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// --------------------------------------------------
// AUTHENTICATION MIDDLEWARE
// --------------------------------------------------

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Get Authorization header.
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})

			c.Abort()
			return
		}

		// Expected format:
		//
		// Bearer <token>
		//
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			strings.ToLower(parts[0]) != "bearer" {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header",
			})

			c.Abort()
			return
		}

		tokenString := parts[1]

		// Create claims object.
		claims := &Claims{}

		// Parse token.
		token, err := jwt.ParseWithClaims(
			tokenString,
			claims,
			func(token *jwt.Token) (interface{}, error) {

				// Make sure we use HMAC.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		// Validate token.
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired access token",
			})

			c.Abort()
			return
		}

		// Put user information into Gin context.
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// Authentication successful.
		c.Next()
	}
}

// --------------------------------------------------
// AUTHORIZATION MIDDLEWARE
// --------------------------------------------------

func RequireRole(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		role, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Role information not found",
			})

			c.Abort()
			return
		}

		userRole, ok := role.(string)

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid role",
			})

			c.Abort()
			return
		}

		// Check allowed roles.
		for _, allowedRole := range allowedRoles {

			if userRole == allowedRole {

				c.Next()
				return
			}
		}

		// Authenticated but not authorized.
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})

		c.Abort()
	}
}