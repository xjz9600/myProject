package queue

import (
	"context"
)

type SliceQueue[T any] struct {
	data []T
}

func NewSliceQueue[T any]() *SliceQueue[T] {
	return &SliceQueue[T]{
		data: make([]T, 0),
	}
}

func (s *SliceQueue[T]) In(ctx context.Context, t T) error {
	s.data = append(s.data, t)
	return nil
}

func (s *SliceQueue[T]) Out(ctx context.Context) (T, error) {
	if len(s.data) == 0 {
		return *new(T), nil
	}
	re := s.data[0]
	s.data = s.data[1:]
	return re, nil
}
