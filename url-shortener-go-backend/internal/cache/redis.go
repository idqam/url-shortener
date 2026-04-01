package cache

import (
	"context"
	"crypto/tls"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, password string, db int, defaultTTL time.Duration) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  password,
		DB:        0,
		TLSConfig: &tls.Config{},
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("redis ping failed", "error", err)
	} else {
		slog.Info("redis connected successfully")
	}
	return &RedisCache{client: rdb, ttl: defaultTTL}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		slog.Debug("redis cache miss", "key", key)
		return "", false, nil
	}
	if err != nil {
		slog.Error("redis get error", "error", err)
		return "", false, err
	}
	slog.Debug("redis cache hit", "key", key)
	return val, true, nil
}

func (c *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.ttl
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

func (c *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.client.Expire(ctx, key, ttl).Err()
}

func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
