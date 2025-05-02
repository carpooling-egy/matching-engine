package model

import "time"

// Point represents a location with temporal information, and associated request ID
type Point struct {
	requestId string
	point     *Coordinate
	time      time.Time
	pointType PointType
}

// NewPoint creates a new Point with the given parameters
func NewPoint(requestId string, coordinate *Coordinate, time time.Time, pointType PointType) Point {
	return Point{
		requestId: requestId,
		point:     coordinate,
		time:      time,
		pointType: pointType,
	}
}

// RequestID returns the request ID associated with this point
func (p *Point) RequestID() string {
	return p.requestId
}

// SetRequestID sets the request ID for this point
func (p *Point) SetRequestID(requestId string) {
	p.requestId = requestId
}

// Coordinate returns the geographic coordinate of this point
func (p *Point) Coordinate() *Coordinate {
	return p.point
}

// SetCoordinate sets the geographic coordinate for this point
func (p *Point) SetCoordinate(coordinate *Coordinate) {
	p.point = coordinate
}

// Time returns the time associated with this point
func (p *Point) Time() time.Time {
	return p.time
}

// SetTime sets the time for this point
func (p *Point) SetTime(time time.Time) {
	p.time = time
}

// PointType returns the type of this point
func (p *Point) PointType() PointType {
	return p.pointType
}

// SetPointType sets the type for this point
func (p *Point) SetPointType(pointType PointType) {
	p.pointType = pointType
}
