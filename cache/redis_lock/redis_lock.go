package redis_lock

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

//go:embed lua/redis-unlock.lua
var unlockScript string

//go:embed lua/redis-refresh.lua
var refreshScript string

//go:embed lua/redis-lock.lua
var lockScript string
var (
	ErrLockNotExist        = errors.New("redis-lock：解锁失败，锁不存在")
	ErrFailedToPreemptLock = errors.New("redis-lock：抢锁失败")
)

type RedisLock struct {
	client redis.Cmdable
}

func NewRedisLock(client redis.Cmdable) *RedisLock {
	return &RedisLock{
		client: client,
	}
}

func (r *RedisLock) Lock(ctx context.Context, key string, expiration, timeout time.Duration, retry Retry) (*lock, error) {
	val := uuid.New().String()
	var timer *time.Timer
	for {
		cx, cancel := context.WithTimeout(context.Background(), timeout)
		res, err := r.client.Eval(cx, lockScript, []string{key}, val, expiration.Seconds()).Result()
		cancel()
		if err != nil && err != context.DeadlineExceeded {
			return nil, err
		}
		if res == "OK" {
			return &lock{
				client:     r.client,
				key:        key,
				value:      val,
				expiration: expiration,
				stopChan:   make(chan struct{}),
			}, nil
		}
		duration, ok := retry.Next()
		if !ok {
			return nil, fmt.Errorf("redis-lock：超过重试限制，%w", ErrFailedToPreemptLock)
		}
		if timer == nil {
			timer = time.NewTimer(duration)
		} else {
			timer.Reset(duration)
		}
		select {
		case <-timer.C:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (r *RedisLock) TryLock(ctx context.Context, key string, expiration time.Duration) (*lock, error) {
	val := uuid.New().String()
	ok, err := r.client.SetNX(ctx, key, val, expiration).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrFailedToPreemptLock
	}
	return &lock{
		client:     r.client,
		key:        key,
		value:      val,
		expiration: expiration,
		stopChan:   make(chan struct{}),
	}, nil
}

type lock struct {
	client     redis.Cmdable
	key        string
	value      string
	expiration time.Duration
	stopChan   chan struct{}
	once       sync.Once
}

func (l *lock) AutoRefresh(interval time.Duration, timeout time.Duration) error {
	ticker := time.NewTicker(interval)
	timeoutChan := make(chan struct{}, 1)
	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			// 超时重试
			if err == context.DeadlineExceeded {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}
		case <-l.stopChan:
			return nil
		case <-timeoutChan:
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			err := l.Refresh(ctx)
			cancel()
			// 超时重试
			if err == context.DeadlineExceeded {
				timeoutChan <- struct{}{}
				continue
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *lock) Refresh(ctx context.Context) error {
	res, err := l.client.Eval(ctx, refreshScript, []string{l.key}, l.value, l.expiration.Seconds()).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotExist
	}
	return nil
}

func (l *lock) UnLock(ctx context.Context) error {
	res, err := l.client.Eval(ctx, unlockScript, []string{l.key}, l.value).Int64()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotExist
	}
	l.once.Do(func() {
		close(l.stopChan)
	})
	return nil
}
