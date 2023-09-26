package leetcode

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

// ConcurrentLinkedQueue 无界并发安全队列
type ConcurrentLinkedQueue[T any] struct {
	header unsafe.Pointer
	tail   unsafe.Pointer
}

func NewConcurrentLinkedQueue[T any]() *ConcurrentLinkedQueue[T] {
	head := &node[T]{}
	ptr := unsafe.Pointer(head)
	return &ConcurrentLinkedQueue[T]{
		header: ptr,
		tail:   ptr,
	}
}

func (c *ConcurrentLinkedQueue[T]) Enqueue(t T) error {
	newNode := &node[T]{val: t}
	newPtr := unsafe.Pointer(newNode)
	for {
		tailPtr := atomic.LoadPointer(&c.tail)
		tail := (*node[T])(tailPtr)
		tailNext := atomic.LoadPointer(&tail.next)
		if tailNext != nil {
			// 证明已经有人修改过了
			continue
		}
		if atomic.CompareAndSwapPointer(&tail.next, tailNext, newPtr) {
			if atomic.CompareAndSwapPointer(&c.tail, tailPtr, newPtr) {
				return nil
			}
		}
	}
}

func (c *ConcurrentLinkedQueue[T]) Dequeue() (T, error) {
	for {
		headPtr := atomic.LoadPointer(&c.header)
		head := (*node[T])(headPtr)
		tailPtr := atomic.LoadPointer(&c.tail)
		if headPtr == tailPtr {
			var t T
			return t, errors.New("没有数据")
		}
		headNextPtr := atomic.LoadPointer(&head.next)
		if atomic.CompareAndSwapPointer(&c.header, headPtr, headNextPtr) {
			return (*node[T])(headNextPtr).val, nil
		}
	}
}

type node[T any] struct {
	val  T
	next unsafe.Pointer
}
