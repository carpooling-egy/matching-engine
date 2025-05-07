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

// GetOffer returns the offer
func (mr *MatchedRequest) GetOffer() *Offer {
	return mr.offer
}

// SetOffer sets the offer
func (mr *MatchedRequest) SetOffer(offer *Offer) {
	mr.offer = offer
}

// GetRequest returns the request
func (mr *MatchedRequest) GetRequest() *Request {
	return mr.request
}

// SetRequest sets the request
func (mr *MatchedRequest) SetRequest(request *Request) {
	mr.request = request
}

// GetPickup returns the pickup point
func (mr *MatchedRequest) GetPickup() *Point {
	return &mr.pickup
}

// SetPickup sets the pickup point
func (mr *MatchedRequest) SetPickup(pickup Point) {
	mr.pickup = pickup
}

// GetDropoff returns the dropoff point
func (mr *MatchedRequest) GetDropoff() *Point {
	return &mr.dropoff
}

// SetDropoff sets the dropoff point
func (mr *MatchedRequest) SetDropoff(dropoff Point) {
	mr.dropoff = dropoff
}
