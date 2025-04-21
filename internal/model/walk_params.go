package model

import "errors"

type WalkParams struct {
	origin, destination Coordinate
}

func NewWalkParams(origin, destination Coordinate) (*WalkParams, error) {
	return &WalkParams{
		origin:      origin,
		destination: destination,
	}, nil
}

func (p *WalkParams) Origin() (Coordinate, error) {
	if p == nil {
		return Coordinate{}, errors.New("nil walk params reference")
	}

	return p.origin, nil
}

func (p *WalkParams) Destination() (Coordinate, error) {
	if p == nil {
		return Coordinate{}, errors.New("nil walk params reference")
	}

	return p.destination, nil
}
