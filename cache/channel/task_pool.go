package channel

import (
	"context"
	"errors"
	"sync"
)

type Task func()

type TaskPool struct {
	tasks  chan Task
	closer chan struct{}
	once   sync.Once
}

func (t *TaskPool) Close() error {
	t.once.Do(func() {
		close(t.tasks)
	})
	return nil
}

func NewTaskPool(numG int, capacity int) *TaskPool {
	res := &TaskPool{
		tasks:  make(chan Task, capacity),
		closer: make(chan struct{}),
	}
	for i := 0; i < numG; i++ {
		go func() {
			select {
			case task := <-res.tasks:
				task()
			case <-res.closer:
				return
			}
		}()
	}
	return res
}

func (t *TaskPool) Submit(ctx context.Context, task Task) error {
	select {
	case t.tasks <- task:
	case <-ctx.Done():
		return errors.New("添加任务超时")
	}
	return nil
}
