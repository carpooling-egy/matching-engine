package model

import "errors"

type DistanceUnit int

const (
	Meter DistanceUnit = iota
	Kilometer
	Mile
)

type Distance struct {
	value float64
	unit  DistanceUnit
}

func NewDistance(value float64, unit DistanceUnit) (*Distance, error) {
	if value < 0 {
		return nil, errors.New("distance cannot be negative")
	}

	if unit != Meter && unit != Kilometer && unit != Mile {
		return nil, errors.New("invalid distance unit")
	}

	return &Distance{value: value, unit: unit}, nil
}

func (d *Distance) Value() (float64, error) {
	if d == nil {
		return 0, errors.New("nil distance reference")
	}

	return d.value, nil
}

func (d *Distance) Unit() (DistanceUnit, error) {
	if d == nil {
		return 0, errors.New("nil distance reference")
	}

	return d.unit, nil
}
