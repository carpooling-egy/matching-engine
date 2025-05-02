package model

// MatchingResult represents the outcome of a matching operation
type MatchingResult struct {
	offerId                 string
	assignedMatchedRequests []MatchedRequest
	path                    []Point
}

// NewMatchingResult creates a new MatchingResult with the given offer ID
func NewMatchingResult(offerId string) *MatchingResult {
	return &MatchingResult{
		offerId:                 offerId,
		assignedMatchedRequests: make([]MatchedRequest, 0),
		path:                    make([]Point, 0),
	}
}

// OfferID returns the offer ID of the matching result
func (m *MatchingResult) OfferID() string {
	return m.offerId
}

// SetOfferID sets the offer ID for the matching result
func (m *MatchingResult) SetOfferID(offerId string) {
	m.offerId = offerId
}

// AssignedMatchedRequests returns the matched requests
func (m *MatchingResult) AssignedMatchedRequests() []MatchedRequest {
	return m.assignedMatchedRequests
}

// SetAssignedMatchedRequests sets the list of assigned matched requests
func (m *MatchingResult) SetAssignedMatchedRequests(requests []MatchedRequest) {
	m.assignedMatchedRequests = requests
}

// AddAssignedMatchedRequest adds a matched request to the list
func (m *MatchingResult) AddAssignedMatchedRequest(request MatchedRequest) {
	m.assignedMatchedRequests = append(m.assignedMatchedRequests, request)
}

// Path returns the path of points for the matching result
func (m *MatchingResult) Path() []Point {
	return m.path
}

// SetPath sets the path for the matching result
func (m *MatchingResult) SetPath(path []Point) {
	m.path = path
}

// AddPoint adds a point to the path
func (m *MatchingResult) AddPoint(point Point) {
	m.path = append(m.path, point)
}
