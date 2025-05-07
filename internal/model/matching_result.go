package model

// MatchingResult represents the result of a matching operation
type MatchingResult struct {
	userID                  string
	offerID                 string
	assignedMatchedRequests []*MatchedRequest
	newPath                 []*Point
}

// NewMatchingResult creates a new MatchingResult
func NewMatchingResult(userID, offerID string, assignedMatchedRequests []*MatchedRequest, newPath []*Point) *MatchingResult {
	if assignedMatchedRequests == nil {
		assignedMatchedRequests = make([]*MatchedRequest, 0)
	}
	if newPath == nil {
		newPath = make([]*Point, 0)
	}
	return &MatchingResult{
		userID:                  userID,
		offerID:                 offerID,
		assignedMatchedRequests: assignedMatchedRequests,
		newPath:                 newPath,
	}
}

// GetUserID returns the user ID
func (mr *MatchingResult) GetUserID() string {
	return mr.userID
}

// SetUserID sets the user ID
func (mr *MatchingResult) SetUserID(userID string) {
	mr.userID = userID
}

// GetOfferID returns the offer ID
func (mr *MatchingResult) GetOfferID() string {
	return mr.offerID
}

// SetOfferID sets the offer ID
func (mr *MatchingResult) SetOfferID(offerID string) {
	mr.offerID = offerID
}

// GetAssignedMatchedRequests returns the assigned matched requests
func (mr *MatchingResult) GetAssignedMatchedRequests() []*MatchedRequest {
	return mr.assignedMatchedRequests
}

// SetAssignedMatchedRequests sets the assigned matched requests
func (mr *MatchingResult) SetAssignedMatchedRequests(requests []*MatchedRequest) {
	mr.assignedMatchedRequests = requests
}

// GetNewPath returns the new path
func (mr *MatchingResult) GetNewPath() []*Point {
	return mr.newPath
}

// SetNewPath sets the new path
func (mr *MatchingResult) SetNewPath(path []*Point) {
	mr.newPath = path
}
