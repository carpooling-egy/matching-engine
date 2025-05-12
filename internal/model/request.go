package model

import (
	"time"
)

type Request struct {
	id                        string
	userID                    string
	source                    Coordinate
	destination               Coordinate
	earliestDepartureTime     time.Time
	latestArrivalTime         time.Time
	maxWalkingDurationMinutes time.Duration
	numberOfRiders            int
	preferences               Preference
}

// NewRequest Creates a new Request, no need to validate parameters as they will be read from database
// This constructor should be only used from database entities
func NewRequest(
	id string,
	userID string,
	source Coordinate,
	destination Coordinate,
	earliestDepartureTime time.Time,
	latestArrivalTime time.Time,
	maxWalkingDurationMinutes time.Duration,
	numberOfRiders int,
	preferences Preference,
) *Request {
	return &Request{
		id:                        id,
		userID:                    userID,
		source:                    source,
		destination:               destination,
		earliestDepartureTime:     earliestDepartureTime,
		latestArrivalTime:         latestArrivalTime,
		maxWalkingDurationMinutes: maxWalkingDurationMinutes,
		numberOfRiders:            numberOfRiders,
		preferences:               preferences,
	}
}

// Getters
func (r *Request) ID() string                               { return r.id }
func (r *Request) UserID() string                           { return r.userID }
func (r *Request) Source() Coordinate                       { return r.source }
func (r *Request) Destination() Coordinate                  { return r.destination }
func (r *Request) EarliestDepartureTime() time.Time         { return r.earliestDepartureTime }
func (r *Request) LatestArrivalTime() time.Time             { return r.latestArrivalTime }
func (r *Request) MaxWalkingDurationMinutes() time.Duration { return r.maxWalkingDurationMinutes }
func (r *Request) NumberOfRiders() int                      { return r.numberOfRiders }
func (r *Request) Preferences() Preference                  { return r.preferences }

func (r *Request) AsOffer() (*Offer, bool) {
	return nil, false
}

func (r *Request) AsRequest() (*Request, bool) {
	return r, true
}
