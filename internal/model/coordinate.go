package model

import "errors"

type Coordinate struct {
	lat, lng float64
}

func NewCoordinate(lat, lng float64) (*Coordinate, error) {
	if err := validateLatitude(lat); err != nil {
		return nil, err
	}

	if err := validateLongitude(lng); err != nil {
		return nil, err
	}

	return &Coordinate{lat: lat, lng: lng}, nil
}

func (c *Coordinate) Lat() float64 {
	return c.lat
}

func (c *Coordinate) Lng() float64 {
	return c.lng
}

func validateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("latitude must be between -90 and 90 degrees")
	}
	return nil
}

func validateLongitude(lng float64) error {
	if lng < -180 || lng > 180 {
		return errors.New("longitude must be between -180 and 180 degrees")
	}
	return nil
}
