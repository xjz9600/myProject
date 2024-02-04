package ratelimit

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"myProject/micro/grpc_demo/gen"
	"testing"
	"time"
)

func TestTokenLeakyBucketLimiter_Tokens(t *testing.T) {
	limiter := NewLeakyBucketLimiter(2 * time.Second)
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
