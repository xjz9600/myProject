package ratelimit

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"myProject/micro/grpc_demo/gen"
	"testing"
	"time"
)

func TestTokenBucketLimiter_BuildServerInterceptor(t *testing.T) {
	testCases := []struct {
		name     string
		ctx      context.Context
		b        func() *TokenBucketLimiter
		wantErr  error
		wantResp any
	}{
		{
			name: "closed",
			ctx:  context.Background(),
			b: func() *TokenBucketLimiter {
				closeChan := make(chan struct{})
				close(closeChan)
				return &TokenBucketLimiter{
					closeChan: closeChan,
				}
			},
			wantErr: errors.New("缺乏保护，拒绝请求"),
		},
		{
			name: "cancel",
			ctx: func() context.Context {
				ctx := context.Background()
				cancelCtx, cancel := context.WithCancel(ctx)
				cancel()
				return cancelCtx
			}(),
			b: func() *TokenBucketLimiter {
				return &TokenBucketLimiter{}
			},
			wantErr: context.Canceled,
		},
		{
			name: "get tokens",
			ctx:  context.Background(),
			b: func() *TokenBucketLimiter {
				ch := make(chan struct{}, 1)
				ch <- struct{}{}
				return &TokenBucketLimiter{
					tokens: ch,
				}
			},
			wantResp: "test",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := tc.b().BuildServerInterceptor()
			resp, err := interceptor(tc.ctx, &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, func(ctx context.Context, req any) (any, error) {
				return "test", nil
			})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantResp, resp)
		})
	}
}

func TestTokenBucketLimiter_Tokens(t *testing.T) {
	limiter := NewTokenBucketLimiter(2*time.Second, 10)
	defer limiter.Close()
	cnt := 0
	handler := func(ctx context.Context, req any) (any, error) {
		cnt++
		return "test", nil
	}
	interceptor := limiter.BuildServerInterceptor()
	_, err := interceptor(context.Background(), &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	assert.NoError(t, err)
	assert.Equal(t, cnt, 1)
	ctx := context.Background()
	cancelCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()
	_, err = interceptor(cancelCtx, &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	assert.Equal(t, err, context.DeadlineExceeded)
}
