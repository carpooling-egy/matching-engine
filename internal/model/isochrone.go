package model

import (
	"errors"
	"fmt"
)

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

func (i *Isochrone) Polygons() [][][]Coordinate {
	if i.ring == nil || len(*i.ring) < 4 {
		return nil
	}

	// For now, we only have one polygon with one outer ring (no holes)
	polygon := make([][][]Coordinate, 1)

	// Clone the coordinates to avoid any potential issues with the original data
	ring := make([]Coordinate, len(*i.ring))
	copy(ring, *i.ring)

	polygon[0] = [][]Coordinate{ring}

	return polygon
}

func (i *Isochrone) String() string {
	return fmt.Sprintf("Isochrone{contour: %s, ring: %s}",
		i.contour.String(), i.ring.String())
}
