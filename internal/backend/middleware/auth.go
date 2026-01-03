package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/philopaterwaheed/passGO/internal/backend/auth"
)

// AuthMiddleware validates JWT tokens and sets user information in context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check for Bearer token format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Verify the token
		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			switch err {
			case auth.ErrTokenExpired:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			case auth.ErrInvalidToken:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			case auth.ErrTokenClaims:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			}
			c.Abort()
			return
		}

		// Set user information in context for handlers to use
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("supabaseUID", claims.SupabaseUID)
		c.Set("claims", claims)

		c.Next()
	}
}

