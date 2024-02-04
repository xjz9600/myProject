package local_cache

import (
	"context"
	"errors"
	"fmt"
	"myProject/cache"
	"sync"
	"time"
)

var (
	ErrKeyNotFound = errors.New("cache：键值不存在")
)

var _ cache.Cache = &LocalCache{}

type LocalCache struct {
	mutex     sync.RWMutex
	items     map[string]*item
	onEvicted func(key string, val any)
	close     chan struct{}
	once      sync.Once
}

type item struct {
	data     any
	deadline time.Time
}

type LocalCacheOpt func(*LocalCache)

func WithOnEvictedLocalCache(onEvicted func(key string, val any)) LocalCacheOpt {
	return func(cache *LocalCache) {
		cache.onEvicted = onEvicted
	}
}

func NewLocalCache(interval time.Duration, opts ...LocalCacheOpt) *LocalCache {
	res := &LocalCache{
		items:     make(map[string]*item, 100),
		onEvicted: func(key string, val any) {},
		close:     make(chan struct{}),
	}
	for _, opt := range opts {
		opt(res)
	}
	tk := time.NewTicker(interval)
	go func() {
		for {
			select {
			case t := <-tk.C:
				res.mutex.Lock()
				i := 0
				for key, im := range res.items {
					if i > 1000 {
						break
					}
					if im.deadlineBefore(t) {
						res.delete(context.Background(), key)
					}
					i++
				}
				res.mutex.Unlock()
			case <-res.close:
				return
			}
		}
	}()
	return res
}

func (l *LocalCache) Close() {
	l.once.Do(func() {
		close(l.close)
	})
}

func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.set(ctx, key, val, expiration)
}

func (l *LocalCache) set(ctx context.Context, key string, val any, expiration time.Duration) error {
	var deadline time.Time
	if expiration > 0 {
		deadline = time.Now().Add(expiration)
	}
	l.items[key] = &item{
		data:     val,
		deadline: deadline,
	}
	return nil
}

func (i *item) deadlineBefore(t time.Time) bool {
	return !i.deadline.IsZero() && i.deadline.Before(t)
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {
	l.mutex.RLock()
	im, ok := l.items[key]
	l.mutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("%w key: %s", ErrKeyNotFound, key)
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	im, ok = l.items[key]
	if !ok {
		return nil, fmt.Errorf("%w key: %s", ErrKeyNotFound, key)
	}
	if im.deadlineBefore(time.Now()) {
		err := l.delete(ctx, key)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%w key: %s", ErrKeyNotFound, key)
	}
	return im.data, nil
}

func (l *LocalCache) Delete(ctx context.Context, key string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.delete(ctx, key)
	return nil
}

func (l *LocalCache) delete(ctx context.Context, key string) error {
	v, ok := l.items[key]
	if !ok {
		return fmt.Errorf("%w key: %s", ErrKeyNotFound, key)
	}
	delete(l.items, key)
	l.onEvicted(key, v)
	return nil
}

func (l *LocalCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	v, ok := l.items[key]
	if !ok {
		return nil, fmt.Errorf("%w key: %s", ErrKeyNotFound, key)
	}
	delete(l.items, key)
	l.onEvicted(key, v)
	return v, nil
}
