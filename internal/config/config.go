package config

type Config struct {
	DatabaseDSN         string
	SlidingWindowConfig SlidingWindowConfig
	TokenBucketConfig   TokenBucketConfig
}

type SlidingWindowConfig struct {
	Enabled      bool
	RequestLimit int
	WindowSize   int
}

type TokenBucketConfig struct {
	Tokens     int
	RefillRate float64
}

func NewConfig() *Config {
	return &Config{
		DatabaseDSN: "file:rateLimiter.db",
		SlidingWindowConfig: SlidingWindowConfig{
			Enabled:      true,
			RequestLimit: 5,
			WindowSize:   60,
		},
		TokenBucketConfig: TokenBucketConfig{
			Tokens:     10,
			RefillRate: 1.0,
		},
	}
}
