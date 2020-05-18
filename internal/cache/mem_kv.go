package cache

import "sync"

type MemKv struct {
	kvs   map[interface{}]interface{}
	mutex sync.RWMutex
}

func NewMemKv() *MemKv {
	return &MemKv{
		kvs: make(map[interface{}]interface{}),
	}
}

func (m *MemKv) Get(key interface{}) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	val, ok := m.kvs[key]
	return val, ok
}

func (m *MemKv) Set(key interface{}, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.kvs[key] = value
}