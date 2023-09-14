package leetcode

import (
	"context"
	"sync"
)

type concurrentDeque[T any] struct {
	data  []T
	front int
	last  int
	count int

	mutex   *sync.RWMutex
	dequeue *sync.Cond
	enqueue *sync.Cond

	zero T
}

func ConcurrentConstructor[T any](k int) *concurrentDeque[T] {
	return &concurrentDeque[T]{
		data:    make([]T, k),
		dequeue: sync.NewCond(&sync.Mutex{}),
		enqueue: sync.NewCond(&sync.Mutex{}),
		mutex:   &sync.RWMutex{},
	}
}
func (this *concurrentDeque[T]) InsertLast(ctx context.Context, t T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	for this.IsFull() {
		ch := make(chan struct{})
		go func() {
			this.enqueue.L.Lock()
			this.enqueue.Wait()
			this.enqueue.L.Unlock()
			select {
			case ch <- struct{}{}:
			default:
				this.enqueue.Signal()
			}
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ch:
		}
	}
	this.mutex.Lock()
	this.count++
	this.data[this.last] = t
	this.last = (this.last + 1) % len(this.data)
	this.mutex.Unlock()
	this.dequeue.Signal()
	return nil
}

func (this *concurrentDeque[T]) IsEmpty() bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.count == 0
}

func (this *concurrentDeque[T]) IsFull() bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.count == cap(this.data)
}

func (c *concurrentDeque[T]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.count
}

func (c *concurrentDeque[T]) AsSlice() []T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	res := make([]T, 0, c.count)
	cnt := 0
	capacity := cap(c.data)
	for cnt < c.count {
		index := (c.front + cnt) % capacity
		res = append(res, c.data[index])
		cnt++
	}
	return res
}

func (this *concurrentDeque[T]) GetFront(ctx context.Context) (T, error) {
	if ctx.Err() != nil {
		return this.zero, ctx.Err()
	}
	for this.IsEmpty() {
		ch := make(chan struct{})
		go func() {
			this.dequeue.L.Lock()
			this.dequeue.Wait()
			this.dequeue.L.Unlock()
			select {
			case ch <- struct{}{}:
			default:
				this.dequeue.Signal()
			}
		}()
		select {
		case <-ctx.Done():
			return this.zero, ctx.Err()
		case <-ch:
		}
	}
	this.mutex.Lock()
	data := this.data[this.front]
	this.data[this.front] = this.zero
	this.count--
	this.front = (this.front + 1) % len(this.data)
	this.mutex.Unlock()
	this.enqueue.Signal()
	return data, nil
}
