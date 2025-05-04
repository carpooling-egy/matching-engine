package errors

import (
	"fmt"
)

// Database-specific errors
var (
	ErrDatabase = fmt.Errorf("%w: database error", ErrInternal)
)

// wrapDatabaseError adds database context to any error
func wrapDatabaseError(err error) error {
	return fmt.Errorf("database: %w", err)
}

// NotFound returns a formatted "not found" error
func NotFound(resource, id string) error {
	baseErr := fmt.Errorf("%s not found: %s: %w", resource, id, ErrNotFound)
	return wrapDatabaseError(baseErr)
}

// EmptyID returns an error for empty ID inputs
func EmptyID(resourceType string) error {
	baseErr := fmt.Errorf("%s ID cannot be empty: %w", resourceType, ErrInvalidInput)
	return wrapDatabaseError(baseErr)
}

// InvalidInput returns a formatted validation error
func InvalidInput(message string) error {
	baseErr := fmt.Errorf("%s: %w", message, ErrInvalidInput)
	return wrapDatabaseError(baseErr)
}

// DatabaseError wraps an error with database operation context
func DatabaseError(operation string, err error) error {
	baseErr := fmt.Errorf("%s operation failed: %w: %w", operation, err, ErrDatabase)
	return wrapDatabaseError(baseErr)
}
