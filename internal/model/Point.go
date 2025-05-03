package model

import (
	"github.com/google/uuid"
	"matching-engine/internal/enums"
	"time"
)

// Point represents a location point with owner, time, and type information
type Point struct {
	ownerID    uuid.UUID
	ownerType  enums.OwnerType
	coordinate Coordinate
	time       time.Time
	pointType  enums.PointType
}

// NewPoint creates a new Point
func NewPoint(ownerID uuid.UUID, ownerType enums.OwnerType, coordinate Coordinate, time time.Time, pointType enums.PointType) *Point {
	return &Point{
		ownerID:    ownerID,
		ownerType:  ownerType,
		coordinate: coordinate,
		time:       time,
		pointType:  pointType,
	}
}

// GetOwnerID returns the owner ID
func (p *Point) GetOwnerID() uuid.UUID {
	return p.ownerID
}

// SetOwnerID sets the owner ID
func (p *Point) SetOwnerID(ownerID uuid.UUID) {
	p.ownerID = ownerID
}

// GetOwnerType returns the owner type
func (p *Point) GetOwnerType() enums.OwnerType {
	return p.ownerType
}

// SetOwnerType sets the owner type
func (p *Point) SetOwnerType(ownerType enums.OwnerType) {
	p.ownerType = ownerType
}

// GetCoordinate returns the coordinate
func (p *Point) GetCoordinate() Coordinate {
	return p.coordinate
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
