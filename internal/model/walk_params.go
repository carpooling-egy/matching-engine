package model

import (
	"errors"
	"fmt"
)

type WalkParams struct {
	origin, destination *Coordinate
}

func NewWalkParams(origin, destination *Coordinate) (*WalkParams, error) {
	if origin == nil {
		return nil, errors.New("origin coordinate is nil")
	}

	if destination == nil {
		return nil, errors.New("destination coordinate is nil")
	}

	return &WalkParams{
		origin:      origin,
		destination: destination,
	}, nil
}

func (p *WalkParams) Origin() *Coordinate {
	return p.origin
}

func (p *WalkParams) Destination() *Coordinate {
	return p.destination
}

func (p *WalkParams) String() string {
	return fmt.Sprintf("WalkParams{origin=%s, destination=%s}", p.origin, p.destination)
}
