package leetcode

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrentQueue_InsertLast(t *testing.T) {
	testCases := []struct {
		name      string
		q         func() *concurrentDeque[int]
		val       int
		timeout   time.Duration
		wantErr   error
		wantData  []int
		wantSlice []int
		wantLen   int
		wantHead  int
		wantTail  int
	}{
		{
			name: "empty and enqueued",
			q: func() *concurrentDeque[int] {
				return ConcurrentConstructor[int](3)
			},
			val:       123,
			timeout:   time.Second,
			wantData:  []int{123, 0, 0},
			wantSlice: []int{123},
			wantLen:   1,
			wantTail:  1,
			wantHead:  0,
		},
		{
			name: "invalid context",
			q: func() *concurrentDeque[int] {
				return ConcurrentConstructor[int](3)
			},
			val:       123,
			timeout:   -time.Second,
			wantData:  []int{0, 0, 0},
			wantSlice: []int{},
			wantErr:   context.DeadlineExceeded,
		},
		{
			// 入队之后就满了，恰好放在切片的最后一个位置
			name: "enqueued full last index",
			q: func() *concurrentDeque[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := ConcurrentConstructor[int](3)
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 234)
				require.NoError(t, err)
				return q
			},
			val:       345,
			timeout:   time.Second,
			wantData:  []int{123, 234, 345},
			wantSlice: []int{123, 234, 345},
			wantLen:   3,
			wantTail:  0,
			wantHead:  0,
		},
		{
			// 入队之后就满了，恰好放在切片的第一个
			name: "enqueued full middle index",
			q: func() *concurrentDeque[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := ConcurrentConstructor[int](3)
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 234)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 345)
				require.NoError(t, err)
				val, err := q.GetFront(ctx)
				require.NoError(t, err)
				require.Equal(t, 123, val)
				return q
			},
			val:       456,
			timeout:   time.Second,
			wantData:  []int{456, 234, 345},
			wantSlice: []int{234, 345, 456},
			wantLen:   3,
			wantTail:  1,
			wantHead:  1,
		},
		{
			// 入队之后就满了，恰好放在中间
			name: "enqueued full first index",
			q: func() *concurrentDeque[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := ConcurrentConstructor[int](3)
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 234)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 345)
				require.NoError(t, err)
				val, err := q.GetFront(ctx)
				require.NoError(t, err)
				require.Equal(t, 123, val)
				val, err = q.GetFront(ctx)
				require.NoError(t, err)
				require.Equal(t, 234, val)
				err = q.InsertLast(ctx, 456)
				require.NoError(t, err)
				return q
			},
			val:       567,
			timeout:   time.Second,
			wantData:  []int{456, 567, 345},
			wantSlice: []int{345, 456, 567},
			wantLen:   3,
			wantTail:  2,
			wantHead:  2,
		},
		{
			// 元素本身就是零值
			name: "all zero value ",
			q: func() *concurrentDeque[int] {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				q := ConcurrentConstructor[int](3)
				err := q.InsertLast(ctx, 0)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 0)
				require.NoError(t, err)
				return q
			},
			val:       0,
			timeout:   time.Second,
			wantData:  []int{0, 0, 0},
			wantSlice: []int{0, 0, 0},
			wantLen:   3,
			wantTail:  0,
			wantHead:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			q := tc.q()
			err := q.InsertLast(ctx, tc.val)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantData, q.data)
			assert.Equal(t, tc.wantSlice, q.AsSlice())
			assert.Equal(t, tc.wantLen, q.Len())
			assert.Equal(t, tc.wantHead, q.front)
			assert.Equal(t, tc.wantTail, q.last)
		})
	}

	t.Run("enqueue timeout", func(t *testing.T) {
		q := ConcurrentConstructor[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := q.InsertLast(ctx, 123)
		require.NoError(t, err)
		err = q.InsertLast(ctx, 234)
		require.NoError(t, err)
		err = q.InsertLast(ctx, 345)
		require.NoError(t, err)
		err = q.InsertLast(ctx, 456)
		require.Equal(t, context.DeadlineExceeded, err)
	})

	// 入队阻塞，而后出队，于是入队成功
	t.Run("enqueue blocking and dequeue", func(t *testing.T) {
		q := ConcurrentConstructor[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		go func() {
			time.Sleep(time.Millisecond * 100)
			val, err := q.GetFront(ctx)
			require.NoError(t, err)
			require.Equal(t, 123, val)
		}()
		err := q.InsertLast(ctx, 123)
		require.NoError(t, err)
		err = q.InsertLast(ctx, 234)
		require.NoError(t, err)
		err = q.InsertLast(ctx, 345)
		require.NoError(t, err)
		err = q.InsertLast(ctx, 456)
		require.NoError(t, err)
	})
}

func TestConcurrentQueue_GetFront(t *testing.T) {
	testCases := []struct {
		name      string
		q         func() *concurrentDeque[int]
		val       int
		timeout   time.Duration
		wantErr   error
		wantVal   int
		wantData  []int
		wantSlice []int
		wantLen   int
		wantHead  int
		wantTail  int
	}{
		{
			name: "dequeued",
			q: func() *concurrentDeque[int] {
				q := ConcurrentConstructor[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 234)
				require.NoError(t, err)
				return q
			},
			wantVal:   123,
			timeout:   time.Second,
			wantData:  []int{0, 234, 0},
			wantSlice: []int{234},
			wantLen:   1,
			wantHead:  1,
			wantTail:  2,
		},
		{
			name: "invalid context",
			q: func() *concurrentDeque[int] {
				q := ConcurrentConstructor[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 234)
				require.NoError(t, err)
				return q
			},
			wantErr:   context.DeadlineExceeded,
			timeout:   -time.Second,
			wantData:  []int{123, 234, 0},
			wantSlice: []int{123, 234},
			wantLen:   2,
			wantHead:  0,
			wantTail:  2,
		},
		{
			name: "dequeue and empty first",
			q: func() *concurrentDeque[int] {
				q := ConcurrentConstructor[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				return q
			},
			wantVal:   123,
			timeout:   time.Second,
			wantData:  []int{0, 0, 0},
			wantSlice: []int{},
			wantLen:   0,
			wantHead:  1,
			wantTail:  1,
		},
		{
			name: "dequeue and empty middle",
			q: func() *concurrentDeque[int] {
				q := ConcurrentConstructor[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 234)
				require.NoError(t, err)
				val, err := q.GetFront(ctx)
				require.NoError(t, err)
				require.Equal(t, 123, val)
				return q
			},
			wantVal:   234,
			timeout:   time.Second,
			wantData:  []int{0, 0, 0},
			wantSlice: []int{},
			wantLen:   0,
			wantHead:  2,
			wantTail:  2,
		},
		{
			name: "dequeue and empty last",
			q: func() *concurrentDeque[int] {
				q := ConcurrentConstructor[int](3)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err := q.InsertLast(ctx, 123)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 234)
				require.NoError(t, err)
				err = q.InsertLast(ctx, 345)
				require.NoError(t, err)
				val, err := q.GetFront(ctx)
				require.NoError(t, err)
				require.Equal(t, 123, val)
				val, err = q.GetFront(ctx)
				require.NoError(t, err)
				require.Equal(t, 234, val)
				return q
			},
			wantVal:   345,
			timeout:   time.Second,
			wantData:  []int{0, 0, 0},
			wantSlice: []int{},
			wantLen:   0,
			wantHead:  0,
			wantTail:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()
			q := tc.q()
			val, err := q.GetFront(ctx)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantVal, val)
			assert.Equal(t, tc.wantData, q.data)
			assert.Equal(t, tc.wantSlice, q.AsSlice())
			assert.Equal(t, tc.wantLen, q.Len())
			assert.Equal(t, tc.wantHead, q.front)
			assert.Equal(t, tc.wantTail, q.last)
		})
	}

	t.Run("dequeue timeout", func(t *testing.T) {
		q := ConcurrentConstructor[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		val, err := q.GetFront(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
		require.Equal(t, 0, val)
	})

	// 出队阻塞，然后入队，然后出队成功
	t.Run("dequeue blocking and enqueue", func(t *testing.T) {
		q := ConcurrentConstructor[int](3)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		go func() {
			time.Sleep(time.Millisecond * 100)
			err := q.InsertLast(ctx, 123)
			require.NoError(t, err)
		}()
		val, err := q.GetFront(ctx)
		require.NoError(t, err)
		require.Equal(t, 123, val)
	})
}

func TestConcurrentQueue(t *testing.T) {
	// 并发测试，只是测试有没有死锁之类的问题
	// 先进先出这个特性依赖于其它单元测试
	// 也依赖于代码审查
	q := ConcurrentConstructor[int](100)
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			val := rand.Int()
			err := q.InsertLast(ctx, val)
			cancel()
			require.NoError(t, err)
		}()
	}
	go func() {
		for i := 0; i < 1000; i++ {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				_, err := q.GetFront(ctx)
				cancel()
				require.NoError(t, err)
				wg.Done()
			}()
		}
	}()
	wg.Wait()
}

func ExampleConcurrentConstructor() {
	q := ConcurrentConstructor[int](10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = q.InsertLast(ctx, 22)
	val, err := q.GetFront(ctx)
	// 这是例子，实际中你不需要写得那么复杂
	switch err {
	case context.Canceled:
		// 有人主动取消了，即调用了 cancel 方法。在这个例子里不会出现这个情况
	case context.DeadlineExceeded:
		// 超时了
	case nil:
		fmt.Println(val)
	default:
		// 其它乱七八糟的
	}
	// Output:
	// 22
}
