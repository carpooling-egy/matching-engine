package model

// MatchedRequest represents a request that has been matched with an offer
type MatchedRequest struct {
	offer   *Offer
	request *Request
	pickup  Point
	dropoff Point
}

// NewMatchedRequest creates a new MatchedRequest
func NewMatchedRequest(offer *Offer, request *Request, pickup, dropoff Point) *MatchedRequest {
	return &MatchedRequest{
		offer:   offer,
		request: request,
		pickup:  pickup,
		dropoff: dropoff,
	}
}

// Offer returns the offer
func (mr *MatchedRequest) Offer() *Offer {
	return mr.offer
}

// SetOffer sets the offer
func (mr *MatchedRequest) SetOffer(offer *Offer) {
	mr.offer = offer
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
func (mr *MatchedRequest) Pickup() *Point {
	return &mr.pickup
}

// SetPickup sets the pickup point
func (mr *MatchedRequest) SetPickup(pickup Point) {
	mr.pickup = pickup
}

// Dropoff returns the dropoff point
func (mr *MatchedRequest) Dropoff() *Point {
	return &mr.dropoff
}

// SetDropoff sets the dropoff point
func (mr *MatchedRequest) SetDropoff(dropoff Point) {
	mr.dropoff = dropoff
}
