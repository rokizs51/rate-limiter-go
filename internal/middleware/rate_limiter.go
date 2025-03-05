package middleware

import (
	"fmt"
	"log"
	"net/http"
	"rateLimiter/internal/config"
	"rateLimiter/internal/factory"
	service "rateLimiter/internal/service/rate_limit"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(cfg *config.Config, algorithm factory.Algorithm) gin.HandlerFunc {
	rateLimiter := service.NewRedisTokenBucketService(&cfg.TokenBucketConfig)
	fmt.Println(rateLimiter)
	return func(c *gin.Context) {
		// Add fallback mechanism
		identifier := c.ClientIP()
		allowed, count, resetAt, err := rateLimiter.IsAllowed(c.Request.Context(), identifier)
		if err != nil {
			// If Redis is not available, fallback to allow traffic
			log.Printf("Rate limiter error: %v, falling back to allow traffic", err)
			c.Next()
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "Rate limiting error",
				"details": err.Error(),
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
