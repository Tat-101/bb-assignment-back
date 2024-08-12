package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tat-101/bb-assignment-back/domain"
	"github.com/tat-101/bb-assignment-back/internal/rest/service"
)

func AuthMiddleware(svc service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// TODO: improve cache
		user, err := svc.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		userStruct, ok := user.(*domain.User)
		if !ok || userStruct.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied, admin role required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
