package model

import (
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
