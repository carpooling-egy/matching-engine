package model

// MatchedRequest represents a request that has been matched with an offer
type MatchedRequest struct {
	request *Request
	pickup  PathPoint
	dropoff PathPoint
}

// NewMatchedRequest creates a new MatchedRequest
func NewMatchedRequest(request *Request, pickup, dropoff PathPoint) *MatchedRequest {
	return &MatchedRequest{
		request: request,
		pickup:  pickup,
		dropoff: dropoff,
	}
}

// Request returns the request
func (mr *MatchedRequest) Request() *Request {
	return mr.request
}

// SetRequest sets the request
func (mr *MatchedRequest) SetRequest(request *Request) {
	mr.request = request
}

// Pickup returns the pickup point
func (mr *MatchedRequest) Pickup() *PathPoint {
	return &mr.pickup
}

// SetPickup sets the pickup point
func (mr *MatchedRequest) SetPickup(pickup PathPoint) {
	mr.pickup = pickup
}

// Dropoff returns the dropoff point
func (mr *MatchedRequest) Dropoff() *PathPoint {
	return &mr.dropoff
}

// SetDropoff sets the dropoff point
func (mr *MatchedRequest) SetDropoff(dropoff PathPoint) {
	mr.dropoff = dropoff
}
