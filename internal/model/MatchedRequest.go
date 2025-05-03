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

// GetOfferID returns the offer ID
func (mr *MatchedRequest) GetOfferID() uuid.UUID {
	return mr.offerID
}

// SetOfferID sets the offer ID
func (mr *MatchedRequest) SetOfferID(offerID uuid.UUID) {
	mr.offerID = offerID
}

// GetRequestID returns the request ID
func (mr *MatchedRequest) GetRequestID() uuid.UUID {
	return mr.requestID
}

// SetRequestID sets the request ID
func (mr *MatchedRequest) SetRequestID(requestID uuid.UUID) {
	mr.requestID = requestID
}

// GetPickup returns the pickup point
func (mr *MatchedRequest) GetPickup() *Point {
	return mr.pickup
}

// SetPickup sets the pickup point
func (mr *MatchedRequest) SetPickup(pickup *Point) {
	mr.pickup = pickup
}

// GetDropoff returns the dropoff point
func (mr *MatchedRequest) GetDropoff() *Point {
	return mr.dropoff
}

// SetDropoff sets the dropoff point
func (mr *MatchedRequest) SetDropoff(dropoff *Point) {
	mr.dropoff = dropoff
}
