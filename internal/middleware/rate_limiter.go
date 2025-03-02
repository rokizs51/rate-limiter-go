package middleware

import (
	"net/http"
	"rateLimiter/internal/config"
	service "rateLimiter/internal/service/rate_limit"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(cfg *config.RateLimitConfig) gin.HandlerFunc {
	rateLimitService := service.NewRateLimitService(cfg)

	return func(c *gin.Context) {
		identifier := c.ClientIP()
		allowed, count, resetAt, err := rateLimitService.IsAllowed(identifier)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error" : "Rate limiting error",
			})
			return
		}
		//set rate limit headers
		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.RequestLimit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(cfg.RequestLimit - count))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

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