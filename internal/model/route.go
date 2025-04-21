package model

import (
	"errors"
	"time"
)

type Route struct {
	polyline Polyline
	distance Distance
	time     time.Duration
}

func NewRoute(polyline Polyline, distance Distance, time time.Duration) (*Route, error) {
	if polyline == (Polyline{}) {
		return nil, errors.New("polyline is nil")
	}

	if distance == (Distance{}) {
		return nil, errors.New("distance cannot be negative")
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

func (r *Route) Polyline() (Polyline, error) {
	if r == nil {
		return Polyline{}, errors.New("nil route reference")
	}

	return r.polyline, nil
}

func (r *Route) Distance() (Distance, error) {
	if r == nil {
		return Distance{}, errors.New("nil route reference")
	}

	return r.distance, nil
}

func (r *Route) Time() (time.Duration, error) {
	if r == nil {
		return 0, errors.New("nil route reference")
	}

	return r.time, nil
}
