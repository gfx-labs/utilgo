package generic

import "sync"

type Set[K comparable] struct {
	inner Map[K, struct{}]

	mu sync.RWMutex
}

func (s *Set[K]) Store(k K) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.inner.Store(k, struct{}{})
}

func (s *Set[K]) Has(k K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.inner.Load(k)
	return ok
}

// runs function fn if not there and adds k to set. bool returns true if it was new
func (s *Set[K]) HasOrDo(k K, fn func(K)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.inner.Load(k)
	if !ok {
		s.inner.Store(k, struct{}{})
		fn(k)
		return true
	}
	return false
}
