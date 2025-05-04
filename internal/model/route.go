package model

import (
	"errors"
	"time"
)

type Route struct {
	polyline *Polyline
	distance *Distance
	time     time.Duration
}

func NewRoute(
	polyline *Polyline,
	distance *Distance,
	time time.Duration,
) (*Route, error) {

	if polyline == nil {
		return nil, errors.New("polyline is nil")
	}

	if distance == nil {
		return nil, errors.New("distance is nil")
	}

	if time < 0 {
		return nil, errors.New("time cannot be negative")
	}

	return &Route{
		polyline: polyline,
		distance: distance,
		time:     time,
	}, nil
}

func (r *Route) Polyline() *Polyline {
	return r.polyline
}

func (r *Route) Distance() *Distance {
	return r.distance
}

func (r *Route) Time() time.Duration {
	return r.time
}
