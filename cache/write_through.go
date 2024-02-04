package cache

import (
	"context"
	"golang.org/x/sync/singleflight"
	"log"
	"time"
)

type WriteThroughCache struct {
	Cache
	storeFunc func(ctx context.Context, key string, val any, expiration time.Duration) error
	g         singleflight.Group
}

func (w *WriteThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.storeFunc(ctx, key, val, expiration)
	if err != nil {
		return err
	}
	return w.Cache.Set(ctx, key, val, expiration)
}

func (w *WriteThroughCache) SetSingleFlight(ctx context.Context, key string, val any, expiration time.Duration) error {
	go func() {
		w.g.Do(key, func() (interface{}, error) {
			err := w.storeFunc(ctx, key, val, expiration)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			err = w.Cache.Set(ctx, key, val, expiration)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}
			return nil, nil
		})
	}()
	return nil
}
