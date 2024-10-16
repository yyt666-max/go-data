package cache

import (
	"context"
	"fmt"
	"time"
)

const (
	defaultExpiration = time.Minute
)

type IKVCache[T any, K comparable] interface {
	Get(ctx context.Context, k K) (*T, error)
	Set(ctx context.Context, k K, t *T) error
	Delete(ctx context.Context, keys ...K) error
}
type kvCache[T any, K comparable] struct {
	client        ICommonCache `autowired:""`
	formatHandler func(K) string
	expiration    time.Duration
}

func (r *kvCache[T, K]) Get(ctx context.Context, k K) (*T, error) {
	kv := r.formatHandler(k)

	bytes, err := r.client.Get(ctx, kv)
	if err != nil {
		return nil, err
	}

	return decode[T](bytes)

}
func (r *kvCache[T, K]) Set(ctx context.Context, k K, t *T) error {

	kv := r.formatHandler(k)

	bytes, err := encode(t)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, kv, bytes, r.expiration)
}

func (r *kvCache[T, K]) Delete(ctx context.Context, ks ...K) error {
	for _, k := range ks {
		key := r.formatHandler(k)
		if err := r.client.Del(ctx, key); err != nil {
			return err
		}

	}
	return nil
}
func CreateKvCache[T any, K comparable](expiration time.Duration, format ...func(k K) string) IKVCache[T, K] {

	if expiration == 0 {
		expiration = defaultExpiration
	}
	r := &kvCache[T, K]{
		expiration: expiration,
		client:     client,
	}

	if len(format) > 0 {
		r.formatHandler = format[0]
	} else {
		r.formatHandler = func(k K) string {
			return fmt.Sprint(k)
		}
	}

	return r
}
