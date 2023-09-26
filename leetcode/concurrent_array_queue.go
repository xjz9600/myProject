package leetcode

import (
	"context"
	"golang.org/x/sync/semaphore"
	"sync"
)

// ConcurrentArrayBlockingQueue 有界并发阻塞队列
type ConcurrentArrayBlockingQueue[T any] struct {
	data  []T
	mutex sync.RWMutex

	// 队头元素下标
	head int
	// 队尾元素下标
	tail int
	// 包含多少个元素
	count int

	enqueueCap *semaphore.Weighted
	dequeueCap *semaphore.Weighted

	// zero 不能作为返回值返回，防止用户篡改
	zero T
}

// NewConcurrentArrayBlockingQueue 创建一个有界阻塞队列
// 容量会在最开始的时候就初始化好
// capacity 必须为正数
func NewConcurrentArrayBlockingQueue[T any](capacity int) *ConcurrentArrayBlockingQueue[T] {
	data := make([]T, capacity)
	enqueueCap := semaphore.NewWeighted(int64(capacity))
	dequeueCap := semaphore.NewWeighted(int64(capacity))
	dequeueCap.Acquire(context.Background(), int64(capacity))
	return &ConcurrentArrayBlockingQueue[T]{
		data:       data,
		enqueueCap: enqueueCap,
		dequeueCap: dequeueCap,
	}
}

// Enqueue 入队
// 通过sema来控制容量、超时、阻塞问题
func (c *ConcurrentArrayBlockingQueue[T]) Enqueue(ctx context.Context, t T) error {
	// 上来判断超时
	if ctx.Err() != nil {
		return ctx.Err()
	}
	err := c.enqueueCap.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	// 再次判断超时，防止释放后已经超时了
	if ctx.Err() != nil {
		// 超时后，释放一个位置
		c.enqueueCap.Release(1)
		return ctx.Err()
	}
	c.data[c.tail] = t
	c.count++
	c.tail = (c.tail + 1) % len(c.data)
	c.mutex.Unlock()
	// 告诉出口可以取出一个数据了
	c.dequeueCap.Release(1)
	return nil
}

// Dequeue 出队
// 通过sema来控制容量、超时、阻塞问题
func (c *ConcurrentArrayBlockingQueue[T]) Dequeue(ctx context.Context) (T, error) {
	// 上来判断超时
	if ctx.Err() != nil {
		return c.zero, ctx.Err()
	}
	err := c.dequeueCap.Acquire(ctx, 1)
	if err != nil {
		return c.zero, err
	}
	c.mutex.Lock()
	if ctx.Err() != nil {
		// 超时后，释放一个位置
		c.dequeueCap.Release(1)
		return c.zero, ctx.Err()
	}
	res := c.data[c.head]
	// 释放内存
	c.data[c.head] = c.zero
	c.count--
	c.head = (c.head + 1 + len(c.data)) % len(c.data)
	c.mutex.Unlock()
	c.enqueueCap.Release(1)
	return res, nil
}

func (c *ConcurrentArrayBlockingQueue[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

func (c *ConcurrentArrayBlockingQueue[T]) AsSlice() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	res := make([]T, 0, c.count)
	cnt := 0
	capacity := cap(c.data)
	for cnt < c.count {
		index := (c.head + cnt) % capacity
		res = append(res, c.data[index])
		cnt++
	}
	return res
}
