package model

import (
	"errors"
	"fmt"
	"strings"
)

type LineString []Coordinate

func (ls LineString) IsClosed() bool {
	// A closed ring must have at least three distinct points,
	// with the first and last point equal.
	if len(ls) < 4 {
		return false
	}

	first, last := ls[0], ls[len(ls)-1]
	return first.Lat() == last.Lat() && first.Lng() == last.Lng()
}

func (ls LineString) String() string {
	coords := make([]string, len(ls))
	for i, pt := range ls {
		coords[i] = pt.String()
	}
	return fmt.Sprintf("LineString[%s]", strings.Join(coords, ", "))
}

type Contour struct {
	value float32
	// Unit describes what Value measures (e.g. "minutes" or "kilometers").
	unit string
}

func NewContour(value float32, unit string) (*Contour, error) {
	if value < 0 {
		return nil, errors.New("contour value must be non-negative")
	}
	if unit == "" {
		return nil, errors.New("contour unit must be non-empty")
	}
	return &Contour{
		value: value,
		unit:  unit,
	}, nil
}

func (c *Contour) Value() float32 {
	return c.value
}

func (c *Contour) Unit() string {
	return c.unit
}

func (c *Contour) String() string {
	return fmt.Sprintf("%.2f %s", c.value, c.unit)
}

type Isochrone struct {
	contour *Contour
	ring    *LineString
}

func NewIsochrone(contour *Contour, ring *LineString) (*Isochrone, error) {
	if contour == nil {
		return nil, errors.New("contour must be non-nil")
	}

	if ring == nil {
		return nil, errors.New("ring must be non-nil")
	}

	if !ring.IsClosed() {
		return nil, errors.New("ring must be a closed LineString with at least 4 points")
	}

	return &Isochrone{contour: contour, ring: ring}, nil
}

func (i *Isochrone) Contour() *Contour {
	return i.contour
}

func (i *Isochrone) Geometry() *LineString {
	return i.ring
}

func (i *Isochrone) String() string {
	return fmt.Sprintf("Isochrone{contour: %s, ring: %s}",
		i.contour.String(), i.ring.String())
}
