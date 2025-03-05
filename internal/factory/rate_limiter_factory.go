package factory

import (
	"context"
	"rateLimiter/internal/config"
	service "rateLimiter/internal/service/rate_limit"
	"time"
)

type RateLimiter interface {
	IsAllowed(ctx context.Context, identifier string) (bool, int, time.Time, error)
}

type Algorithm string

const (
	SlidingWindow Algorithm = "sliding_window"
	TokenBucket   Algorithm = "token_bucket"
)

func NewRateLimiter(algorithm Algorithm, cfg *config.Config) RateLimiter {
	switch algorithm {
	case TokenBucket:
		return service.NewTokenBucketService(&cfg.TokenBucketConfig)
	default:
		return service.NewSlidingWindowService(&cfg.SlidingWindowConfig)
	}
}
