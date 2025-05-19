package model

import (
	"fmt"
	"matching-engine/internal/errors"
)

// MatchingResult represents the result of a matching operation
type MatchingResult struct {
	userID                  string
	offerID                 string
	assignedMatchedRequests []*Request
	newPath                 []PathPoint
	currentNumberOfRequests int
}

// NewMatchingResult creates a new MatchingResult
func NewMatchingResult(userID, offerID string, assignedMatchedRequests []*Request, newPath []PathPoint, currentNumberOfRequests int) *MatchingResult {
	if assignedMatchedRequests == nil {
		assignedMatchedRequests = make([]*Request, 0)
	}
	if newPath == nil {
		newPath = make([]PathPoint, 0)
	}
	return &MatchingResult{
		userID:                  userID,
		offerID:                 offerID,
		assignedMatchedRequests: assignedMatchedRequests,
		newPath:                 newPath,
		currentNumberOfRequests: currentNumberOfRequests,
	}
}

func NewMatchingResultFromOfferNode(node *OfferNode) (*MatchingResult, error) {
	if node == nil {
		return nil, fmt.Errorf(errors.ErrNilOfferNode)
	}
	if node.NewlyAssignedMatchedRequests() == nil {
		return nil, fmt.Errorf(errors.ErrNilMatchedRequests)
	}
	if len(node.NewlyAssignedMatchedRequests()) == 0 {
		return nil, fmt.Errorf(errors.ErrEmptyMatchedRequests)
	}
	currentNumberOfRequests := len(node.Offer().matchedRequests) + len(node.NewlyAssignedMatchedRequests())
	return &MatchingResult{
		userID:                  node.Offer().UserID(),
		offerID:                 node.Offer().ID(),
		assignedMatchedRequests: node.NewlyAssignedMatchedRequests(),
		newPath:                 node.Offer().Path(),
		currentNumberOfRequests: currentNumberOfRequests,
	}, nil
}

// UserID returns the user ID
func (mr *MatchingResult) UserID() string {
	return mr.userID
}

// SetUserID sets the user ID
func (mr *MatchingResult) SetUserID(userID string) {
	mr.userID = userID
}

// OfferID returns the offer ID
func (mr *MatchingResult) OfferID() string {
	return mr.offerID
}

// SetOfferID sets the offer ID
func (mr *MatchingResult) SetOfferID(offerID string) {
	mr.offerID = offerID
}

// AssignedMatchedRequests returns the assigned matched requests
func (mr *MatchingResult) AssignedMatchedRequests() []*Request {
	return mr.assignedMatchedRequests
}

// SetAssignedMatchedRequests sets the assigned matched requests
func (mr *MatchingResult) SetAssignedMatchedRequests(requests []*Request) {
	mr.assignedMatchedRequests = requests
}

// NewPath returns the new path
func (mr *MatchingResult) NewPath() []PathPoint {
	return mr.newPath
}

// SetNewPath sets the new path
func (mr *MatchingResult) SetNewPath(path []PathPoint) {
	mr.newPath = path
}

// CurrentNumberOfRequests returns the current number of requests
func (mr *MatchingResult) CurrentNumberOfRequests() int {
	return mr.currentNumberOfRequests
}

// SetCurrentNumberOfRequests sets the current number of requests
func (mr *MatchingResult) SetCurrentNumberOfRequests(count int) {
	mr.currentNumberOfRequests = count
}
