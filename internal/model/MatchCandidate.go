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
