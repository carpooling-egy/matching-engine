package models

// Preference represents user preferences for rides
type Preference struct {
	sameGender    bool 
	allowsSmoking bool 
	allowsPets    bool 
}

// No need to validate parameters as they will be read from database
// This constructor should be only used from database entities
func NewPreference(sameGender, allowsSmoking, allowsPets bool) *Preference {
	return &Preference{
		sameGender:   sameGender,
		allowsSmoking: allowsSmoking,
		allowsPets:   allowsPets,
	}
}

func (p *Preference) SameGender() bool {
	return p.sameGender
}

func (p *Preference) AllowsSmoking() bool {
	return p.allowsSmoking
}

func (p *Preference) AllowsPets() bool {
	return p.allowsPets
}
