package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyMiddleware(validAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validAPIKey == "" {
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing X-API-Key header"})
			c.Abort()
			return
		}

		if apiKey != validAPIKey {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid X-API-Key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
