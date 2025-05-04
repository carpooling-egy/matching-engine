package models

import (
	"time"
	// "matching-engine/internal/errors"
)

// PointType represents the type of path point (pickup or dropoff)
type PointType string

const (
	PointTypePickup  PointType = "pickup"
	PointTypeDropoff PointType = "dropoff"
)

// PathPoint represents a point in a driver's path
type PathPoint struct {
	coordinate          Coordinate
	pointType           PointType
	expectedArrivalTime time.Time
	riderRequest        *RiderRequest
}

func NewPathPoint(
    coordinate Coordinate, pointType PointType, expectedArrivalTime time.Time, riderRequest *RiderRequest) *PathPoint {
    
	// TODO: Validate parameters
    return &PathPoint{
        coordinate:          coordinate,
        pointType:           pointType,
        expectedArrivalTime: expectedArrivalTime,
        riderRequest:        riderRequest,
    }
}

func (p *PathPoint) Coordinate() Coordinate {
	return p.coordinate
}
func (p *PathPoint) PointType() PointType {
	return p.pointType
}
func (p *PathPoint) ExpectedArrivalTime() time.Time {
	return p.expectedArrivalTime
}
func (p *PathPoint) RiderRequest() *RiderRequest {
	return p.riderRequest
}
