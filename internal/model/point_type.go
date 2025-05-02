package model

// PointType represents the type of a point in a route
type PointType int

const (
	// PickupPoint represents a pickup location
	PickupPoint PointType = iota
	// DropoffPoint represents a dropoff location
	DropoffPoint
)

// String returns the string representation of the point type
func (pt PointType) String() string {
	switch pt {
	case PickupPoint:
		return "Pickup"
	case DropoffPoint:
		return "Dropoff"
	default:
		return "Unknown"
	}
}
