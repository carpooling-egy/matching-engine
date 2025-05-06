package model

// MatchCandidate represents a potential match between a request and an offer
type MatchCandidate struct {
	request *Request
	offer   *Offer
}

// NewMatchCandidate creates a new MatchCandidate
func NewMatchCandidate(request *Request, offer *Offer) *MatchCandidate {
	return &MatchCandidate{
		request: request,
		offer:   offer,
	}
}

// GetRequest returns the request
func (mc *MatchCandidate) GetRequest() *Request {
	return mc.request
}

// SetRequest sets the request
func (mc *MatchCandidate) SetRequest(request *Request) {
	mc.request = request
}

// GetOffer returns the offer
func (mc *MatchCandidate) GetOffer() *Offer {
	return mc.offer
}

// SetOffer sets the offer
func (mc *MatchCandidate) SetOffer(offer *Offer) {
	mc.offer = offer
}
