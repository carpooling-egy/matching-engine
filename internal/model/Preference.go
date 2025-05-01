package model

import "matching-engine/internal/enums"

// Preference represents user preferences for matching
type Preference struct {
	Gender enums.Gender
	Smoker bool
	Pets   bool
}

// NewPreference creates a new Preference
func NewPreference(gender enums.Gender, smoker, pets bool) *Preference {
	return &Preference{
		Gender: gender,
		Smoker: smoker,
		Pets:   pets,
	}
}
