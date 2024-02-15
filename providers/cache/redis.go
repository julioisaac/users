package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/julioisaac/users/config"
	"time"
)

type redisCache struct {
	redisCli *redis.Client
}

func NewRedisCache() *redisCache {
	opts := &redis.Options{
		Addr:         config.GetString("REDIS_ADDR"),
		ReadTimeout:  config.GetDuration("REDIS_TIMEOUT_MS"),
		WriteTimeout: config.GetDuration("REDIS_WRITE_TIMEOUT_MS"),
		DialTimeout:  config.GetDuration("REDIS_TIMEOUT_MS"),
		MaxRetries:   config.GetInt("REDIS_MAX_RETRIES"),
		MinIdleConns: config.GetInt("REDIS_POOL_SIZE"),
		PoolTimeout:  config.GetDuration("REDIS_TIMEOUT_MS"),
	}
	r := redis.NewClient(opts)

	return &redisCache{redisCli: r}
}

func (r *redisCache) RemoveAll(ctx context.Context, keys ...string) (int64, error) {
	return r.redisCli.Del(ctx, keys...).Result()
}

func (r *redisCache) Add(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return r.redisCli.Set(ctx, key, value, expiration).Err()
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.redisCli.Get(ctx, key).Bytes()
}
