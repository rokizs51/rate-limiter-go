package repository

import (
	"rateLimiter/internal/database"
	"rateLimiter/internal/models"
	"time"

	"gorm.io/gorm"
)

type TokenBucketRepository struct {
	db *gorm.DB
}

func NewTokenBucketRepository() *TokenBucketRepository {
	return &TokenBucketRepository{db: database.GetDB()}
}

func (r *TokenBucketRepository) GetBucket(identifier string) (*models.TokenBucket, error) {
	var bucket models.TokenBucket
	result := r.db.Where("identifier = ?", identifier).First(&bucket)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &bucket, nil
}

func (r *TokenBucketRepository) CreateBucket(identifier string, token float64) (*models.TokenBucket, error) {
	bucket := models.TokenBucket{
		Identifier: identifier,
		Tokens:     token,
		LastRefill: time.Now(),
	}

	err := r.db.Create(&bucket).Error
	if err != nil {
		return nil, err
	}

	return &bucket, nil
}

func (r *TokenBucketRepository) UpdateBucket(bucket *models.TokenBucket) error {
	return r.db.Save(bucket).Error
}
