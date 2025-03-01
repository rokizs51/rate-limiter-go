package config

type Config struct {
	DatabaseDSN string
	RateLimit RateLimitConfig
}

type RateLimitConfig struct {
	Enabled bool
	RequestLimit int
	WindowSize int
}


func NewConfig() *Config {
	return &Config{
		DatabaseDSN: "file:rateLimiter.db",
		RateLimit: RateLimitConfig{
			Enabled: true,
			RequestLimit: 5,
			WindowSize: 60,
		},
	}
}