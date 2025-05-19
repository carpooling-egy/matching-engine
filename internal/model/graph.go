package model

import (
	"matching-engine/internal/collections"
)

type Graph struct {
	offerNodes   *collections.SyncMap[string, *OfferNode]
	requestNodes *collections.SyncMap[string, *RequestNode]
	edges        *collections.SyncMap[OfferRequestKey, *Edge]
}

func NewGraph() *Graph {
	return &Graph{
		offerNodes:   collections.NewSyncMap[string, *OfferNode](),
		requestNodes: collections.NewSyncMap[string, *RequestNode](),
		edges:        collections.NewSyncMap[OfferRequestKey, *Edge](),
	}
}

// OfferNodes returns the offer nodes
func (g *Graph) OfferNodes() *collections.SyncMap[string, *OfferNode] {
	return g.offerNodes
}

// SetOfferNodes sets the offer nodes
func (g *Graph) SetOfferNodes(offerNodes *collections.SyncMap[string, *OfferNode]) {
	g.offerNodes = offerNodes

}

func (g *Graph) AddOfferNode(offerNode *OfferNode) {
	g.offerNodes.Set(offerNode.offer.id, offerNode)
}

func (g *Graph) RemoveOfferNode(offerNode *OfferNode) {
	g.offerNodes.Delete(offerNode.offer.id)
}

func (g *Graph) Clear() {
	g.OfferNodes().ForEach(func(node *OfferNode) {
		node.ClearEdges()
	})
	g.offerNodes.Clear()
}

// RequestNodes returns the request nodes
func (g *Graph) RequestNodes() *collections.SyncMap[string, *RequestNode] {
	return g.requestNodes
}

// SetRequestNodes sets the request nodes
func (g *Graph) SetRequestNodes(requestNodes *collections.SyncMap[string, *RequestNode]) {
	g.requestNodes = requestNodes
}

func (g *Graph) AddRequestNode(requestNode *RequestNode) {
	g.requestNodes.Set(requestNode.request.id, requestNode)
}

func (g *Graph) RemoveRequestNode(requestNode *RequestNode) {
	g.requestNodes.Delete(requestNode.request.id)
}

func (g *Graph) ClearRequestNodes() {
	g.requestNodes.Clear()
}

// Clear clears the graph
func (g *Graph) Clear() {
	g.offerNodes.Clear()
	g.requestNodes.Clear()
	g.ClearEdges()
}

// Edges returns the edges
func (g *Graph) Edges() *collections.SyncMap[OfferRequestKey, *Edge] {
	return g.edges
}

// SetEdges sets the edges
func (g *Graph) SetEdges(edges *collections.SyncMap[OfferRequestKey, *Edge]) {
	g.edges = edges
}

func (g *Graph) AddEdge(offer *Offer, request *Request, edge *Edge) {
	key := NewOfferRequestKey(
		offer.id,
		request.id,
	)
	g.edges.Set(key, edge)
}

func (g *Graph) RemoveEdge(offer *Offer, request *Request) {
	key := NewOfferRequestKey(
		offer.id,
		request.id,
	)
	g.edges.Delete(key)
}

func (g *Graph) GetEdge(offer *Offer, request *Request) (*Edge, bool) {
	key := NewOfferRequestKey(
		offer.id,
		request.id,
	)
	return g.edges.Get(key)

}

func (g *Graph) ClearEdges() {
	g.edges.Clear()
}
