package cache

import (
	"context"
	"github.com/eolinker/go-common/autowire"
	"time"
)

var (
	client ICommonCache
)

func init() {
	autowire.Autowired(&client)
}

type ICommonCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	GetInt(ctx context.Context, key string) (int64, error)
	Del(ctx context.Context, keys ...string) error
	Set(ctx context.Context, key string, val []byte, expiration time.Duration) error

	HMSet(ctx context.Context, key string, value map[string][]byte, expiration time.Duration) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) error

	Incr(ctx context.Context, key string, expiration time.Duration) error
	IncrBy(ctx context.Context, key string, val int64, expiration time.Duration) error

	SetNX(ctx context.Context, key string, val interface{}, expiration time.Duration) (bool, error)

	Clone() ICommonCache
}
