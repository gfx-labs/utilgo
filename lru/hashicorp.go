package lru

import (
	"git.tuxpa.in/a/lambda"
	hashicorp "github.com/hashicorp/golang-lru"
)

type HashiCorp[K comparable, V any] struct {
	c hashicorp.Cache
}

// Adds a value to the cache, returns true if an eviction occurred and
// updates the "recently used"-ness of the key.
func (h *HashiCorp[K, V]) Add(key K, value V) bool {
	return h.c.Add(key, value)
}

// Returns key's value from the cache and
// updates the "recently used"-ness of the key. #value, isFound
func (h *HashiCorp[K, V]) Get(key K) (value V, ok bool) {
	val, ok := h.c.Get(key)
	if !ok {
		return *new(V), false
	}
	if resp, ok := val.(V); ok {
		return resp, true
	}
	return *new(V), false
}

// Checks if a key exists in cache without updating the recent-ness.
func (h *HashiCorp[K, V]) Contains(key K) (ok bool) {
	return h.c.Contains(key)
}

// Returns key's value without updating the "recently used"-ness of the key.
func (h *HashiCorp[K, V]) Peek(key K) (value V, ok bool) {
	val, ok := h.c.Peek(key)
	if !ok {
		return *new(V), false
	}
	if resp, ok := val.(V); ok {
		return resp, true
	}
	return *new(V), false
}

// Removes a key from the cache.
func (h *HashiCorp[K, V]) Remove(key K) bool {
	return h.Remove(key)
}

// Removes the oldest entry from cache.
func (h *HashiCorp[K, V]) RemoveOldest() (K, V, bool) {
	k, v, ok := h.c.RemoveOldest()
	if ok {
		return k.(K), v.(V), ok
	}
	return *new(K), *new(V), false
}

// Returns the oldest entry from the cache. #key, value, isFound
func (h *HashiCorp[K, V]) GetOldest() (K, V, bool) {
	k, v, ok := h.c.GetOldest()
	if ok {
		return k.(K), v.(V), ok
	}
	return *new(K), *new(V), false
}

// Returns a slice of the keys in the cache, from oldest to newest.
func (h *HashiCorp[K, V]) Keys() []K {
	return lambda.MapV(h.c.Keys(), func(v any) K {
		return v.(K)
	})
}

// Returns the number of items in the cache.
func (h *HashiCorp[K, V]) Len() int {
	return h.c.Len()
}

// Clears all cache entries.
func (h *HashiCorp[K, V]) Purge() {
	h.c.Purge()
}

// Resizes cache, returning number evicted
func (h *HashiCorp[K, V]) Resize(i int) int {
	return h.c.Resize(i)
}
