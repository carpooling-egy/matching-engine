package model

import (
	"errors"
	"fmt"
	"strings"
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

func (dtm *DistanceTimeMatrixParams) String() string {
	coords := make([]string, len(dtm.points))
	for i, pt := range dtm.points {
		coords[i] = pt.String()
	}
	return fmt.Sprintf(
		"Points: [%s], Departure: %s, Profile: %s",
		strings.Join(coords, ", "),
		dtm.departureTime.Format(time.RFC3339),
		dtm.profile.String(),
	)
}
