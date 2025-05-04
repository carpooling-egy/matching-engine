package model

import (
	"errors"
	"time"
)

type RouteParams struct {
	waypoints []Coordinate // [source, pickup/drop-off points of riders, destination]
	// TODO check timezone stuff
	departureTime time.Time
}

func NewRouteParams(
	waypoints []Coordinate,
	departureTime time.Time,
) (*RouteParams, error) {

	if len(waypoints) < 2 {
		return nil, errors.New("at least two waypoints are required")
	}

	if departureTime.Before(time.Now()) {
		return nil, errors.New("departure time cannot be in the past")
	}

	return &RouteParams{
		waypoints:     waypoints,
		departureTime: departureTime,
	}, nil
}

func (p *RouteParams) Waypoints() []Coordinate {
	return p.waypoints
}

func (p *RouteParams) DepartureTime() time.Time {
	return p.departureTime
}
