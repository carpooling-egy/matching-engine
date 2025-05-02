package model

// MatchedRequest represents a request that has been matched with pickup and dropoff points
type MatchedRequest struct {
	requestId    string
	pickupPoint  Point
	dropoffPoint Point
}

// NewMatchedRequest creates a new MatchedRequest with the given request ID, pickup, and dropoff points
func NewMatchedRequest(requestId string, pickupPoint, dropoffPoint Point) *MatchedRequest {
	return &MatchedRequest{
		requestId:    requestId,
		pickupPoint:  pickupPoint,
		dropoffPoint: dropoffPoint,
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

// PickupPoint PickupPoints returns the pickup points for this matched request
func (m *MatchedRequest) PickupPoint() Point {
	return m.pickupPoint
}

// SetPickupPoint SetPickupPoints sets the pickup points for this matched request
func (m *MatchedRequest) SetPickupPoint(point Point) {
	m.pickupPoint = point
}

// DropoffPoint returns the dropoff points for this matched request
func (m *MatchedRequest) DropoffPoint() Point {
	return m.dropoffPoint
}

// SetDropoffPoint sets the dropoff points for this matched request
func (m *MatchedRequest) SetDropoffPoint(points Point) {
	m.dropoffPoint = points
}
