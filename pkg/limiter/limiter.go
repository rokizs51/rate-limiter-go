package limiter

import "time"

type RateLimiter interface {
	IsAllowed(identifier string) (bool, int, time.Time, error)
}
