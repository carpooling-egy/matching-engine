package model

import "github.com/google/uuid"

// MatchedRequest represents a request that has been matched with an offer
type MatchedRequest struct {
	OfferID   uuid.UUID
	RequestID uuid.UUID
	Pickup    Point
	Dropoff   Point
}

// NewMatchedRequest creates a new MatchedRequest
func NewMatchedRequest(offerID, requestID uuid.UUID, pickup, dropoff Point) *MatchedRequest {
	return &MatchedRequest{
		OfferID:   offerID,
		RequestID: requestID,
		Pickup:    pickup,
		Dropoff:   dropoff,
	}
}
