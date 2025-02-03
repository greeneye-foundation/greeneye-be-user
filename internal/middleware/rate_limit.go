package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/errors"
)

type RateLimiter struct {
	redisClient *redis.Client
	limit       int
	window      time.Duration
}

func RateLimit(redisClient *redis.Client) gin.HandlerFunc {
	// Define rate limit parameters
	limit := 100              // max requests
	window := 1 * time.Minute // per minute

	rl := &RateLimiter{
		redisClient: redisClient,
		limit:       limit,
		window:      window,
	}

	return rl.HandleRateLimit
}

func (rl *RateLimiter) HandleRateLimit(c *gin.Context) {
	clientIP := c.ClientIP()
	key := "rate_limit:" + clientIP

	// Increment the count for this IP
	count, err := rl.redisClient.Incr(c.Request.Context(), key).Result()
	if err != nil {
		// If Redis fails, allow the request but log the error
		c.Next()
		return
	}

	if count == 1 {
		// Set expiration
		rl.redisClient.Expire(c.Request.Context(), key, rl.window)
	}

	if count > int64(rl.limit) {
		// Exceeded the limit
		c.JSON(http.StatusTooManyRequests, errors.New(
			http.StatusTooManyRequests,
			"Too Many Requests",
			"Rate limit exceeded. Try again later.",
		))
		c.Abort()
		return
	}

	c.Next()
}
