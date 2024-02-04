package cache

import (
	"context"
	"errors"
	"golang.org/x/sync/singleflight"
	"myProject/cache/local_cache"
	"time"
)

type ReadThroughCache struct {
	Cache
	storeFunc  func(ctx context.Context, key string, expiration time.Duration) (any, error)
	expiration time.Duration
	g          singleflight.Group
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if !errors.Is(err, local_cache.ErrKeyNotFound) {
		return nil, err
	}
	val, err = r.storeFunc(ctx, key, r.expiration)
	if err != nil {
		return nil, err
	}
	err = r.Cache.Set(ctx, key, val, r.expiration)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *ReadThroughCache) GetSingleFlight(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if !errors.Is(err, local_cache.ErrKeyNotFound) {
		return nil, err
	}
	val, err = r.storeFunc(ctx, key, r.expiration)
	if err != nil {
		return nil, err
	}
	val, err, _ = r.g.Do(key, func() (interface{}, error) {
		err = r.Cache.Set(ctx, key, val, r.expiration)
		if err != nil {
			return nil, err
		}
		return val, nil
	})

	return val, err
}
