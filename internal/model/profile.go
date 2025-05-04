package model

const (
	Auto       Profile = "auto"
	Pedestrian Profile = "pedestrian"
)

type Profile string

func (p Profile) IsValid() bool {
	switch p {
	case Auto, Pedestrian:
		return true
	}
	return false
}

func (p Profile) String() string {
	switch p {
	case Auto:
		return "auto"
	case Pedestrian:
		return "pedestrian"
	default:
		return ""
	}
}
