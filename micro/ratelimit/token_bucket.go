package ratelimit

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"time"
)

type TokenBucketLimiter struct {
	tokens    chan struct{}
	closeChan chan struct{}
}

func NewTokenBucketLimiter(interval time.Duration, capacity int) *TokenBucketLimiter {
	tokenChan := make(chan struct{}, capacity)
	closeChan := make(chan struct{})
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case tokenChan <- struct{}{}:
				default:
					// 没有人使用
				}
			case <-closeChan:
				return
			}
		}
	}()
	return &TokenBucketLimiter{
		tokens: tokenChan,
	}
}

func (t *TokenBucketLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-t.tokens:
			return handler(ctx, req)
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-t.closeChan:
			err = errors.New("缺乏保护，拒绝请求")
			return
		}
	}
}

func (t *TokenBucketLimiter) Close() error {
	close(t.closeChan)
	return nil
}
