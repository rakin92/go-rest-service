package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/logger"
)

// Cache is an exported cache object
type Cache struct {
	client redis.UniversalClient
	ttl    time.Duration
}

// Init creates a new redis cache
func Init(c *cfg.Cache) (*Cache, error) {
	logger.Info("[Cache.Init] Connecting to cache")
	client := redis.NewClient(&redis.Options{
		Addr:     c.Server,
		Password: c.Password,
	})
	timout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return nil, err
	}
	err = client.Ping().Err()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	logger.Info("[Cache.Init] Connected to cache")

	return &Cache{client: client, ttl: timout}, nil
}

// Add inserts items to cache
func (c *Cache) Add(ctx context.Context, key string, value string) (string, error) {
	s, err := c.client.Set(key, value, c.ttl).Result()
	if err == redis.Nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return s, nil
}

// Get returns items from cache given key
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	s, err := c.client.Get(key).Result()
	if err == redis.Nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return s, nil
}

// Del removes items from cache given key
func (c *Cache) Del(ctx context.Context, key string) (string, error) {
	_, err := c.client.Del(key).Result()
	if err != nil {
		return "", err
	}
	return "", nil
}
