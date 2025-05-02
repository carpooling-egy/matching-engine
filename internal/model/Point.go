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
