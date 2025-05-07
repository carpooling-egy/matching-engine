package model

import (
	"matching-engine/internal/enums"
	"time"
)

// Point represents a location point with owner, time, and type information
type Point struct {
	owner      Role
	coordinate Coordinate
	time       time.Time
	pointType  enums.PointType
}

// NewPoint creates a new Point
func NewPoint(owner Role, coordinate Coordinate, time time.Time, pointType enums.PointType) *Point {
	return &Point{
		owner:      owner,
		coordinate: coordinate,
		time:       time,
		pointType:  pointType,
	}
}

// Owner returns the owner of the point
func (p *Point) Owner() Role {
	return p.owner
}

// SetOwner sets the owner of the point
func (p *Point) SetOwner(owner Role) {
	p.owner = owner
}

// Coordinate returns the coordinate
func (p *Point) Coordinate() *Coordinate {
	return &p.coordinate
}

// SetCoordinate sets the coordinate
func (p *Point) SetCoordinate(coordinate Coordinate) {
	p.coordinate = coordinate
}

// Time returns the time
func (p *Point) Time() time.Time {
	return p.time
}

// SetTime sets the time
func (p *Point) SetTime(timestamp time.Time) {
	p.time = timestamp
}

// PointType returns the point type
func (p *Point) PointType() enums.PointType {
	return p.pointType
}

// SetPointType sets the point type
func (p *Point) SetPointType(pointType enums.PointType) {
	p.pointType = pointType
}
