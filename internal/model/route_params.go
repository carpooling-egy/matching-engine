package model

import (
	"time"
)

type RouteParams struct {
	waypoints []Coordinate // [source, pickup/drop-off points of riders, destination]
	// TODO check timezone stuff
	departureTime time.Time
}

func NewRouteParams(waypoints []Coordinate, departureTime time.Time) (*RouteParams, error) {
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
