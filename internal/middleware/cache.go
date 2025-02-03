package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func CacheMiddleware(redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check cache
		key := c.Request.URL.Path
		val, err := redis.Get(c.Request.Context(), key).Result()
		if err == nil {
			c.Header("X-Cache", "HIT")
			c.Data(http.StatusOK, "application/json", []byte(val))
			c.Abort()
			return
		}

		// Continue if not found in cache
		c.Next()
	}
}
