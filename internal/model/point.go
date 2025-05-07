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

// GetOwner returns the owner of the point
func (p *Point) GetOwner() Role {
	return p.owner
}

// SetOwner sets the owner of the point
func (p *Point) SetOwner(owner Role) {
	p.owner = owner
}

// GetCoordinate returns the coordinate
func (p *Point) GetCoordinate() *Coordinate {
	return &p.coordinate
}

// SetCoordinate sets the coordinate
func (p *Point) SetCoordinate(coordinate Coordinate) {
	p.coordinate = coordinate
}

// GetTime returns the time
func (p *Point) GetTime() time.Time {
	return p.time
}

// SetTime sets the time
func (p *Point) SetTime(timestamp time.Time) {
	p.time = timestamp
}

// GetPointType returns the point type
func (p *Point) GetPointType() enums.PointType {
	return p.pointType
}

// SetPointType sets the point type
func (p *Point) SetPointType(pointType enums.PointType) {
	p.pointType = pointType
}
