package enums

// RoleType represents the type of entity that owns a point
type RoleType string

const (
	// Offer represents an offer owner type
	Offer RoleType = "offer"
	// Request represents a request owner type
	Request RoleType = "request"
)

// IsValid checks if the RoleType value is valid
func (r RoleType) IsValid() bool {
	switch r {
	case Offer, Request:
		return true
	default:
		return false
	}
}

// String returns the string representation of the RoleType
func (r RoleType) String() string {
	return string(r)
}
