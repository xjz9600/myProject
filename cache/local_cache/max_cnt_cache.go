package local_cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var (
	errOverCapacity = errors.New("超过容量限制")
)

type MaxCntCache struct {
	*LocalCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(l *LocalCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		LocalCache: l,
		maxCnt:     maxCnt,
	}
	origin := res.onEvicted
	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		origin(key, val)
	}
	return res
}

func (l *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.mutex.Lock()
	_, ok := l.items[key]
	defer l.mutex.Unlock()
	if !ok {
		if l.cnt+1 > l.maxCnt {
			return errOverCapacity
		}
		l.cnt++
	}
	return l.LocalCache.set(ctx, key, val, expiration)
}
