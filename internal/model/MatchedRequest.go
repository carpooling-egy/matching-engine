package model

import "github.com/google/uuid"

// MatchedRequest represents a request that has been matched with an offer
type MatchedRequest struct {
	offerID   uuid.UUID
	requestID uuid.UUID
	pickup    *Point
	dropoff   *Point
}

// NewMatchedRequest creates a new MatchedRequest
func NewMatchedRequest(offerID, requestID uuid.UUID, pickup, dropoff *Point) *MatchedRequest {
	return &MatchedRequest{
		offerID:   offerID,
		requestID: requestID,
		pickup:    pickup,
		dropoff:   dropoff,
	}
}
