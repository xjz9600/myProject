package ratelimit

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"time"
)

//go:embed lua/slide_window.lua
var slideWindowScript string

type RedisSlideWindowLimiter struct {
	rate     int
	client   redis.Cmdable
	interval time.Duration
	service  string
}

func NewRedisSlideWindowLimiter(client redis.Cmdable, service string, interval time.Duration, rate int) *RedisSlideWindowLimiter {
	return &RedisSlideWindowLimiter{
		rate:     rate,
		client:   client,
		service:  service,
		interval: interval,
	}
}

func (f *RedisSlideWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		limit, err := f.limit(ctx)
		if err != nil {
			return
		}
		if limit {
			err = errors.New("触及了瓶颈")
			return
		}
		return handler(ctx, req)
	}
}

func (f *RedisSlideWindowLimiter) limit(ctx context.Context) (bool, error) {
	return f.client.Eval(ctx, slideWindowScript, []string{f.service}, time.Now().UnixMilli(), f.interval.Milliseconds(), f.rate).Bool()
}
