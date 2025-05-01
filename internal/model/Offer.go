package model

import "github.com/google/uuid"

// Offer represents a service provider's offer
type Offer struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	MatchedRequests []MatchedRequest
}

// NewOffer creates a new Offer
func NewOffer(id, userID uuid.UUID, matchedRequests []MatchedRequest) *Offer {
	if matchedRequests == nil {
		matchedRequests = make([]MatchedRequest, 0)
	}
	return &Offer{
		ID:              id,
		UserID:          userID,
		MatchedRequests: matchedRequests,
	}
}
