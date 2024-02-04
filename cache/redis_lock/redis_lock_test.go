package redis_lock

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestRedisLock_TryLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name     string
		mock     func() redis.Cmdable
		key      string
		wantErr  error
		wantLock *lock
	}{
		{
			name: "redis err",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewBoolCmd(context.Background())
				res.SetErr(errors.New("redis err"))
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Second*3).Return(res)
				return cmd
			},
			key:     "key1",
			wantErr: errors.New("redis err"),
		},
		{
			name: "FailedToPreemptLock",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewBoolCmd(context.Background())
				res.SetVal(false)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Second*3).Return(res)
				return cmd
			},
			key:     "key1",
			wantErr: ErrFailedToPreemptLock,
		},
		{
			name: "redis err",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewBoolCmd(context.Background())
				res.SetVal(true)
				cmd.EXPECT().SetNX(context.Background(), "key1", gomock.Any(), time.Second*3).Return(res)
				return cmd
			},
			key: "key1",
			wantLock: &lock{
				key: "key1",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lk := NewRedisLock(tc.mock())
			locked, err := lk.TryLock(context.Background(), tc.key, time.Second*3)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, locked.key, tc.wantLock.key)
		})
	}
}

func TestRedisLock_UnLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name    string
		mock    func() redis.Cmdable
		key     string
		wantErr error
		value   string
	}{
		{
			name: "redis err",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("redis err"))
				cmd.EXPECT().Eval(context.Background(), unlockScript, []string{"key1"}, "value1").Return(res)
				return cmd
			},
			key:     "key1",
			value:   "value1",
			wantErr: errors.New("redis err"),
		},
		{
			name: "LockNotExist",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(context.Background(), unlockScript, []string{"key1"}, "value1").Return(res)
				return cmd
			},
			key:     "key1",
			value:   "value1",
			wantErr: ErrLockNotExist,
		},
		{
			name: "unlock",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))
				cmd.EXPECT().Eval(context.Background(), unlockScript, []string{"key1"}, "value1").Return(res)
				return cmd
			},
			key:   "key1",
			value: "value1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			locked := &lock{
				client:   tc.mock(),
				key:      tc.key,
				value:    tc.value,
				stopChan: make(chan struct{}),
			}
			err := locked.UnLock(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestRedisLock_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name       string
		mock       func() redis.Cmdable
		key        string
		wantErr    error
		value      string
		expiration time.Duration
	}{
		{
			name: "redis err",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("redis err"))
				cmd.EXPECT().Eval(context.Background(), refreshScript, []string{"key1"}, "value1", time.Second.Seconds()).Return(res)
				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
			wantErr:    errors.New("redis err"),
		},
		{
			name: "LockNotExist",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(context.Background(), refreshScript, []string{"key1"}, "value1", time.Second.Seconds()).Return(res)
				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
			wantErr:    ErrLockNotExist,
		},
		{
			name: "unlock",
			mock: func() redis.Cmdable {
				cmd := NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(1))
				cmd.EXPECT().Eval(context.Background(), refreshScript, []string{"key1"}, "value1", time.Second.Seconds()).Return(res)
				return cmd
			},
			key:        "key1",
			value:      "value1",
			expiration: time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			locked := &lock{
				client:     tc.mock(),
				key:        tc.key,
				value:      tc.value,
				expiration: tc.expiration,
			}
			err := locked.Refresh(context.Background())
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func ExampleLock_Refresh() {
	var lk *lock
	// 关闭续约信号
	stopChan := make(chan struct{})
	errChan := make(chan error)
	timeoutChan := make(chan struct{}, 1)
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ticker.C:
				err := lk.Refresh(context.Background())
				// 超时重试
				if err == context.DeadlineExceeded {
					timeoutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
					return
				}
			case <-stopChan:
				return
			case <-timeoutChan:
				err := lk.Refresh(context.Background())
				// 超时重试
				if err == context.DeadlineExceeded {
					timeoutChan <- struct{}{}
					continue
				}
				if err != nil {
					errChan <- err
					return
				}
			}
		}
	}()
	// 业务自己打点判断是否有错误，继续下一步
	select {
	case er := <-errChan:
		log.Fatal(er)
	default:
		// 业务处理
	}
	close(stopChan)
}

func ExampleLock_AutoRefresh() {
	var lk *lock
	go lk.AutoRefresh(time.Second*3, time.Second)
}
