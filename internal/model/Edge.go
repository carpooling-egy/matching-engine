package model

// Edge represents a connection between an offer node and a request node
type Edge struct {
	requestNode *RequestNode
	newPath     []*Point
	pickup      *Point
	dropoff     *Point
}

// NewEdge creates a new Edge
func NewEdge(requestNode *RequestNode, pickup, dropoff *Point, newPath []*Point) *Edge {
	if newPath == nil {
		newPath = make([]*Point, 0)
	}
	return &Edge{
		requestNode: requestNode,
		newPath:     newPath,
		pickup:      pickup,
		dropoff:     dropoff,
	}
}
