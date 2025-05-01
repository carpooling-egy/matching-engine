package model

// OfferNode represents a node in the offer graph
type OfferNode struct {
	Offer                        *Offer
	NewlyAssignedMatchedRequests []MatchedRequest
	Edges                        []Edge
	Matched                      bool
}

// NewOfferNode creates a new OfferNode
func NewOfferNode(offer *Offer) *OfferNode {
	return &OfferNode{
		Offer:                        offer,
		NewlyAssignedMatchedRequests: make([]MatchedRequest, 0),
		Edges:                        make([]Edge, 0),
		Matched:                      false,
	}
}

// GetAllRequests returns all matched requests, both existing and newly assigned
func (node *OfferNode) GetAllRequests() []MatchedRequest {
	return append(node.Offer.MatchedRequests, node.NewlyAssignedMatchedRequests...)
}
