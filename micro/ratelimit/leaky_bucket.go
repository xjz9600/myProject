package ratelimit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

type LeakyBucketLimiter struct {
	tick      *time.Ticker
	closeChan chan struct{}
}

func NewLeakyBucketLimiter(interval time.Duration) *LeakyBucketLimiter {
	closeChan := make(chan struct{})
	return &LeakyBucketLimiter{
		tick:      time.NewTicker(interval),
		closeChan: closeChan,
	}
}

func (l *LeakyBucketLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-l.tick.C:
			return handler(ctx, req)
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-l.closeChan:
			err = errors.New("缺乏保护，拒绝请求")
			return
		}
	}
}

func (l *LeakyBucketLimiter) Close() error {
	l.tick.Stop()
	close(l.closeChan)
	return nil
}
