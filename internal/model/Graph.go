package model

type Graph struct {
	offerNodes []*OfferNode
}

func NewGraph(offerNodes []*OfferNode) *Graph {
	return &Graph{
		offerNodes: offerNodes,
	}
}

func (g *Graph) GetOfferNodes() []*OfferNode {
	return g.offerNodes
}

func (g *Graph) AddOfferNode(offerNode *OfferNode) {
	g.offerNodes = append(g.offerNodes, offerNode)
}

func (g *Graph) RemoveOfferNode(offerNode *OfferNode) {
	for i, node := range g.offerNodes {
		if node == offerNode {
			g.offerNodes = append(g.offerNodes[:i], g.offerNodes[i+1:]...)
			break
		}
	}
}
