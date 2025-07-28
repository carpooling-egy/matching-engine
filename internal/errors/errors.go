package errors

import (
	"errors"
	"fmt"
)

// Standard error types
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInternal     = errors.New("internal error")
	ErrUnauthorized = errors.New("unauthorized access")
)

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsInvalidInput checks if an error is an invalid input error
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsInternal checks if an error is an internal error
func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}

// IsUnauthorized checks if an error is an authorization error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// GetErrorMessage extracts a user-friendly error message
func GetErrorMessage(err error) string {
	// You can customize error messages based on type
	if IsNotFound(err) {
		return "The requested resource was not found"
	}
	if IsInvalidInput(err) {
		return "Invalid input provided"
	}
	if IsInternal(err) {
		return "An internal error occurred"
	}
	return err.Error()
}

// Wrap adds context to an error while preserving its type
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// New creates a new error with the given message
func New(message string) error {
	return errors.New(message)
}

// Unwrap is a convenience function that calls errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
