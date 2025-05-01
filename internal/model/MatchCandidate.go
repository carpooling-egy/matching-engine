package model

// MatchCandidate represents a potential match between a request and an offer
type MatchCandidate struct {
	Request *Request
	Offer   *Offer
}

// NewMatchCandidate creates a new MatchCandidate
func NewMatchCandidate(request *Request, offer *Offer) *MatchCandidate {
	return &MatchCandidate{
		Request: request,
		Offer:   offer,
	}
}
