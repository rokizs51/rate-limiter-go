package service

import (
	"context"
	"rateLimiter/internal/config"
	"rateLimiter/internal/repository"
	"time"
)

type SlidingWindowService struct {
	repo   *repository.SlidingWindowRepository
	config *config.SlidingWindowConfig
}

func NewSlidingWindowService(config *config.SlidingWindowConfig) *SlidingWindowService {
	return &SlidingWindowService{repo: repository.NewRateLimitRepository(), config: config}
}

func (s *SlidingWindowService) IsAllowed(ctx context.Context, identifier string) (bool, int, time.Time, error) {
	if !s.config.Enabled {
		return true, 0, time.Time{}, nil
	}

	rateLimit, err := s.repo.GetRateLimit(ctx, identifier)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	now := time.Now()

	//iff no rate limit exist create one
	if rateLimit == nil {
		rateLimit, err = s.repo.CreateRateLimit(ctx, identifier, s.config.WindowSize)
		if err != nil {
			return false, 0, time.Time{}, err
		}
		return true, 1, rateLimit.ResetAt, nil
	}

	//if rate limit has expired, reset it
	if now.After(rateLimit.ResetAt) {
		err := s.repo.ResetRateLimit(ctx, rateLimit, s.config.WindowSize)
		if err != nil {
			return false, 0, time.Time{}, err
		}
		return true, 1, rateLimit.ResetAt, nil
	}

	//check if rate limit exceeded
	if rateLimit.Count >= s.config.RequestLimit {
		return false, rateLimit.Count, rateLimit.ResetAt, nil
	}

	// increment rate limit
	err = s.repo.IncrementRateLimit(ctx, rateLimit)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	return true, rateLimit.Count, rateLimit.ResetAt, nil
}
