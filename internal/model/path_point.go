package model

import (
	"matching-engine/internal/enums"
	"time"
	// "matching-engine/internal/errors"
)

// PathPointType represents the type of path PathPoint (pickup or dropoff)
type PathPointType string

// PathPoint represents a PathPoint in a driver's path
type PathPoint struct {
	owner               Role
	coordinate          Coordinate
	pointType           enums.PointType
	expectedArrivalTime time.Time
}

func NewPathPoint(
	coordinate Coordinate, pointType enums.PointType, expectedArrivalTime time.Time, owner Role) *PathPoint {

	// TODO: Validate parameters
	return &PathPoint{
		coordinate:          coordinate,
		pointType:           pointType,
		expectedArrivalTime: expectedArrivalTime,
		owner:               owner,
	}
}

// Owner returns the owner of the PathPoint
func (p *PathPoint) Owner() Role {
	return p.owner
}

// SetOwner sets the owner of the PathPoint
func (p *PathPoint) SetOwner(owner Role) {
	p.owner = owner
}

// Coordinate returns the coordinate
func (p *PathPoint) Coordinate() *Coordinate {
	return &p.coordinate
}

// SetCoordinate sets the coordinate
func (p *PathPoint) SetCoordinate(coordinate Coordinate) {
	p.coordinate = coordinate
}

// ExpectedArrivalTime Time returns the time
func (p *PathPoint) ExpectedArrivalTime() time.Time {
	return p.expectedArrivalTime
}

// SetExpectedArrivalTime sets the time
func (p *PathPoint) SetExpectedArrivalTime(timestamp time.Time) {
	p.expectedArrivalTime = timestamp
}

// PointType returns the PathPoint type
func (p *PathPoint) PointType() enums.PointType {
	return p.pointType
}

// SetPointType sets the PathPoint type
func (p *PathPoint) SetPointType(pointType enums.PointType) {
	p.pointType = pointType
}
