package models

import (
    "time"	
)

type RiderRequest struct {
    id                        string
    userID                    string
    source                    Coordinate
    destination               Coordinate
    earliestDepartureTime     time.Time
    latestArrivalTime         time.Time
    maxWalkingDurationMinutes time.Duration
    numberOfRiders            int
    preferences               Preference
    isMatched                 bool
}

// No need to validate parameters as they will be read from database
// This constructor should be only used from database entities
func NewRiderRequest(
    id string,
    userID string,
    source Coordinate,
    destination Coordinate,
    earliestDepartureTime time.Time,
    latestArrivalTime time.Time,
    maxWalkingDurationMinutes time.Duration,
    numberOfRiders int,
    preferences Preference,
	isMatched bool,
) (*RiderRequest) {
    return &RiderRequest{
        id:                        id,
        userID:                    userID,
        source:                    source,
        destination:               destination,
        earliestDepartureTime:     earliestDepartureTime,
        latestArrivalTime:         latestArrivalTime,
        maxWalkingDurationMinutes: maxWalkingDurationMinutes,
        numberOfRiders:            numberOfRiders,
        preferences:               preferences,
        isMatched:                 false,
    }
}

// Getters
func (r *RiderRequest) ID() string                       { return r.id }
func (r *RiderRequest) UserID() string                   { return r.userID }
func (r *RiderRequest) Source() Coordinate               { return r.source }
func (r *RiderRequest) Destination() Coordinate          { return r.destination }
func (r *RiderRequest) EarliestDepartureTime() time.Time { return r.earliestDepartureTime }
func (r *RiderRequest) LatestArrivalTime() time.Time     { return r.latestArrivalTime }
func (r *RiderRequest) MaxWalkingDurationMinutes() time.Duration { return r.maxWalkingDurationMinutes }
func (r *RiderRequest) NumberOfRiders() int              { return r.numberOfRiders }
func (r *RiderRequest) Preferences() Preference          { return r.preferences }
func (r *RiderRequest) IsMatched() bool                  { return r.isMatched }

// WithMatched returns a new RiderRequest with isMatched set to the given value
func (r *RiderRequest) WithMatched(isMatched bool) *RiderRequest {
    newRequest := *r
    newRequest.isMatched = isMatched
    return &newRequest
}