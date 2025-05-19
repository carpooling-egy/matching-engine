package model

import (
	"matching-engine/internal/enums"
)

// Preference represents user preferences for rides
type Preference struct {
	gender        enums.Gender
	sameGender    bool
}

// NewPreference Creates a new Preference. No need to validate parameters as they will be read from database
// This constructor should be only used from database entities
func NewPreference(gender enums.Gender, sameGender bool) *Preference {
	return &Preference{
		gender:        gender,
		sameGender:    sameGender,
	}
}

func (p *Preference) Gender() enums.Gender {
	return p.gender
}

func (p *Preference) SameGender() bool {
	return p.sameGender
}
