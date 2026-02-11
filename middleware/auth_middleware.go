package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-supabase-api/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"message": "Token tidak ada"})
			c.Abort()
			return
		}

		tokenStr := strings.Replace(auth, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_SECRET, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"message": "Token tidak valid"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("role", claims["role"])
		c.Next()
	}
}
