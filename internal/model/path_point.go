package model

import (
	"matching-engine/internal/enums"
	"sync/atomic"
	"time"
)

// PathPointID represents a unique identifier for a path point
type PathPointID int64

// Global atomic counter for thread-safe ID generation
var NextPointID int64 = 1

// PathPoint represents a PathPoint in a driver's path
type PathPoint struct {
	id                  PathPointID
	owner               Role
	coordinate          Coordinate
	pointType           enums.PointType
	expectedArrivalTime time.Time
	walkingDuration     time.Duration
}

func NewPathPoint(
	coordinate Coordinate, pointType enums.PointType, expectedArrivalTime time.Time, owner Role, walkingDuration time.Duration) *PathPoint {

	id := atomic.AddInt64(&NextPointID, 1) - 1

	// TODO: Validate parameters
	return &PathPoint{
		id:                  PathPointID(id),
		coordinate:          coordinate,
		pointType:           pointType,
		expectedArrivalTime: expectedArrivalTime,
		owner:               owner,
		walkingDuration:     walkingDuration,
	}
}

// ID returns the ID of the PathPoint
func (p *PathPoint) ID() PathPointID {
	return p.id
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

// WalkingDuration returns the walking duration
func (p *PathPoint) WalkingDuration() time.Duration {
	return p.walkingDuration
}

// SetWalkingDuration sets the walking duration
func (p *PathPoint) SetWalkingDuration(duration time.Duration) {
	p.walkingDuration = duration
}
