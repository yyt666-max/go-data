package cache

import (
	"context"
	"time"
)

type IListCache[T any] interface {
	SetAll(ctx context.Context, t []T) error
	GetAll(ctx context.Context) ([]T, error)
	Delete(ctx context.Context) error
}
type listCache[T any] struct {
	client     ICommonCache
	key        string
	expiration time.Duration
}

func CreateListCache[T any](expiration time.Duration, key string) IListCache[T] {

	r := &listCache[T]{
		key:        key,
		expiration: expiration,
		client:     client,
	}

	return r
}
func (r *listCache[T]) Delete(ctx context.Context) error {
	return r.client.Del(ctx, r.key)
}

func (r *listCache[T]) GetAll(ctx context.Context) ([]T, error) {

	bytes, err := r.client.Get(ctx, r.key)
	if err != nil {
		return nil, err
	}

	return decodeList[T](bytes)

}

func (r *listCache[T]) SetAll(ctx context.Context, t []T) error {

	bytes, err := encode(t)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.key, bytes, r.expiration)
}
