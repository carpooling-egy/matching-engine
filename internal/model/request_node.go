package model

// RequestNode represents a node in the request graph
type RequestNode struct {
	request *Request
}

// NewRequestNode creates a new RequestNode
func NewRequestNode(request *Request) *RequestNode {
	return &RequestNode{
		request: request,
	}
}

// GetRequest returns the request
func (rn *RequestNode) GetRequest() *Request {
	return rn.request
}

// SetRequest sets the request
func (rn *RequestNode) SetRequest(request *Request) {
	rn.request = request
}
