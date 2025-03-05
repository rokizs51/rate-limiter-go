package config

type Config struct {
	DatabaseDSN         string
	RedisConfig         RedisConfig
	SlidingWindowConfig SlidingWindowConfig
	TokenBucketConfig   TokenBucketConfig
}

type SlidingWindowConfig struct {
	Enabled      bool
	RequestLimit int
	WindowSize   int
}

type TokenBucketConfig struct {
	Enabled    bool
	Tokens     int
	RefillRate float64
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewConfig() *Config {
	return &Config{
		DatabaseDSN: "file:rateLimiter.db",
		RedisConfig: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		SlidingWindowConfig: SlidingWindowConfig{
			Enabled:      false,
			RequestLimit: 5,
			WindowSize:   60,
		},
		TokenBucketConfig: TokenBucketConfig{
			Enabled:    true,
			Tokens:     10,
			RefillRate: 0.5,
		},
	}
}
