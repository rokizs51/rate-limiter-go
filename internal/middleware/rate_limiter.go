package middleware

import (
	"net/http"
	"rateLimiter/internal/config"
	"rateLimiter/internal/factory"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(cfg *config.Config, algorithm factory.Algorithm) gin.HandlerFunc {
	rateLimiter := factory.NewRateLimiter(algorithm, cfg)

	return func(c *gin.Context) {
		identifier := c.ClientIP()
		allowed, count, resetAt, err := rateLimiter.IsAllowed(identifier)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiting error",
			})
			return
		}
		//set rate limit headers
		switch algorithm {
		case factory.TokenBucket:
			c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.TokenBucketConfig.Tokens))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(count))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))
			c.Header("X-Rate-Limit-Type", "token-bucket")
		case factory.SlidingWindow:
			c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.SlidingWindowConfig.RequestLimit))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(cfg.SlidingWindowConfig.RequestLimit-count))
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))
			c.Header("X-Rate-Limit-Type", "sliding-window")
		}

		if !allowed {
			retryAfter := int(resetAt.Sub(time.Now()).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
