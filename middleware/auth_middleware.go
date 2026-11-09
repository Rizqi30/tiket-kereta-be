package middleware

import (
	"net/http"
	"os"
	"strings"

	"tiket-kereta-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth Middleware - Mengecek user sudah login (token JWT valid)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// take the header Authorizarion
		authHeader := c.GetHeader("Authorization")
		
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not found, please login first"})
			c.Abort()
			return
		}

		// split string "Bearer <token>" to just "<token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token is invalid"})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		// take secret key from .env
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "default_secret_key_if_env_empty"
		}

		// parse and validated token
		claims := &utils.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalid or expired"})
			c.Abort()
			return 
		}

		// save user data from token to context Gin
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// forward to next handler/controller
		c.Next()
	}
}

// AdminMiddleware - Mengecek apakah user yang login memiliki role ADMIN
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied, just ADMIN who perrmitted"})
			c.Abort()
			return 
		}
		c.Next()
	}
}