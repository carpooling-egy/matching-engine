package model

import (
	"github.com/google/uuid"
	"time"
)

// Offer represents a service provider's offer
type Offer struct {
	id              uuid.UUID
	userID          uuid.UUID
	source          Coordinate
	destination     Coordinate
	detourTime      time.Duration
	departureTime   time.Time
	matchedRequests []MatchedRequest
	preference      Preference
	path            []*Point
}

// NewOffer creates a new Offer
func NewOffer(id, userID uuid.UUID, source, destination Coordinate, detourTime time.Duration, departureTime time.Time, matchedRequests []MatchedRequest, preference Preference, path []*Point) *Offer {
	if matchedRequests == nil {
		matchedRequests = make([]MatchedRequest, 0)
	}
	if path == nil {
		path = make([]*Point, 0)
	}
	return &Offer{
		id:              id,
		userID:          userID,
		source:          source,
		destination:     destination,
		detourTime:      detourTime,
		departureTime:   departureTime,
		matchedRequests: matchedRequests,
		preference:      preference,
		path:            path,
	}
}
