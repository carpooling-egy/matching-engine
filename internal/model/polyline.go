package model

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
)

type Polyline struct {
	encoded     string
	precision   int
	coordinates LineString

	once      sync.Once
	decodeErr error
}

type Option func(*Polyline) error

func WithPrecision(precision int) Option {
	return func(p *Polyline) error {
		if precision < 0 {
			return errors.New("precision cannot be negative")
		}
		p.precision = precision
		return nil
	}
}

func NewPolyline(encoded string, opts ...Option) (*Polyline, error) {
	if encoded == "" {
		return nil, errors.New("encoded string is empty")
	}

	p := &Polyline{
		encoded:   encoded,
		precision: 6,
	}

	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *Polyline) Encoded() string {
	return p.encoded
}

func (p *Polyline) Precision() int {
	return p.precision
}

func (p *Polyline) Coordinates() (LineString, error) {
	p.once.Do(func() {
		p.coordinates, p.decodeErr = p.decode()
	})
	return p.coordinates, p.decodeErr
}

func (p *Polyline) decode() (LineString, error) {
	factor := math.Pow10(p.precision)
	var coordinates LineString
	index := 0
	lat, lng := 0, 0

	for index < len(p.encoded) {
		deltaLat, err := p.decodeDelta(&index)
		if err != nil {
			return nil, err
		}
		lat += deltaLat

		deltaLng, err := p.decodeDelta(&index)
		if err != nil {
			return nil, err
		}
		lng += deltaLng

		coord, err := NewCoordinate(
			float64(lat)/factor,
			float64(lng)/factor,
		)
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, *coord)
	}

	return coordinates, nil
}

func (p *Polyline) decodeDelta(index *int) (int, error) {
	var result, shift int
	b := 0x20

	for *index < len(p.encoded) && b >= 0x20 {
		if *index >= len(p.encoded) {
			return 0, errors.New("malformed polyline encoding")
		}
		b = int(p.encoded[*index]) - 63
		result |= (b & 0x1F) << shift
		shift += 5
		*index++
	}

	if (result & 1) != 0 {
		return ^(result >> 1), nil
	}
	return result >> 1, nil
}

func (p *Polyline) String() string {
	// Truncate encoded string if too long
	encoded := p.encoded
	if len(encoded) > 20 {
		encoded = encoded[:17] + "..."
	}

	// Always decode coordinates
	coords, err := p.Coordinates()
	if err != nil {
		return fmt.Sprintf("Polyline{encoded: %q, precision: %d, coords: error(%s)}", encoded, p.precision, err)
	}

	// Format coordinates (limit to 3 for brevity if many)
	coordsInfo := fmt.Sprintf("%d points", len(coords))
	if len(coords) > 0 {
		coordsInfo += ": ["
		maxDisplay := 3
		for i, coord := range coords {
			if i >= maxDisplay {
				coordsInfo += fmt.Sprintf(", ...+%d more", len(coords)-maxDisplay)
				break
			}
			if i > 0 {
				coordsInfo += ", "
			}
			coordsInfo += coord.String()
		}
		coordsInfo += "]"
	} else {
		coordsInfo += ": []"
	}

	return fmt.Sprintf("Polyline{encoded: %q, precision: %d, coords: %s}", encoded, p.precision, coordsInfo)
}

func (p *Polyline) ToWKT() (string, error) {
	coords, err := p.Coordinates()
	if err != nil {
		return "", err
	}
	if len(coords) < 2 {
		return "", errors.New("LINESTRING requires at least 2 points")
	}

	var parts []string
	for _, coord := range coords {
		parts = append(parts, fmt.Sprintf("%f %f", coord.Lng(), coord.Lat()))
	}
	return "LINESTRING(" + strings.Join(parts, ",") + ")", nil
}
