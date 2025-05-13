// Package collections SyncMap provides a generic, thread-safe syncMap implementation
package collections

import "sync"

// SyncMap is a generic thread-safe syncMap generic structure
type SyncMap[K comparable, V any] struct {
	store sync.Map
}

// New creates a new syncMap instance
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
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

// Range iterates over the SyncMap and calls the function for each key-value pair.
// If the function returns false, the iteration stops.
func (c *SyncMap[K, V]) Range(f func(K, V) bool) {
	c.store.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

// ForEach executes a function for each key-value pair in the SyncMap.
// Note: The iteration order over the SyncMap is unspecified and may vary.
func (c *SyncMap[K, V]) ForEach(f func(K, V)) {
	c.store.Range(func(key, value interface{}) bool {
		f(key.(K), value.(V))
		return true
	})
}

// Contains checks if a key exists in the SyncMap
func (c *SyncMap[K, V]) Contains(key K) bool {
	_, exists := c.store.Load(key)
	return exists
}

// Size returns the number of key-value pairs in the SyncMap
func (c *SyncMap[K, V]) Size() int {
	count := 0
	c.store.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Clear removes all items from the syncMap
func (c *SyncMap[K, V]) Clear() {
	c.store = sync.Map{}
}
