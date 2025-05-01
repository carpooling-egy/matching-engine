package model

import (
	"github.com/google/uuid"
	"matching-engine/internal/enums"
	"time"
)

// Point represents a location point with owner, time, and type information
type Point struct {
	OwnerID    uuid.UUID
	OwnerType  enums.OwnerType
	Coordinate Coordinate
	Time       time.Time
	PointType  enums.PointType
}

// NewPoint creates a new Point
func NewPoint(ownerID uuid.UUID, ownerType enums.OwnerType, coordinate Coordinate, time time.Time, pointType enums.PointType) *Point {
	return &Point{
		OwnerID:    ownerID,
		OwnerType:  ownerType,
		Coordinate: coordinate,
		Time:       time,
		PointType:  pointType,
	}
}
