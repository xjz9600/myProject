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

func TestTokenSlideWindowLimiter_Tokens(t *testing.T) {
	limiter := NewSlideWindowLimiter(2*time.Second.Nanoseconds(), 1)
	cnt := 0
	handler := func(ctx context.Context, req any) (any, error) {
		cnt++
		return "test", nil
	}
	interceptor := limiter.BuildServerInterceptor()
	_, err := interceptor(context.Background(), &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	assert.NoError(t, err)
	assert.Equal(t, cnt, 1)
	_, err = interceptor(context.Background(), &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	assert.Equal(t, err, errors.New("触发瓶颈了"))
	time.Sleep(2 * time.Second)
	_, err = interceptor(context.Background(), &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	assert.NoError(t, err)
	assert.Equal(t, cnt, 2)
}
