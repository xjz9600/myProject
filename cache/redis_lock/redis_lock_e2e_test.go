//go:build e2e

package redis_lock

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisLock_e2e_TryLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name     string
		key      string
		wantErr  error
		wantLock *lock
		before   func()
		after    func()
	}{
		{
			name:    "FailedToPreemptLock",
			key:     "key1",
			wantErr: ErrFailedToPreemptLock,
			before: func() {
				res, err := rdb.Set(context.Background(), "key1", "value1", time.Second*3).Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "OK")
			},
			after: func() {
				res, err := rdb.Get(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "value1")
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
		},
		{
			name:   "locked",
			key:    "key1",
			before: func() {},
			after: func() {
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
			wantLock: &lock{
				key: "key1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			locked := NewRedisLock(rdb)
			lk, err := locked.TryLock(context.Background(), tc.key, time.Second*3)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, lk.key, tc.wantLock.key)
			assert.NotEmpty(t, lk.value)
		})
	}
}

func TestRedisLock_e2e_UnLock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name    string
		key     string
		wantErr error
		value   string
		before  func()
		after   func()
	}{
		{
			name:    "ErrLockNotExist",
			before:  func() {},
			after:   func() {},
			key:     "key1",
			value:   "value1",
			wantErr: ErrLockNotExist,
		},
		{
			name: "Lock_hold_by_other",
			before: func() {
				re, err := rdb.Set(context.Background(), "key1", "otherValue", time.Second*3).Result()
				assert.NoError(t, err)
				assert.Equal(t, re, "OK")

			},
			after: func() {
				res, err := rdb.Get(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "otherValue")
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
			key:     "key1",
			value:   "value1",
			wantErr: ErrLockNotExist,
		},
		{
			name: "Lock",
			before: func() {
				re, err := rdb.Set(context.Background(), "key1", "value1", time.Second*3).Result()
				assert.NoError(t, err)
				assert.Equal(t, re, "OK")

			},
			after: func() {
				res, err := rdb.Exists(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, res, int64(0))

			},
			key:   "key1",
			value: "value1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			locked := &lock{
				key:      tc.key,
				value:    tc.value,
				client:   rdb,
				stopChan: make(chan struct{}),
			}
			err := locked.UnLock(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRedisLock_e2e_Refresh(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name       string
		key        string
		wantErr    error
		value      string
		before     func()
		after      func()
		expiration time.Duration
	}{
		{
			name:       "ErrLockNotExist",
			before:     func() {},
			after:      func() {},
			key:        "key1",
			value:      "value1",
			wantErr:    ErrLockNotExist,
			expiration: time.Minute,
		},
		{
			name: "Lock_hold_by_other",
			before: func() {
				re, err := rdb.Set(context.Background(), "key1", "otherValue", time.Second*3).Result()
				assert.NoError(t, err)
				assert.Equal(t, re, "OK")

			},
			after: func() {
				res, err := rdb.Get(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "otherValue")
				duration, err := rdb.TTL(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.True(t, duration.Seconds() <= time.Second.Seconds()*3)
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
			key:        "key1",
			value:      "value1",
			wantErr:    ErrLockNotExist,
			expiration: time.Minute,
		},
		{
			name: "Refresh",
			before: func() {
				re, err := rdb.Set(context.Background(), "key1", "value1", time.Second*3).Result()
				assert.NoError(t, err)
				assert.Equal(t, re, "OK")

			},
			after: func() {
				duration, err := rdb.TTL(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.True(t, duration.Seconds() > time.Second.Seconds()*3)
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))

			},
			key:        "key1",
			value:      "value1",
			expiration: time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			locked := &lock{
				key:        tc.key,
				value:      tc.value,
				client:     rdb,
				expiration: tc.expiration,
			}
			err := locked.Refresh(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRedisLock_e2e_Lock(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	testCases := []struct {
		name       string
		key        string
		wantErr    error
		wantLock   *lock
		retry      Retry
		before     func()
		after      func()
		expiration time.Duration
		timeout    time.Duration
		context    context.Context
	}{
		{
			name:   "locked",
			key:    "key1",
			before: func() {},
			after: func() {
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
			timeout:    time.Second,
			expiration: time.Second * 3,
			wantLock: &lock{
				key:        "key1",
				expiration: time.Second * 3,
			},
			context: context.Background(),
		},
		{
			name:    "timeout",
			key:     "key1",
			wantErr: context.DeadlineExceeded,
			before: func() {
				res, err := rdb.Set(context.Background(), "key1", "value1", time.Second*3).Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "OK")
			},
			after: func() {
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
			timeout:    time.Second,
			expiration: time.Second * 3,
			context: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return ctx
			}(),
			retry: &TestRetry{
				maxCnt:   10,
				duration: time.Second,
			},
		},
		{
			name:    "Lock_hold_by_other",
			key:     "key1",
			wantErr: fmt.Errorf("redis-lock：超过重试限制，%w", ErrFailedToPreemptLock),
			before: func() {
				res, err := rdb.Set(context.Background(), "key1", "value1", time.Minute).Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "OK")
			},
			after: func() {
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
			timeout:    time.Second * 2,
			expiration: time.Second * 3,
			context:    context.Background(),
			retry: &TestRetry{
				maxCnt:   3,
				duration: time.Second,
			},
		},
		{
			name: "retry and locked",
			key:  "key1",
			before: func() {
				res, err := rdb.Set(context.Background(), "key1", "value1", time.Second*2).Result()
				assert.NoError(t, err)
				assert.Equal(t, res, "OK")
			},
			after: func() {
				duration, err := rdb.TTL(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.True(t, time.Second.Seconds()*30 < duration.Seconds())
				re, err := rdb.Del(context.Background(), "key1").Result()
				assert.NoError(t, err)
				assert.Equal(t, re, int64(1))
			},
			timeout:    time.Second * 2,
			expiration: time.Minute,
			context:    context.Background(),
			retry: &TestRetry{
				maxCnt:   3,
				duration: time.Second,
			},
			wantLock: &lock{
				key:        "key1",
				expiration: time.Second * 3,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			locked := NewRedisLock(rdb)
			lk, err := locked.Lock(tc.context, tc.key, tc.expiration, tc.timeout, tc.retry)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, lk.key, tc.wantLock.key)
			assert.NotEmpty(t, lk.value)
		})
	}
}
