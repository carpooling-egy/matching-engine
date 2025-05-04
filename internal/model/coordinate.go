package models

type Coordinate struct {
	lat, lng float64
}

func NewCoordinate(lat, lng float64) (*Coordinate, error) {
	return &Coordinate{lat: lat, lng: lng}, nil
}

func (c *Coordinate) Lat() float64 {
	return c.lat
}

func (c *Coordinate) Lng() float64 {
	return c.lng
}
