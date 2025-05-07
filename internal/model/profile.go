package model

const (
	ProfileAuto       Profile = "auto"
	ProfilePedestrian Profile = "pedestrian"
)

type Profile string

func (p Profile) IsValid() bool {
	switch p {
	case ProfileAuto, ProfilePedestrian:
		return true
	}
	return false
}

func (p Profile) String() string {
	return string(p)
}
