package utils

import (
	"sync"
)

// SyncMap 是一个线程安全的泛型映射
type SyncMap[K comparable, V any] struct {
	mu    sync.RWMutex
	store map[K]V
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.store == nil {
		m.store = make(map[K]V)
	}
	m.store[key] = value
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.store[key]
	return
}

func (m *SyncMap[K, V]) Get(key K) (value V) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value = m.store[key]
	return
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.store) == 0 {
		return
	}
	delete(m.store, key)
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if actual, loaded = m.store[key]; loaded {
		return
	}
	if m.store == nil {
		m.store = make(map[K]V)
	}
	m.store[key] = value
	actual = value
	return
}

func (m *SyncMap[K, V]) LoadOrStoreFunc(key K, load func() (V, error)) (actual V, err error) {
	if value, ok := m.Load(key); ok {
		return value, nil
	}
	if actual, err = load(); err != nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if value, ok := m.store[key]; ok {
		return value, nil
	}
	if m.store == nil {
		m.store = make(map[K]V)
	}
	m.store[key] = actual
	return
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.store {
		if !f(k, v) {
			break
		}
	}
}
