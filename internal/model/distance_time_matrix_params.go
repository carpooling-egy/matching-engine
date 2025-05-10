package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DistanceTimeMatrixParams struct {
	sources       []Coordinate
	targets       []Coordinate
	departureTime time.Time
	profile       Profile
}

type MatrixOption func(*DistanceTimeMatrixParams)

func WithTargets(targets []Coordinate) MatrixOption {
	return func(p *DistanceTimeMatrixParams) {
		p.targets = targets
	}
}

func WithDepartureTime(departureTime time.Time) MatrixOption {
	return func(p *DistanceTimeMatrixParams) {
		p.departureTime = departureTime
	}
}

func NewDistanceTimeMatrixParams(
	sources []Coordinate,
	profile Profile,
	opts ...MatrixOption,
) (*DistanceTimeMatrixParams, error) {

	if len(sources) == 0 {
		return nil, fmt.Errorf("sources cannot be empty")
	}

	if !profile.IsValid() {
		return nil, errors.New("invalid profile")
	}

	dtm := &DistanceTimeMatrixParams{
		sources: sources,
		targets: sources,
		profile: profile,
	}

	for _, opt := range opts {
		opt(dtm)
	}

	if len(dtm.targets) == 0 {
		return nil, errors.New("targets cannot be empty")
	}

	if dtm.departureTime.Before(time.Now()) {
		return nil, errors.New("departure time is in the past")
	}

	return dtm, nil
}

func (dtm *DistanceTimeMatrixParams) Sources() []Coordinate {
	return dtm.sources
}

func (dtm *DistanceTimeMatrixParams) Targets() []Coordinate {
	return dtm.targets
}

func (dtm *DistanceTimeMatrixParams) DepartureTime() time.Time {
	return dtm.departureTime
}

func (dtm *DistanceTimeMatrixParams) Profile() Profile {
	return dtm.profile
}

func (dtm *DistanceTimeMatrixParams) String() string {
	formatCoords := func(coords []Coordinate) string {
		parts := make([]string, len(coords))
		for i, pt := range coords {
			parts[i] = pt.String()
		}
		return "[" + strings.Join(parts, ", ") + "]"
	}

	return fmt.Sprintf(
		"Sources: %s, Targets: %s, Departure: %s, Profile: %s",
		formatCoords(dtm.sources),
		formatCoords(dtm.targets),
		dtm.departureTime,
		dtm.profile.String(),
	)
}
