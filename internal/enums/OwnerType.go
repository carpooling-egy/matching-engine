package enums

// OwnerType represents the type of entity that owns a point
type OwnerType string

const (
	// Offer represents an offer owner type
	Offer OwnerType = "offer"
	// Request represents a request owner type
	Request OwnerType = "request"
)

// IsValid checks if the OwnerType value is valid
func (o OwnerType) IsValid() bool {
	switch o {
	case Offer, Request:
		return true
	default:
		return false
	}
}

// String returns the string representation of the OwnerType
func (o OwnerType) String() string {
	return string(o)
}
