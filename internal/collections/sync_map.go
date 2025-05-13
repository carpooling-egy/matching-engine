// Package collections SyncMap provides a generic, thread-safe syncMap implementation
package collections

import "sync"

// SyncMap is a generic thread-safe syncMap generic structure
type SyncMap[K comparable, V any] struct {
	store sync.Map
}

// New creates a new syncMap instance
func New[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// Set stores a value in the syncMap with the given key
func (c *SyncMap[K, V]) Set(key K, value V) {
	c.store.Store(key, value)
}

// Get retrieves a value from the syncMap by key
func (c *SyncMap[K, V]) Get(key K) (V, bool) {
	val, ok := c.store.Load(key)
	if !ok {
		var zero V
		return zero, false
	}

	// Type assertion to ensure the value is of the expected type
	value, ok := val.(V)
	if !ok {
		var zero V
		return zero, false
	}

	return value, true
}

// Delete removes a value from the syncMap by key
func (c *SyncMap[K, V]) Delete(key K) {
	c.store.Delete(key)
}

// Clear removes all items from the syncMap
func (c *SyncMap[K, V]) Clear() {
	c.store = sync.Map{}
}
