package repository

import (
	"context"
	"encoding/json"
	"rateLimiter/internal/database"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenBucketRepository struct {
	redis *redis.Client
}

type BucketState struct {
	Tokens     float64   `json:"tokens"`
	LastRefill time.Time `json:"last_refill"`
}

func NewRedisTokenBucketRepository() *RedisTokenBucketRepository {
	return &RedisTokenBucketRepository{redis: database.GetRedis()}
}

func (r *RedisTokenBucketRepository) GetBucket(ctx context.Context, identifier string) (*BucketState, error) {
	data, err := r.redis.Get(ctx, "bucket:"+identifier).Result()
	// This condition is wrong - it should check if err == redis.Nil
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	var bucket BucketState
	err = json.Unmarshal([]byte(data), &bucket)
	if err != nil {
		return nil, err
	}
	return &bucket, nil
}

func (r *RedisTokenBucketRepository) UpdateBucket(ctx context.Context, identifier string, bucket *BucketState, expiry time.Duration) error {
	data, err := json.Marshal(bucket)
	if err != nil {
		return err
	}

	return r.redis.Set(ctx, "bucket:"+identifier, data, expiry).Err()
}
