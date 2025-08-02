package cache

import (
	"context"
	"crypto/tls"
	"log"
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
		log.Printf("[ERROR] Redis ping failed: %v", err)
	} else {
		log.Println("[INFO] Redis connected successfully")
	}
	return &RedisCache{client: rdb, ttl: defaultTTL}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[Redis] MISS for key: %s", key)
		return "", false, nil
	}
	if err != nil {
		log.Printf("[Redis] ERROR: %v", err)
		return "", false, err
	}
	log.Printf("[Redis] HIT for key: %s", key)
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
