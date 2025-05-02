package model

import "matching-engine/internal/enums"

// Preference represents user preferences for matching
type Preference struct {
	gender enums.Gender
	smoker bool
	pets   bool
}

// NewPreference creates a new Preference
func NewPreference(gender enums.Gender, smoker, pets bool) *Preference {
	return &Preference{
		gender: gender,
		smoker: smoker,
		pets:   pets,
	}
}
