package model

// OfferNode represents a node in the offer graph
type OfferNode struct {
	offer                        *Offer
	newlyAssignedMatchedRequests []*Request
	edges                        []*Edge
	isMatched                    bool
}

// NewOfferNode creates a new OfferNode
func NewOfferNode(offer *Offer) *OfferNode {
	return &OfferNode{
		offer:                        offer,
		newlyAssignedMatchedRequests: make([]*Request, 0),
		edges:                        make([]*Edge, 0),
		isMatched:                    false,
	}
}

// Offer returns the offer
func (node *OfferNode) Offer() *Offer {
	return node.offer
}

// SetOffer sets the offer
func (node *OfferNode) SetOffer(offer *Offer) {
	node.offer = offer
}

// NewlyAssignedMatchedRequests returns the newly assigned matched requests
func (node *OfferNode) NewlyAssignedMatchedRequests() []*Request {
	return node.newlyAssignedMatchedRequests
}

// SetNewlyAssignedMatchedRequests sets the newly assigned matched requests
func (node *OfferNode) SetNewlyAssignedMatchedRequests(requests []*Request) {
	node.newlyAssignedMatchedRequests = requests
}

// Edges returns the edges
func (node *OfferNode) Edges() []*Edge {
	return node.edges
}

// SetEdges sets the edges
func (node *OfferNode) SetEdges(edges []*Edge) {
	node.edges = edges
}

// IsMatched returns whether the node is matched
func (node *OfferNode) IsMatched() bool {
	return node.isMatched
}

// SetMatched sets whether the node is matched
func (node *OfferNode) SetMatched(isMatched bool) {
	node.isMatched = isMatched
}

// GetAllRequests returns all matched requests, both existing and newly assigned
func (node *OfferNode) GetAllRequests() []*Request {
	return append(node.offer.matchedRequests, node.newlyAssignedMatchedRequests...)
}
