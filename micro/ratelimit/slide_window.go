package ratelimit

import (
	"container/list"
	"context"
	"errors"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type SlideWindowLimiter struct {
	list     *list.List
	interval int64
	rate     int
	mutex    sync.Mutex
}

func NewSlideWindowLimiter(interval int64, rate int) *SlideWindowLimiter {
	return &SlideWindowLimiter{
		list:     list.New(),
		interval: interval,
		rate:     rate,
		mutex:    sync.Mutex{},
	}
}

func (s *SlideWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		now := time.Now().UnixNano()
		boundary := now - s.interval
		s.mutex.Lock()
		timeStamp := s.list.Front()
		for timeStamp != nil && timeStamp.Value.(int64) < boundary {
			s.list.Remove(timeStamp)
			timeStamp = s.list.Front()
		}
		length := s.list.Len()
		if length >= s.rate {
			err = errors.New("触发瓶颈了")
			s.mutex.Unlock()
			return
		}
		s.list.PushBack(now)
		s.mutex.Unlock()
		return handler(ctx, req)
	}
}
