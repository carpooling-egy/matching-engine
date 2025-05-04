package model

import (
	"errors"
)

type DistanceUnit int

const (
	Kilometer DistanceUnit = iota
	Mile
)

func (d DistanceUnit) IsValid() bool {
	switch d {
	case Kilometer, Mile:
		return true
	}
	return false
}

func (d DistanceUnit) String() string {
	switch d {
	case Kilometer:
		return "Kilometer"
	case Mile:
		return "Mile"
	default:
		return ""
	}
}

type Distance struct {
	value float32
	unit  DistanceUnit
}

func NewDistance(value float32, unit DistanceUnit) (*Distance, error) {
	if value < 0 {
		return nil, errors.New("distance cannot be negative")
	}

	if !unit.IsValid() {
		return nil, errors.New("invalid distance unit")
	}

	return &Distance{value: value, unit: unit}, nil
}

func (d *Distance) Value() float32 {
	return d.value
}

func (d *Distance) Unit() DistanceUnit {
	return d.unit
}
