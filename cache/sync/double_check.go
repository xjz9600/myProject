package sync

import "sync"

type safeMap[k comparable, v any] struct {
	data map[k]v
	sync sync.RWMutex
}

func (s *safeMap[k, v]) Get(key k) (any, bool) {
	s.sync.RLock()
	defer s.sync.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *safeMap[k, v]) Put(key k, val v) {
	s.sync.Lock()
	defer s.sync.Unlock()
	s.data[key] = val
}

func (s *safeMap[k, v]) LoadOrStore(key k, val v) (any, bool) {
	s.sync.RLock()
	oldVal, ok := s.data[key]
	s.sync.RUnlock()
	if ok {
		return oldVal, ok
	}
	s.sync.Lock()
	defer s.sync.Unlock()
	oldVal, ok = s.data[key]
	if ok {
		return oldVal, ok
	}
	s.data[key] = val
	return val, false
}
