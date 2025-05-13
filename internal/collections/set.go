package collections

// Set is a generic set implementation that stores unique elements
type Set[T comparable] struct {
	elements map[T]struct{}
}

// NewSet creates a new empty set
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		elements: make(map[T]struct{}),
	}
}

// Add adds an element to the set
func (s *Set[T]) Add(element T) {
	s.elements[element] = struct{}{}
}

// Remove removes an element from the set
func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}

// Contains checks if an element exists in the set
func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

// Size returns the number of elements in the set
func (s *Set[T]) Size() int {
	return len(s.elements)
}

// Clear removes all elements from the set
func (s *Set[T]) Clear() {
	s.elements = make(map[T]struct{})
}

// ToSlice converts the set to a slice
func (s *Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s.elements))
	for element := range s.elements {
		result = append(result, element)
	}
	return result
}

// ForEach executes a function for each element in the set.
// Note: The iteration order over the set is unspecified and may vary.
func (s *Set[T]) ForEach(f func(T)) {
	for element := range s.elements {
		f(element)
	}
}
