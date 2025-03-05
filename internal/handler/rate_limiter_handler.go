package handlers

import (
	"net/http"
	"rateLimiter/internal/database"
	"rateLimiter/internal/repository"

	"github.com/gin-gonic/gin"
)

func GetRateLimiterData() gin.HandlerFunc {
	repo := repository.NewRedisTokenBucketRepository()

	return func(c *gin.Context) {
		// Get all keys matching the bucket pattern
		keys, err := database.RedisClient.Keys(c.Request.Context(), "bucket:*").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve rate limiter data",
			})
			return
		}

		data := make(map[string]interface{})
		for _, key := range keys {
			bucket, err := repo.GetBucket(c.Request.Context(), key[7:]) // Remove "bucket:" prefix
			if err != nil {
				continue
			}
			data[key[7:]] = bucket
		}

		c.JSON(http.StatusOK, gin.H{
			"rate_limits": data,
		})
	}
}
