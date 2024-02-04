package limit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func NewErrKeyNotFound(key string) error {
	return fmt.Errorf("key %s 不存在", key)
}

type MaxMemoryCache struct {
	Cache
	max  int64
	used int64
	LRUCache[[]byte]
	mutex sync.RWMutex
}

func NewMaxMemoryCache(max int64, cache Cache) *MaxMemoryCache {
	res := &MaxMemoryCache{
		Cache:    cache,
		max:      max,
		LRUCache: ConstructorNode[[]byte](),
	}
	res.Cache.OnEvicted(res.onEvicted)
	return res
}

func (m *MaxMemoryCache) Delete(ctx context.Context, key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Cache.Delete(ctx, key)
}

func (m *MaxMemoryCache) LoadAndDelete(ctx context.Context, key string) ([]byte, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.Cache.LoadAndDelete(ctx, key)
}

func (m *MaxMemoryCache) onEvicted(key string, val []byte) {
	if _, ok := m.LRUCache.data[key]; ok {
		node := m.LRUCache.RemoveNode(key)
		m.used = m.used - int64(len(node.value))
	}
}

func (m *MaxMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	_ = m.LRUCache.Get(key)
	return m.Cache.Get(ctx, key)
}

func (m *MaxMemoryCache) Set(ctx context.Context, key string, val []byte,
	expiration time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Cache.LoadAndDelete(ctx, key)
	for m.used+int64(len(val)) > m.max {
		m.Cache.Delete(ctx, m.tail.preNode.key)
	}
	m.used = m.used + int64(len(val))
	m.Put(key, val)
	return m.Cache.Set(ctx, key, val, expiration)
}

func (m *MaxMemoryCache) GetAllKeys() []string {
	var result []string
	pre := m.head.nextNode
	for pre != m.tail {
		result = append(result, pre.key)
		pre = pre.nextNode
	}
	return result
}
