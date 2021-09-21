package credis

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"sync"
	"time"
	"url-shortner/domain"
)

type RedisCache struct {
	m     sync.Mutex
	Cache *cache.Cache
}

func NewRedisCache(client *redis.Client) *RedisCache {
	cache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Hour),
	})
	return &RedisCache{
		Cache: cache,
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, value []domain.Url, expiration time.Duration) error {
	c.m.Lock()
	defer c.m.Unlock()
	err := c.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   expiration,
	})
	if err != nil {
		return errors.Wrapf(err, "could not cache the %v", value)
	}
	return nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]domain.Url, error) {
	c.m.Lock()
	defer c.m.Unlock()
	cachedUrls := make([]domain.Url, 0)
	err := c.Cache.Get(ctx, key, &cachedUrls)
	if err != nil {
		return nil, errors.Wrap(err, "could not get the cached value")
	}
	return cachedUrls, nil
}
