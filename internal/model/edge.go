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

// GetRequestNode returns the request node
func (e *Edge) GetRequestNode() *RequestNode {
	return e.requestNode
}

// SetRequestNode sets the request node
func (e *Edge) SetRequestNode(requestNode *RequestNode) {
	e.requestNode = requestNode
}

// GetNewPath returns the new path
func (e *Edge) GetNewPath() []*Point {
	return e.newPath
}

// SetNewPath sets the new path
func (e *Edge) SetNewPath(newPath []*Point) {
	e.newPath = newPath
}

// GetPickup returns the pickup point
func (e *Edge) GetPickup() *Point {
	return e.pickup
}

// SetPickup sets the pickup point
func (e *Edge) SetPickup(pickup *Point) {
	e.pickup = pickup
}

// GetDropoff returns the dropoff point
func (e *Edge) GetDropoff() *Point {
	return e.dropoff
}

// SetDropoff sets the dropoff point
func (e *Edge) SetDropoff(dropoff *Point) {
	e.dropoff = dropoff
}
