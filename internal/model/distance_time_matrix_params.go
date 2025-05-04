package model

import (
	"errors"
	"time"
)

type DistanceTimeMatrixParams struct {
	points        []Coordinate
	departureTime time.Time
	profile       Profile
}

func NewDistanceTimeMatrixParams(
	points []Coordinate,
	departureTime time.Time,
	profile Profile,
) (*DistanceTimeMatrixParams, error) {

	if len(points) == 0 {
		return nil, errors.New("points list is empty")
	}

	if departureTime.Before(time.Now()) {
		return nil, errors.New("departure time is in the past")
	}

	if !profile.IsValid() {
		return nil, errors.New("invalid profile")
	}

	return &DistanceTimeMatrixParams{
		points:        points,
		departureTime: departureTime,
		profile:       profile,
	}, nil
}

func (dtm *DistanceTimeMatrixParams) Points() []Coordinate {
	return dtm.points
}

func (dtm *DistanceTimeMatrixParams) DepartureTime() time.Time {
	return dtm.departureTime
}

func (dtm *DistanceTimeMatrixParams) Profile() Profile {
	return dtm.profile
}
