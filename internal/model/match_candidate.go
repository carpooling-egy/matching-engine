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

// Request returns the request
func (mc *MatchCandidate) Request() *Request {
	return mc.request
}

// SetRequest sets the request
func (mc *MatchCandidate) SetRequest(request *Request) {
	mc.request = request
}

// Offer returns the offer
func (mc *MatchCandidate) Offer() *Offer {
	return mc.offer
}

// SetOffer sets the offer
func (mc *MatchCandidate) SetOffer(offer *Offer) {
	mc.offer = offer
}
