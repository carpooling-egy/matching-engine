package model

// MatchingResult represents the result of a matching operation
type MatchingResult struct {
	userID                  string
	offerID                 string
	assignedMatchedRequests []*MatchedRequest
	newPath                 []*PathPoint
}

// NewMatchingResult creates a new MatchingResult
func NewMatchingResult(userID, offerID string, assignedMatchedRequests []*MatchedRequest, newPath []*PathPoint) *MatchingResult {
	if assignedMatchedRequests == nil {
		assignedMatchedRequests = make([]*MatchedRequest, 0)
	}
	if newPath == nil {
		newPath = make([]*PathPoint, 0)
	}
	return &MatchingResult{
		userID:                  userID,
		offerID:                 offerID,
		assignedMatchedRequests: assignedMatchedRequests,
		newPath:                 newPath,
	}
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
func (mr *MatchingResult) AssignedMatchedRequests() []*MatchedRequest {
	return mr.assignedMatchedRequests
}

// SetAssignedMatchedRequests sets the assigned matched requests
func (mr *MatchingResult) SetAssignedMatchedRequests(requests []*MatchedRequest) {
	mr.assignedMatchedRequests = requests
}

// NewPath returns the new path
func (mr *MatchingResult) NewPath() []*PathPoint {
	return mr.newPath
}

// SetNewPath sets the new path
func (mr *MatchingResult) SetNewPath(path []*PathPoint) {
	mr.newPath = path
}
