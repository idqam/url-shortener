package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
}
