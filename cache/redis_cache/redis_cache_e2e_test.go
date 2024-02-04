package redis_cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisCacheE2E(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	c := NewRedisCache(rdb)
	err := c.Set(context.Background(), "key1", "value1", time.Second*3)
	require.NoError(t, err)
	val, err := c.Get(context.Background(), "key1")
	require.NoError(t, err)
	require.Equal(t, val, "value1")
}
