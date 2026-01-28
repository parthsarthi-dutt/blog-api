package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/parthsarthi-dutt/blog-api/internal/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "AUTH_MISSING", "Missing authorization header")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, http.StatusUnauthorized, "AUTH_INVALID", "Invalid authorization header")
			c.Abort()
			return
		}

		tokenStr := parts[1]

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			utils.Error(c, http.StatusInternalServerError, "SERVER_ERROR", "JWT secret not configured")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// ✅ Ensure correct signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			utils.Error(c, http.StatusUnauthorized, "AUTH_INVALID", "Invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.Error(c, http.StatusUnauthorized, "AUTH_INVALID", "Invalid token claims")
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			utils.Error(c, http.StatusUnauthorized, "AUTH_INVALID", "Invalid token payload")
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Next()
	}
}
