package model

// Edge represents a connection between an offer node and a request node
type Edge struct {
	requestNode *RequestNode
	newPath     []PathPoint
}

// NewEdge creates a new Edge
func NewEdge(requestNode *RequestNode, newPath []PathPoint) *Edge {
	if newPath == nil {
		newPath = make([]PathPoint, 0)
	}
	return &Edge{
		requestNode: requestNode,
		newPath:     newPath,
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
func (e *Edge) NewPath() []PathPoint {
	return e.newPath
}

// SetNewPath sets the new path
func (e *Edge) SetNewPath(newPath []PathPoint) {
	e.newPath = newPath
}
