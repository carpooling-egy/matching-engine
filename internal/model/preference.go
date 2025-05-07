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

// Gender returns the gender preference
func (p *Preference) Gender() enums.Gender {
	return p.gender
}

// SetGender sets the gender preference
func (p *Preference) SetGender(gender enums.Gender) {
	p.gender = gender
}

// IsSmoker returns whether smoking is allowed
func (p *Preference) IsSmoker() bool {
	return p.smoker
}

// SetSmoker sets whether smoking is allowed
func (p *Preference) SetSmoker(smoker bool) {
	p.smoker = smoker
}

// HasPets returns whether pets are allowed
func (p *Preference) HasPets() bool {
	return p.pets
}

// SetPets sets whether pets are allowed
func (p *Preference) SetPets(pets bool) {
	p.pets = pets
}
