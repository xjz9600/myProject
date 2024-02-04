package ratelimit

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"time"
)

//go:embed lua/fix_window.lua
var fixWindowScript string

type RedisFixWindowLimit struct {
	client   redis.Cmdable
	service  string
	rate     int
	interval time.Duration
}

func NewRedisFixWindowLimit(client redis.Cmdable, service string, rate int, interval time.Duration) *RedisFixWindowLimit {
	return &RedisFixWindowLimit{
		client:   client,
		service:  service,
		rate:     rate,
		interval: interval,
	}
}

func (f *RedisFixWindowLimit) BuildServerInterceptor() grpc.UnaryServerInterceptor {
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

func (f *RedisFixWindowLimit) limit(ctx context.Context) (bool, error) {
	return f.client.Eval(ctx, fixWindowScript, []string{f.service}, f.rate, f.interval.Milliseconds()).Bool()
}
