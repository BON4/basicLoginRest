package redis

import (
	"basicLoginRest/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisCache(cfg *config.Config) *redis.Client {
	if cfg == nil {
		panic("redis options not specified")
	}

	return redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
		DB: cfg.Redis.Database,
		Password: cfg.Redis.Password,
		MaxRetries: cfg.Redis.MaxRetries,
		MaxRetryBackoff: cfg.Redis.MaxRetryBackoff,
	})
}
