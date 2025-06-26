package collections

import "sync"

// Set is a generic set implementation that stores unique elements
type Set[T comparable] struct {
	elements sync.Map
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{}
}

// Add adds an element to the set
func (s *Set[T]) Add(element T) {
	s.elements.Store(element, struct{}{})
}

// Remove removes an element from the set
func (s *Set[T]) Remove(element T) {
	s.elements.Delete(element)
}

// Contains checks if an element exists in the set
func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements.Load(element)
	return exists
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	count := 0
	s.elements.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Note: This is thread-safe, but if other goroutines are concurrently iterating
// over or accessing the old map, they may still see stale data briefly, since the old map
// is not explicitly locked or cleared.

// Clear removes all elements from the set
func (s *Set[T]) Clear() {
	s.elements = sync.Map{}
}

// ToSlice converts the set to a slice
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0)
	s.elements.Range(func(key, value interface{}) bool {
		result = append(result, key.(T))
		return true
	})
	return result
}

// ForEach executes a function for each element in the set.
// Note: The iteration order over the set is unspecified and may vary.
func (s *Set[T]) ForEach(f func(T)) {
	s.elements.Range(func(key, value interface{}) bool {
		f(key.(T))
		return true
	})
}
