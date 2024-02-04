package local_cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLocalCache_Get(t *testing.T) {
	testCases := []struct {
		name    string
		wantErr error
		wantVal any
		cache   func() *LocalCache
		key     string
	}{
		{
			name: "key not found",
			key:  "keyNotFound",
			cache: func() *LocalCache {
				return NewLocalCache(time.Second * 3)
			},
			wantErr: fmt.Errorf("%w key: %s", ErrKeyNotFound, "keyNotFound"),
		},
		{
			name: "expired",
			key:  "key1",
			cache: func() *LocalCache {
				res := NewLocalCache(time.Second * 5)
				err := res.Set(context.Background(), "key1", 123, time.Second)
				assert.NoError(t, err)
				time.Sleep(time.Second * 3)
				return res
			},
			wantErr: fmt.Errorf("%w key: %s", ErrKeyNotFound, "key1"),
		},
		{
			name: "get value",
			key:  "key1",
			cache: func() *LocalCache {
				res := NewLocalCache(time.Second * 3)
				err := res.Set(context.Background(), "key1", 123, time.Second)
				assert.NoError(t, err)
				return res
			},
			wantVal: 123,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.cache().Get(context.Background(), tc.key)
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, val, tc.wantVal)
		})
	}
}

func TestLocalCache_GetDelete(t *testing.T) {
	cnt := 0
	res := NewLocalCache(time.Second*100, WithOnEvictedLocalCache(func(key string, val any) {
		cnt++
	}))
	res.Set(context.Background(), "key1", 123, time.Second)
	time.Sleep(time.Second * 3)
	_, err := res.Get(context.Background(), "key1")
	assert.Equal(t, err, fmt.Errorf("%w key: %s", ErrKeyNotFound, "key1"))
	assert.Equal(t, 1, cnt)
}

func TestLocalCache_Loop(t *testing.T) {
	cnt := 0
	res := NewLocalCache(time.Second*2, WithOnEvictedLocalCache(func(key string, val any) {
		cnt++
	}))
	res.Set(context.Background(), "key1", 123, time.Second)
	time.Sleep(time.Second * 3)
	res.mutex.RLock()
	_, ok := res.items["key1"]
	res.mutex.RUnlock()
	assert.False(t, ok)
	assert.Equal(t, 1, cnt)
}
