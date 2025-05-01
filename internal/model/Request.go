package model

import "github.com/google/uuid"

// Request represents a user's request for a service
type Request struct {
	ID uuid.UUID
	// UserID is the ID of the user who created the request
	UserID      uuid.UUID
	Source      Coordinate
	Destination Coordinate
}

// NewRequest creates a new Request
func NewRequest(id, userID uuid.UUID, source, destination Coordinate) *Request {
	return &Request{
		ID:          id,
		UserID:      userID,
		Source:      source,
		Destination: destination,
	}
}
