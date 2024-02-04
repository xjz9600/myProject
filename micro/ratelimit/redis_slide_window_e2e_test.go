package ratelimit

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"myProject/micro/grpc_demo/gen"
	"testing"
	"time"
)

func TestRedisSLideWindowLimiter_BuildServerInterceptor(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cnt := 0
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		cnt++
		return &gen.GetByIdResp{}, nil
	}
	interceptor := NewRedisSlideWindowLimiter(rdb, "user-service", time.Second*3, 1).BuildServerInterceptor()
	resp, err := interceptor(context.Background(), &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	require.NoError(t, err)
	assert.Equal(t, resp, &gen.GetByIdResp{})

	resp, err = interceptor(context.Background(), &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	assert.Equal(t, err, errors.New("触及了瓶颈"))
	require.Nil(t, resp)

	time.Sleep(3 * time.Second)
	resp, err = interceptor(context.Background(), &gen.GetByIdReq{}, &grpc.UnaryServerInfo{}, handler)
	require.NoError(t, err)
	assert.Equal(t, resp, &gen.GetByIdResp{})
}
