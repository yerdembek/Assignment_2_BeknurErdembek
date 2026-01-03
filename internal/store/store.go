package store

import "sync"

type Store[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func NewStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		m: make(map[K]V),
	}
}

//s := store.NewStore[string, string]() <- I will use it in code in that way

func (s *Store[K, V]) Set(k K, v V) {
	s.mu.Lock()
	s.m[k] = v
	s.mu.Unlock()
}

func (s *Store[K, V]) Get(k K) (v V, ok bool) {
	s.mu.RLock()
	v, ok = s.m[k]
	s.mu.RUnlock()
	return
}

func (s *Store[K, V]) Delete(k K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.m[k]; !ok {
		return false
	}
	delete(s.m, k)
	return true
}

func (s *Store[K, V]) Snapshot() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copyMap := make(map[K]V, len(s.m))
	for k, v := range s.m {
		copyMap[k] = v
	}
	return copyMap
}
