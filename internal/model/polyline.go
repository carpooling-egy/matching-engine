package model

import (
	"errors"
	"math"
)

type Polyline struct {
	encoded   string
	precision int
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

func (p *Polyline) Decode() ([]Coordinate, error) {
	factor := math.Pow10(p.precision)
	var coordinates []Coordinate
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
