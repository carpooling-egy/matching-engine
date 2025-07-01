package model

import (
	"fmt"
)

const (
	OSRMProfileCar  OSRMProfile = "car"
	OSRMProfileFoot OSRMProfile = "foot"
)

type OSRMProfile string

func (p OSRMProfile) IsValid() bool {
	switch p {
	case OSRMProfileCar, OSRMProfileFoot:
		return true
	}
	return false
}

func (p OSRMProfile) String() string {
	return string(p)
}

var osrmMap = map[RoutingProfile]OSRMProfile{
	ProfileCar:        OSRMProfileCar,
	ProfilePedestrian: OSRMProfileFoot,
}

func ToOSRMProfile(profile RoutingProfile) (OSRMProfile, error) {
	if v, ok := osrmMap[profile]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unsupported profile for OSRM: %s", profile)
}
