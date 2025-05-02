package model

import (
	"github.com/google/uuid"
	"time"
)

// Request represents a user's request for a service
type Request struct {
	id                    uuid.UUID
	userID                uuid.UUID
	source                Coordinate
	destination           Coordinate
	earliestDepartureTime time.Time
	latestArrivalTime     time.Time
	maxWalkingTime        time.Duration
	preference            Preference
	numberOfRiders        int
}

// NewRequest creates a new Request
func NewRequest(id, userID uuid.UUID, source, destination Coordinate, earliestDepartureTime, latestArrivalTime time.Time, maxWalkingTime time.Duration, preference Preference, numberOfRiders int) *Request {
	return &Request{
		id:                    id,
		userID:                userID,
		source:                source,
		destination:           destination,
		earliestDepartureTime: earliestDepartureTime,
		latestArrivalTime:     latestArrivalTime,
		maxWalkingTime:        maxWalkingTime,
		preference:            preference,
		numberOfRiders:        numberOfRiders,
	}
}
