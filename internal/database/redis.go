package database

import (
	"context"
	"fmt"
	"log"
	"rateLimiter/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedis(cfg *config.Config) error {
	var redisErr error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.RedisConfig.Host, cfg.RedisConfig.Port),
			Password: cfg.RedisConfig.Password,
			DB:       cfg.RedisConfig.DB,
		})

		// Test connection
		if err := RedisClient.Ping(context.Background()).Err(); err != nil {
			redisErr = err
			log.Printf("Failed to connect to Redis (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(time.Second * 2) // Wait before retrying
			continue
		}

		log.Println("Redis connection established successfully")
		return nil
	}

	return fmt.Errorf("failed to connect to Redis after %d attempts: %v", maxRetries, redisErr)
}

func GetRedis() *redis.Client {
	return RedisClient
}
