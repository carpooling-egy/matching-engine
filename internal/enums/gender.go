package enums

// Gender represents the gender of a user
type Gender string

const (
	// Male represents male gender
	Male Gender = "male"
	// Female represents female gender
	Female Gender = "female"
)

// IsValid checks if the Gender value is valid
func (g Gender) IsValid() bool {
	switch g {
	case Male, Female:
		return true
	default:
		return false
	}
}

// String returns the string representation of the Gender
func (g Gender) String() string {
	return string(g)
}
