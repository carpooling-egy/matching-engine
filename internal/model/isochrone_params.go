package model

import "errors"

type IsochroneParams struct {
	origin   *Coordinate
	distance *Distance
	profile  Profile
}

func NewIsochroneParams(
	origin *Coordinate,
	distance *Distance,
	profile Profile,
) (*IsochroneParams, error) {
	if origin == nil {
		return nil, errors.New("origin is nil")
	}

	if distance == nil {
		return nil, errors.New("distance is nil")
	}

	if profile == "" {
		return nil, errors.New("profile is empty")
	}

	return &IsochroneParams{
		origin:   origin,
		distance: distance,
		profile:  profile,
	}, nil
}

func (ip *IsochroneParams) Origin() *Coordinate {
	return ip.origin
}

func (ip *IsochroneParams) Distance() *Distance {
	return ip.distance
}

func (ip *IsochroneParams) Profile() Profile {
	return ip.profile
}
