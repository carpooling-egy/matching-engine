package enums

// PointType represents the type of a point (pickup or dropoff)
type PointType string

const (
	// Pickup represents a pickup point
	Pickup PointType = "pickup"
	// Dropoff represents a dropoff point
	Dropoff PointType = "dropoff"
)

// IsValid checks if the PointType value is valid
func (p PointType) IsValid() bool {
	switch p {
	case Pickup, Dropoff:
		return true
	default:
		return false
	}
}

// String returns the string representation of the PointType
func (p PointType) String() string {
	return string(p)
}
