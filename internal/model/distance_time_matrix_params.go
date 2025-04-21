package model

import (
	"errors"
	"time"
)

type DistanceTimeMatrixParams struct {
	sources, targets []Coordinate
	departureTime    time.Time
	profile          Profile
}

func NewDistanceTimeMatrixParams(
	sources, targets []Coordinate,
	departureTime time.Time, profile Profile,
) (*DistanceTimeMatrixParams, error) {

	if len(sources) == 0 {
		return nil, errors.New("sources list is empty")
	}

	if len(targets) == 0 {
		return nil, errors.New("targets list is empty")
	}

	if len(sources) != len(targets) {
		return nil, errors.New("sources and targets lists must have the same length")
	}

	if departureTime.Before(time.Now()) {
		return nil, errors.New("departure time is in the past")
	}

	if profile == "" {
		return nil, errors.New("profile is empty")
	}

	return &DistanceTimeMatrixParams{
		sources:       sources,
		targets:       targets,
		departureTime: departureTime,
		profile:       profile,
	}, nil
}

func (dtmp *DistanceTimeMatrixParams) Sources() []Coordinate {
	return dtmp.sources
}

func (dtmp *DistanceTimeMatrixParams) Targets() []Coordinate {
	return dtmp.targets
}
