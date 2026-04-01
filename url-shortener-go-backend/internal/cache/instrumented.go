package cache

import (
	"context"
	"time"

	"url-shortener-go-backend/internal/metrics"
)

type InstrumentedCache struct {
	inner Cache
}

func NewInstrumentedCache(inner Cache) Cache {
	return &InstrumentedCache{inner: inner}
}

func (c *InstrumentedCache) Get(ctx context.Context, key string) (string, bool, error) {
	val, found, err := c.inner.Get(ctx, key)
	if err == nil {
		if found {
			metrics.CacheHitsTotal.WithLabelValues("get").Inc()
		} else {
			metrics.CacheMissesTotal.WithLabelValues("get").Inc()
		}
	}
	return val, found, err
}

func (c *InstrumentedCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.inner.Set(ctx, key, value, ttl)
}

func (c *InstrumentedCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.inner.Incr(ctx, key)
}

func (c *InstrumentedCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return c.inner.Expire(ctx, key, ttl)
}

func (c *InstrumentedCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.inner.TTL(ctx, key)
}

func (c *InstrumentedCache) Delete(ctx context.Context, key string) error {
	return c.inner.Delete(ctx, key)
}

func (c *InstrumentedCache) Ping(ctx context.Context) error {
	return c.inner.Ping(ctx)
}

func (c *InstrumentedCache) Close() error {
	return c.inner.Close()
}
