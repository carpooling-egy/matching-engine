package model

import (
	"errors"
	"fmt"
)

type DistanceUnit string

const (
	DistanceUnitKilometer DistanceUnit = "kilometer"
	DistanceUnitMile      DistanceUnit = "mile"
)

func (du DistanceUnit) IsValid() bool {
	switch du {
	case DistanceUnitKilometer, DistanceUnitMile:
		return true
	}
	return false
}

func (du DistanceUnit) String() string {
	return string(du)
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

func (d *Distance) String() string {
	return fmt.Sprintf("%.2f %s", d.value, d.unit)
}
