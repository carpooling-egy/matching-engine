// Package collections SyncMap provides a generic, thread-safe syncMap implementation
package collections

import (
	"fmt"
	"sync"
)

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

// Range processes key-value pairs, stopping at the first error.
// Returns the first error encountered, or nil if none.
// Safely handles type assertions to prevent panics.
func (c *SyncMap[K, V]) Range(f func(K, V) error) error {
	var err error
	c.store.Range(func(key, value interface{}) bool {
		k, ok1 := key.(K)
		if !ok1 {
			err = fmt.Errorf("type assertion failed for key: %v", key)
			return false
		}

		v, ok2 := value.(V)
		if !ok2 {
			err = fmt.Errorf("type assertion failed for value: %v", value)
			return false
		}

		err = f(k, v)
		return err == nil // Stop if error occurs
	})
	return err
}

// ForEach processes all key-value pairs, collecting any errors.
// Returns a slice of all errors encountered, or empty if none.
// Safely handles type assertions to prevent panics.
func (c *SyncMap[K, V]) ForEach(f func(K, V) error) []error {
	var errors []error
	c.store.Range(func(key, value interface{}) bool {
		k, ok1 := key.(K)
		if !ok1 {
			errors = append(errors, fmt.Errorf("type assertion failed for key: %v", key))
			return true // Continue despite error
		}

		v, ok2 := value.(V)
		if !ok2 {
			errors = append(errors, fmt.Errorf("type assertion failed for value: %v", value))
			return true // Continue despite error
		}

		if err := f(k, v); err != nil {
			errors = append(errors, err)
		}
		return true // Always continue
	})
	return errors
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
