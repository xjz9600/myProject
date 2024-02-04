package channel

import (
	"errors"
	"sync"
)

type Broker struct {
	msg   []chan Msg
	mutex sync.RWMutex
}

func (b *Broker) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	oldMsg := b.msg
	b.msg = nil
	for _, om := range oldMsg {
		close(om)
	}
	return nil
}

func (b *Broker) Send(msg Msg) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for _, m := range b.msg {
		select {
		case m <- msg:
		default:
			return errors.New("队列已满")
		}
	}
	return nil
}

func (b *Broker) subscribe(capacity int) (<-chan Msg, error) {
	ms := make(chan Msg, capacity)
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.msg = append(b.msg, ms)
	return ms, nil
}

type Msg struct {
	content string
}
