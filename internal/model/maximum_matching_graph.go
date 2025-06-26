package model

import (
	"matching-engine/internal/collections"
)

type MaximumMatchingGraph struct {
	offerNodes   *collections.SyncMap[string, *OfferNode]
	requestNodes *collections.SyncMap[string, *RequestNode]
	edges        *collections.SyncMap[OfferRequestKey, *Edge]
}

func NewMaximumMatchingGraph() *MaximumMatchingGraph {
	return &MaximumMatchingGraph{
		offerNodes:   collections.NewSyncMap[string, *OfferNode](),
		requestNodes: collections.NewSyncMap[string, *RequestNode](),
		edges:        collections.NewSyncMap[OfferRequestKey, *Edge](),
	}
}

// OfferNodes returns the offer nodes
func (g *MaximumMatchingGraph) OfferNodes() *collections.SyncMap[string, *OfferNode] {
	return g.offerNodes
}

// SetOfferNodes sets the offer nodes
func (g *MaximumMatchingGraph) SetOfferNodes(offerNodes *collections.SyncMap[string, *OfferNode]) {
	g.offerNodes = offerNodes

}

func (g *MaximumMatchingGraph) AddOfferNode(offerNode *OfferNode) {
	g.offerNodes.Set(offerNode.Offer().ID(), offerNode)
}

// RequestNodes returns the request nodes
func (g *MaximumMatchingGraph) RequestNodes() *collections.SyncMap[string, *RequestNode] {
	return g.requestNodes
}

// SetRequestNodes sets the request nodes
func (g *MaximumMatchingGraph) SetRequestNodes(requestNodes *collections.SyncMap[string, *RequestNode]) {
	g.requestNodes = requestNodes
}

func (g *MaximumMatchingGraph) AddRequestNode(requestNode *RequestNode) {
	g.requestNodes.Set(requestNode.request.id, requestNode)
}

// Clear clears the graph
func (g *MaximumMatchingGraph) Clear() {
	g.offerNodes.ForEach(func(_ string, node *OfferNode) error {
		node.ClearEdges()
		return nil
	})
	g.offerNodes = collections.NewSyncMap[string, *OfferNode]()
	g.requestNodes = collections.NewSyncMap[string, *RequestNode]()
	g.edges.Clear()
}

// Edges returns the edges
func (g *MaximumMatchingGraph) Edges() *collections.SyncMap[OfferRequestKey, *Edge] {
	return g.edges
}

// SetEdges sets the edges
func (g *MaximumMatchingGraph) SetEdges(edges *collections.SyncMap[OfferRequestKey, *Edge]) {
	g.edges = edges
}

func (g *MaximumMatchingGraph) AddEdge(offerNode *OfferNode, requestNode *RequestNode, edge *Edge) {
	// Set the edge in the offer node
	offerNode.AddEdge(edge)

	key := NewOfferRequestKey(
		offerNode.Offer().ID(),
		requestNode.Request().ID(),
	)
	g.edges.Set(key, edge)
}

func (g *MaximumMatchingGraph) GetEdge(offer *Offer, request *Request) (*Edge, bool) {
	key := NewOfferRequestKey(
		offer.id,
		request.id,
	)
	return g.edges.Get(key)

}
