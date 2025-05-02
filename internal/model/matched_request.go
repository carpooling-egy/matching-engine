package model

// MatchedRequest represents a request that has been matched with pickup and dropoff points
type MatchedRequest struct {
	requestId    string
	pickupPoint  []Point
	dropoffPoint []Point
}

// NewMatchedRequest creates a new MatchedRequest with the given request ID, pickup, and dropoff points
func NewMatchedRequest(requestId string, pickupPoints, dropoffPoints []Point) *MatchedRequest {
	return &MatchedRequest{
		requestId:    requestId,
		pickupPoint:  pickupPoints,
		dropoffPoint: dropoffPoints,
	}
}

// NewMatchedRequestWithRequestId NewMatchedRequest creates a new MatchedRequest with the given request ID
func NewMatchedRequestWithRequestId(requestId string) *MatchedRequest {
	return &MatchedRequest{
		requestId:    requestId,
		pickupPoint:  []Point{},
		dropoffPoint: []Point{},
	}
}

// RequestID returns the ID of this matched request
func (m *MatchedRequest) RequestID() string {
	return m.requestId
}

// SetRequestID sets the ID for this matched request
func (m *MatchedRequest) SetRequestID(requestId string) {
	m.requestId = requestId
}

// PickupPoints returns the pickup points for this matched request
func (m *MatchedRequest) PickupPoints() []Point {
	return m.pickupPoint
}

// SetPickupPoints sets the pickup points for this matched request
func (m *MatchedRequest) SetPickupPoints(points []Point) {
	m.pickupPoint = points
}

// AddPickupPoint adds a pickup point to this matched request
func (m *MatchedRequest) AddPickupPoint(point Point) {
	m.pickupPoint = append(m.pickupPoint, point)
}

// DropoffPoints returns the dropoff points for this matched request
func (m *MatchedRequest) DropoffPoints() []Point {
	return m.dropoffPoint
}

// SetDropoffPoints sets the dropoff points for this matched request
func (m *MatchedRequest) SetDropoffPoints(points []Point) {
	m.dropoffPoint = points
}

// AddDropoffPoint adds a dropoff point to this matched request
func (m *MatchedRequest) AddDropoffPoint(point Point) {
	m.dropoffPoint = append(m.dropoffPoint, point)
}
