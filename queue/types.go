package queue

import "context"

type Queue[T any] interface {
	In(context.Context, T) error
	Out(context.Context) (T, error)
	//Clear()
	//Size() int
}
