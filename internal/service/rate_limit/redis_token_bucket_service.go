package service

import (
	"context"
	"fmt"
	"rateLimiter/internal/config"
	"rateLimiter/internal/repository"
	"time"
)

type RedisTokenBucketService struct {
	repo   *repository.RedisTokenBucketRepository
	config *config.TokenBucketConfig
}

func NewRedisTokenBucketService(config *config.TokenBucketConfig) *RedisTokenBucketService {
	return &RedisTokenBucketService{
		repo:   repository.NewRedisTokenBucketRepository(),
		config: config,
	}
}

func (s *RedisTokenBucketService) IsAllowed(ctx context.Context, identifier string) (bool, int, time.Time, error) {
	bucket, err := s.repo.GetBucket(ctx, identifier)
	if err != nil {
		return false, 0, time.Time{}, err
	}
	fmt.Println(bucket)
	now := time.Now()
	if bucket == nil {
		bucket = &repository.BucketState{
			Tokens:     float64(s.config.Tokens),
			LastRefill: now,
		}
	}

	elapsed := now.Sub(bucket.LastRefill).Seconds()
	tokenToAdd := elapsed * s.config.RefillRate
	bucket.Tokens = min(float64(s.config.Tokens), bucket.Tokens+tokenToAdd)
	bucket.LastRefill = now

	if bucket.Tokens < 1 {
		return false, int(bucket.Tokens), bucket.LastRefill, nil
	}

	bucket.Tokens--
	err = s.repo.UpdateBucket(ctx, identifier, bucket, time.Hour*24)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	return true, int(bucket.Tokens), bucket.LastRefill, nil
}
