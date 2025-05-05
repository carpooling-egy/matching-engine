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
	return fmt.Sprintf("LineString{len=%d}[%s]", len(ls), strings.Join(coords, ", "))
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
