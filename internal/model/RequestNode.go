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
