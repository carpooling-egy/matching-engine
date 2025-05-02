package model

import "github.com/google/uuid"

// MatchingResult represents the result of a matching operation
type MatchingResult struct {
	offerID                 uuid.UUID
	assignedMatchedRequests []MatchedRequest
	newPath                 []Point
}

// NewMatchingResult creates a new MatchingResult
func NewMatchingResult(offerID uuid.UUID, assignedMatchedRequests []MatchedRequest, newPath []Point) *MatchingResult {
	if assignedMatchedRequests == nil {
		assignedMatchedRequests = make([]MatchedRequest, 0)
	}
	if newPath == nil {
		newPath = make([]Point, 0)
	}
	return &MatchingResult{
		offerID:                 offerID,
		assignedMatchedRequests: assignedMatchedRequests,
		newPath:                 newPath,
	}
}
