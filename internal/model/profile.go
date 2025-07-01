package model

type RoutingProfile string

const (
	ProfileCar        RoutingProfile = "car"
	ProfilePedestrian RoutingProfile = "pedestrian"
)

func (p RoutingProfile) IsValid() bool {
	switch p {
	case ProfileCar, ProfilePedestrian:
		return true
	}
	return false
}

func (p RoutingProfile) String() string {
	return string(p)
}
