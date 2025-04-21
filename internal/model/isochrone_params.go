package model

import "errors"

type IsochroneParams struct {
	origin   Coordinate
	distance Distance
	profile  Profile
}

func NewIsochroneParams(
	origin Coordinate,
	distance Distance,
	profile Profile,
) (*IsochroneParams, error) {
	if origin == (Coordinate{}) {
		return nil, errors.New("origin coordinate is empty")
	}

	if distance.value < 0 {
		return nil, errors.New("distance is negative")
	}

	if distance.unit != Meter && distance.unit != Kilometer && distance.unit != Mile {
		return nil, errors.New("invalid distance unit")
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

func (ip *IsochroneParams) Origin() (Coordinate, error) {
	if ip == nil {
		return Coordinate{}, errors.New("nil isochrone params reference")
	}
	return ip.origin, nil
}

func (ip *IsochroneParams) Distance() (Distance, error) {
	if ip == nil {
		return Distance{}, errors.New("nil isochrone params reference")
	}
	return ip.distance, nil
}

func (ip *IsochroneParams) Profile() (Profile, error) {
	if ip == nil {
		return "", errors.New("nil isochrone params reference")
	}
	return ip.profile, nil
}
