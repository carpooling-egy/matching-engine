package model

// OfferNode represents a node in the offer graph
type OfferNode struct {
	offer                        *Offer
	newlyAssignedMatchedRequests []MatchedRequest
	edges                        []Edge
	isMatched                    bool
}

// NewOfferNode creates a new OfferNode
func NewOfferNode(offer *Offer) *OfferNode {
	return &OfferNode{
		offer:                        offer,
		newlyAssignedMatchedRequests: make([]MatchedRequest, 0),
		edges:                        make([]Edge, 0),
		isMatched:                    false,
	}
}

// GetAllRequests returns all matched requests, both existing and newly assigned
func (node *OfferNode) GetAllRequests() []MatchedRequest {
	return append(node.offer.matchedRequests, node.newlyAssignedMatchedRequests...)
}
