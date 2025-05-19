// Package collections SyncMap provides a generic, thread-safe syncMap implementation
package collections

import (
	"fmt"
	"sync"
)

// SyncMap is a generic thread-safe syncMap generic structure
type SyncMap[K comparable, V any] struct {
	store sync.Map
	size  int
	mu    sync.Mutex
}

// NewSyncMap creates a new SyncMap instance
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		size: 0,
	}
}

// Set stores a value in the syncMap with the given key
func (sm *SyncMap[K, V]) Set(key K, value V) {
	sm.store.Store(key, value)
	sm.mu.Lock()
	sm.size++
	sm.mu.Unlock()
}

// Delete removes a value from the syncMap by key
func (sm *SyncMap[K, V]) Delete(key K) {
	sm.store.Delete(key)
	sm.mu.Lock()
	sm.size--
	sm.mu.Unlock()
}

// Size returns the number of key-value pairs in the SyncMap
func (sm *SyncMap[K, V]) Size() int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.size
}

// Get retrieves a value from the syncMap by key
func (sm *SyncMap[K, V]) Get(key K) (V, bool) {
	val, ok := sm.store.Load(key)
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

// Range processes key-value pairs, stopping at the first error.
// Returns the first error encountered, or nil if none.
// Safely handles type assertions to prevent panics.
func (sm *SyncMap[K, V]) Range(f func(K, V) error) error {
	var err error
	sm.store.Range(func(key, value any) bool {
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
func (sm *SyncMap[K, V]) ForEach(f func(K, V) error) []error {
	var errors []error
	sm.store.Range(func(key, value any) bool {
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
func (sm *SyncMap[K, V]) Contains(key K) bool {
	_, exists := sm.store.Load(key)
	return exists
}

// Clear removes all items from the syncMap
func (sm *SyncMap[K, V]) Clear() {
	sm.store = sync.Map{}
}
