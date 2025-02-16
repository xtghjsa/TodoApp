package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strings"
)

// AuthenticationMiddleware middleware authentication
func AuthenticationMiddleware(rc *redis.Client, ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "need authorization header"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "bearer token required"})
			c.Abort()
			return
		}
		userID, err := rc.Get(ctx, token).Result()
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			c.Abort()
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't check authorized session"})
			c.Abort()
			return
		}
		log.Printf("Authenticated user with ID: %v\n", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
