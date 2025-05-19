package model

import "matching-engine/internal/collections"

type Graph struct {
	offerNodes *collections.Set[*OfferNode]
}

func NewGraph() *Graph {
	return &Graph{
		offerNodes: collections.NewSet[*OfferNode](),
	}
}

// OfferNodes returns the offer nodes
func (g *Graph) OfferNodes() *collections.Set[*OfferNode] {
	return g.offerNodes
}

// SetOfferNodes sets the offer nodes
func (g *Graph) SetOfferNodes(offerNodes *collections.Set[*OfferNode]) {
	g.offerNodes = offerNodes

}

func (g *Graph) AddOfferNode(offerNode *OfferNode) {
	g.offerNodes.Add(offerNode)
}

func (g *Graph) RemoveOfferNode(offerNode *OfferNode) {
	g.offerNodes.Remove(offerNode)
}

func (g *Graph) Clear() {
	g.OfferNodes().ForEach(func(node *OfferNode) {
		node.ClearEdges()
	})
	g.offerNodes.Clear()
}
