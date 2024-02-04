package ratelimit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type FixWindowLimiter struct {
	timestamp int64
	interval  int64
	rate      int64
	cnt       int64
}

func NewFixWindowLimiter(interval time.Duration, rate int64) *FixWindowLimiter {
	return &FixWindowLimiter{
		timestamp: time.Now().UnixNano(),
		interval:  interval.Nanoseconds(),
		rate:      rate,
	}
}

func (f *FixWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		now := time.Now().UnixNano()
		timestamp := atomic.LoadInt64(&f.timestamp)
		cnt := atomic.LoadInt64(&f.cnt)
		if now > timestamp+f.interval {
			if atomic.CompareAndSwapInt64(&f.timestamp, timestamp, now) {
				atomic.CompareAndSwapInt64(&f.cnt, cnt, 0)
			}
		}
		newCnt := atomic.AddInt64(&f.cnt, 1)
		if newCnt > f.rate {
			err = errors.New("触发瓶颈了")
			return
		}
		return handler(ctx, req)
	}
}
