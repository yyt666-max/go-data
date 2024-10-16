package cache_redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	redis "github.com/redis/go-redis/v9"

	"github.com/eolinker/go-common/cache"
)

type commonCache struct {
	client    redis.UniversalClient
	keyPrefix string
}

func (c *commonCache) Clone() cache.ICommonCache {

	return &commonCache{
		client:    c.client,
		keyPrefix: c.keyPrefix,
	}
}

func newCommonCache(client redis.UniversalClient, namespace string) cache.ICommonCache {

	if namespace == "" {
		namespace = "apinto"
	}
	namespace = fmt.Sprint(strings.Trim(namespace, ":"), ":")
	return &commonCache{client: client, keyPrefix: namespace}
}

func (c *commonCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, c.key(key)).Bytes()
}

func (c *commonCache) Set(ctx context.Context, key string, val []byte, expiration time.Duration) error {
	return c.client.Set(ctx, c.key(key), val, expiration).Err()
}

func (c *commonCache) Incr(ctx context.Context, key string, expiration time.Duration) error {
	redisKey := c.key(key)
	err := c.client.Incr(ctx, redisKey).Err()
	if err != nil {
		return err
	}
	return c.client.Expire(ctx, redisKey, expiration).Err()
}
func (c *commonCache) key(v string) string {
	return fmt.Sprint(c.keyPrefix, v)
}
func (c *commonCache) IncrBy(ctx context.Context, key string, val int64, expiration time.Duration) error {
	redisKey := c.key(key)
	err := c.client.IncrBy(ctx, redisKey, val).Err()
	if err != nil {
		return err
	}
	return c.client.Expire(ctx, redisKey, expiration).Err()
}

func (c *commonCache) GetInt(ctx context.Context, key string) (int64, error) {
	redisKey := c.key(key)
	return c.client.Get(ctx, redisKey).Int64()
}

func (c *commonCache) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		if err := c.client.Del(ctx, c.key(key)).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (c *commonCache) HMSet(ctx context.Context, key string, value map[string][]byte, expiration time.Duration) error {
	values := make([]interface{}, 0)
	for k, val := range value {
		values = append(values, k, val)
	}
	if err := c.client.HMSet(ctx, c.key(key), values...).Err(); err != nil {
		return err
	}
	c.client.Expire(ctx, c.key(key), expiration)
	return nil
}

func (c *commonCache) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, c.key(key), fields...).Err()
}

func (c *commonCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, c.key(key)).Result()
}

func (c *commonCache) SetNX(ctx context.Context, key string, val interface{}, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, c.key(key), val, expiration).Result()
}
