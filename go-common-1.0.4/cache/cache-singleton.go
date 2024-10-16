package cache

import (
	"context"
	"time"
)

type ISingletonCache[T any] interface {
	Get(ctx context.Context) (*T, error)
	Set(ctx context.Context, t *T) error
	Delete(ctx context.Context) error
}

type cacheSingleton[T any] struct {
	base IKVCache[T, string]
	key  string
}

func (r *cacheSingleton[T]) Get(ctx context.Context) (*T, error) {
	return r.base.Get(ctx, r.key)
}

func (r *cacheSingleton[T]) Set(ctx context.Context, t *T) error {
	return r.base.Set(ctx, r.key, t)
}

func (r *cacheSingleton[T]) Delete(ctx context.Context) error {
	return r.base.Delete(ctx, r.key)
}
func CreateSingletonCache[T any](expiration time.Duration, key string) ISingletonCache[T] {
	return &cacheSingleton[T]{
		base: CreateKvCache[T, string](expiration, func(k string) string {
			return k
		}),
		key: key,
	}
}
