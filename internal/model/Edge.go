package model

// Edge represents a connection between an offer node and a request node
type Edge struct {
	RequestNode *RequestNode
	NewPath     []*Point
	Pickup      *Point
	Dropoff     *Point
}

// NewEdge creates a new Edge
func NewEdge(requestNode *RequestNode, pickup, dropoff *Point, newPath []*Point) *Edge {
	if newPath == nil {
		newPath = make([]*Point, 0)
	}
	return &Edge{
		RequestNode: requestNode,
		NewPath:     newPath,
		Pickup:      pickup,
		Dropoff:     dropoff,
	}
}
