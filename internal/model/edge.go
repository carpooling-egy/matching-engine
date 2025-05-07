package model

// Edge represents a connection between an offer node and a request node
type Edge struct {
	requestNode *RequestNode
	newPath     []*PathPoint
	pickup      *PathPoint
	dropoff     *PathPoint
}

// NewEdge creates a new Edge
func NewEdge(requestNode *RequestNode, pickup, dropoff *PathPoint, newPath []*PathPoint) *Edge {
	if newPath == nil {
		newPath = make([]*PathPoint, 0)
	}
	return &Edge{
		requestNode: requestNode,
		newPath:     newPath,
		pickup:      pickup,
		dropoff:     dropoff,
	}
}

// RequestNode returns the request node
func (e *Edge) RequestNode() *RequestNode {
	return e.requestNode
}

// SetRequestNode sets the request node
func (e *Edge) SetRequestNode(requestNode *RequestNode) {
	e.requestNode = requestNode
}

// NewPath returns the new path
func (e *Edge) NewPath() []*PathPoint {
	return e.newPath
}

// SetNewPath sets the new path
func (e *Edge) SetNewPath(newPath []*PathPoint) {
	e.newPath = newPath
}

// Pickup returns the pickup PathPoint
func (e *Edge) Pickup() *PathPoint {
	return e.pickup
}

// SetPickup sets the pickup PathPoint
func (e *Edge) SetPickup(pickup *PathPoint) {
	e.pickup = pickup
}

// Dropoff returns the dropoff PathPoint
func (e *Edge) Dropoff() *PathPoint {
	return e.dropoff
}

// SetDropoff sets the dropoff PathPoint
func (e *Edge) SetDropoff(dropoff *PathPoint) {
	e.dropoff = dropoff
}
