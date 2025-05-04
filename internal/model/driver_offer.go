package models

import (
	"time"
)

// DriverOffer represents a driver offer entity
type DriverOffer struct {
	id                    string
	userID                string
	source                Coordinate
	destination           Coordinate
	departureTime         time.Time
	detourDurationMinutes time.Duration
	capacity              int
	preferences           Preference

	currentNumberOfRequests int
	pathPoints              []PathPoint
}

// No need to validate parameters as they will be read from database
// This constructor should be only used from database entities
func NewDriverOffer(
	id, userID string,
	source, destination Coordinate,
	departureTime time.Time,
	detourDurationMinutes time.Duration,
	capacity int,
	preferences Preference,
	currentNumberOfRequests int,
	pathPoints []PathPoint,
) *DriverOffer {
	return &DriverOffer{
		id:                      id,
		userID:                  userID,
		source:                  source,
		destination:             destination,
		departureTime:           departureTime,
		detourDurationMinutes:   detourDurationMinutes,
		capacity:                capacity,
		preferences:             preferences,
		currentNumberOfRequests: currentNumberOfRequests,
		pathPoints:              pathPoints,
	}
}

// Getters for immutable fields
func (d *DriverOffer) ID() string                           { return d.id }
func (d *DriverOffer) UserID() string                       { return d.userID }
func (d *DriverOffer) Source() Coordinate                   { return d.source }
func (d *DriverOffer) Destination() Coordinate              { return d.destination }
func (d *DriverOffer) DepartureTime() time.Time             { return d.departureTime }
func (d *DriverOffer) DetourDurationMinutes() time.Duration { return d.detourDurationMinutes }
func (d *DriverOffer) Capacity() int                        { return d.capacity }
func (d *DriverOffer) Preferences() Preference              { return d.preferences }
func (d *DriverOffer) CurrentNumberOfRequests() int         { return d.currentNumberOfRequests }
func (d *DriverOffer) PathPoints() []PathPoint              { return d.pathPoints }

// Methods for controlled mutation
func (d *DriverOffer) SetPathPoints(points []PathPoint) {
	d.pathPoints = points
}

func (d *DriverOffer) SetCurrentNumberOfRequests(count int) {
	d.currentNumberOfRequests = count
}
