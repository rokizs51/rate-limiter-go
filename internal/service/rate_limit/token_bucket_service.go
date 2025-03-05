package service

import (
	"context"
	"rateLimiter/internal/config"
	"rateLimiter/internal/repository"
	"time"
)

type TokenBucketService struct {
	repo   *repository.TokenBucketRepository
	config *config.TokenBucketConfig
}

func NewTokenBucketService(config *config.TokenBucketConfig) *TokenBucketService {
	return &TokenBucketService{
		repo:   repository.NewTokenBucketRepository(),
		config: config,
	}
}

func (s *TokenBucketService) IsAllowed(ctx context.Context, identifier string) (bool, int, time.Time, error) {
	bucket, err := s.repo.GetBucket(ctx, identifier)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	now := time.Now()
	if bucket == nil {
		bucket, err = s.repo.CreateBucket(ctx, identifier, float64(s.config.Tokens))
		if err != nil {
			return false, 0, time.Time{}, err
		}
	}

	elapsed := now.Sub(bucket.LastRefill).Seconds()
	tokenToAdd := elapsed * s.config.RefillRate
	bucket.Tokens = min(float64(s.config.Tokens), bucket.Tokens+tokenToAdd)
	bucket.LastRefill = now

	if bucket.Tokens < 1 {
		nextRefillTime := now.Add(time.Duration((1-bucket.Tokens)/s.config.RefillRate) * time.Second)
		err := s.repo.UpdateBucket(ctx, bucket)
		if err != nil {
			return false, 0, time.Time{}, err
		}
		return false, int(bucket.Tokens), nextRefillTime, nil
	}

	bucket.Tokens--
	err = s.repo.UpdateBucket(ctx, bucket)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	nextRefillTime := now.Add(time.Duration((1-bucket.Tokens)/s.config.RefillRate) * time.Second)
	return true, int(bucket.Tokens), nextRefillTime, nil
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
