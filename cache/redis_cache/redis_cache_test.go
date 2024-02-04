package redis_cache

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name       string
		mock       func() redis.Cmdable
		wantErr    error
		key        string
		val        any
		expiration time.Duration
	}{
		{
			name: "set value",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewStatusCmd(context.Background())
				res.SetVal("OK")
				cmd.EXPECT().Set(context.Background(), "key1", "value1", time.Second).Return(res)
				return cmd
			},
			key:        "key1",
			val:        "value1",
			expiration: time.Second,
		},
		{
			name: "redis err",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewStatusCmd(context.Background())
				res.SetErr(errors.New("redis err"))
				cmd.EXPECT().Set(context.Background(), "key1", "value1", time.Second).Return(res)
				return cmd
			},
			key:        "key1",
			val:        "value1",
			expiration: time.Second,
			wantErr:    errors.New("redis err"),
		},
		{
			name: "redis cache err",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewStatusCmd(context.Background())
				res.SetVal("cache error")
				cmd.EXPECT().Set(context.Background(), "key1", "value1", time.Second).Return(res)
				return cmd
			},
			key:        "key1",
			val:        "value1",
			expiration: time.Second,
			wantErr:    errFailedToSetCache,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewRedisCache(tc.mock()).Set(context.Background(), tc.key, tc.val, tc.expiration)
			assert.Equal(t, err, tc.wantErr)
		})
	}
}
