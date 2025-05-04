package errors

import (
	"fmt"
)

// Service-specific errors
var (
	ErrCapacityExceeded = fmt.Errorf("%w: capacity exceeded", ErrInvalidInput)
	ErrMatchingFailed   = fmt.Errorf("%w: ride matching failed", ErrInternal)
)

// InvalidTimeRange returns an error for invalid time ranges
func InvalidTimeRange() error {
	return fmt.Errorf("end time must be after start time: %w", ErrInvalidInput)
}

// CapacityExceeded returns a formatted capacity exceeded error
func CapacityExceeded(offerID string, capacity int) error {
	return fmt.Errorf("driver offer %s has reached maximum capacity of %d: %w",
		offerID, capacity, ErrCapacityExceeded)
}
