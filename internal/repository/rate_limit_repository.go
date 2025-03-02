package repository

import (
	"log"
	"rateLimiter/internal/database"
	"rateLimiter/internal/models"
	"time"

	"gorm.io/gorm"
)

type RateLimitRepository struct {
	db *gorm.DB
}

func NewRateLimitRepository() *RateLimitRepository {
	return &RateLimitRepository{db: database.GetDB()}
}

func(r *RateLimitRepository) GetRateLimit(identifier string) (*models.RateLimit, error) {
	var rateLimit models.RateLimit
	result := r.db.Raw("SELECT * FROM rate_limits WHERE identifier = ?", identifier).Scan(&rateLimit)
	log.Println(rateLimit)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &rateLimit, nil
}

func (r *RateLimitRepository) CreateRateLimit(identifier string, windowSize int) (*models.RateLimit, error) {
	now := time.Now()
	rateLimit := models.RateLimit{
		Identifier: identifier,
		Count: 1,
		LastRequest: now,
		ResetAt: now.Add(time.Duration(windowSize) * time.Second),
	}
	result := r.db.Create(&rateLimit)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rateLimit, nil
}

func (r *RateLimitRepository) IncrementRateLimit(rateLimit *models.RateLimit) error {
	rateLimit.Count++
	rateLimit.LastRequest = time.Now()
	return r.db.Save(rateLimit).Error
}

func (r *RateLimitRepository) ResetRateLimit(rateLimit *models.RateLimit, windowSize int) error {
	now := time.Now()
	rateLimit.Count = 1
	rateLimit.LastRequest = now
	rateLimit.ResetAt = now.Add(time.Duration(windowSize) * time.Second)
	return r.db.Save(rateLimit).Error
}