package model

import (
	"fmt"
)

const (
	ValhallaProfileAuto       ValhallaProfile = "auto"
	ValhallaProfilePedestrian ValhallaProfile = "pedestrian"
)

type ValhallaProfile string

func (p ValhallaProfile) IsValid() bool {
	switch p {
	case ValhallaProfileAuto, ValhallaProfilePedestrian:
		return true
	}
	return false
}

func (p ValhallaProfile) String() string {
	return string(p)
}

var valhallaMap = map[RoutingProfile]ValhallaProfile{
	ProfileCar:        ValhallaProfileAuto,
	ProfilePedestrian: ValhallaProfilePedestrian,
}

func ToValhallaProfile(p RoutingProfile) (ValhallaProfile, error) {
	if v, ok := valhallaMap[p]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unsupported for Valhalla: %s", p)
}
