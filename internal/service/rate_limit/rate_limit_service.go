package service

import (
	"rateLimiter/internal/config"
	"rateLimiter/internal/repository"
	"time"
)

type RateLimitService struct {
	repo *repository.RateLimitRepository
	config *config.RateLimitConfig
}

func NewRateLimitService(config *config.RateLimitConfig) *RateLimitService {
	return &RateLimitService{repo: repository.NewRateLimitRepository(), config: config}
}

func (s *RateLimitService) IsAllowed(identifier string) (bool, int, time.Time, error ) {
	if !s.config.Enabled {
		return true, 0, time.Time{}, nil
	}
	
	rateLimit, err := s.repo.GetRateLimit(identifier)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	now := time.Now()

	//iff no rate limit exist create one
	if rateLimit == nil {
		rateLimit, err = s.repo.CreateRateLimit(identifier, s.config.WindowSize)
		if err != nil {
			return false, 0, time.Time{}, err
		}
		return true, 1, rateLimit.ResetAt, nil
	}

	//if rate limit has expired, reset it
	if now.After(rateLimit.ResetAt) {
		err := s.repo.ResetRateLimit(rateLimit, s.config.WindowSize)
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
	err = s.repo.IncrementRateLimit(rateLimit)
	if err != nil {
		return false, 0, time.Time{}, err
	}

	return true, rateLimit.Count, rateLimit.ResetAt, nil
}