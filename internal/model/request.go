package model

import (
	"matching-engine/internal/enums"
	"time"
)

// Request represents a user's request for a service
type Request struct {
	id                    string
	userID                string
	source                Coordinate
	destination           Coordinate
	earliestDepartureTime time.Time
	latestArrivalTime     time.Time
	maxWalkingTime        time.Duration
	preference            Preference
	numberOfRiders        int
}

// NewRequest creates a new Request
func NewRequest(id, userID string, source, destination Coordinate, earliestDepartureTime, latestArrivalTime time.Time, maxWalkingTime time.Duration, preference Preference, numberOfRiders int) *Request {
	return &Request{
		id:                    id,
		userID:                userID,
		source:                source,
		destination:           destination,
		earliestDepartureTime: earliestDepartureTime,
		latestArrivalTime:     latestArrivalTime,
		maxWalkingTime:        maxWalkingTime,
		preference:            preference,
		numberOfRiders:        numberOfRiders,
	}
}

// GetID returns the request ID
func (r *Request) GetID() string {
	return r.id
}

// SetID sets the request ID
func (r *Request) SetID(id string) {
	r.id = id
}

// GetUserID returns the user ID
func (r *Request) GetUserID() string {
	return r.userID
}

// SetUserID sets the user ID
func (r *Request) SetUserID(userID string) {
	r.userID = userID
}

// GetSource returns the source coordinate
func (r *Request) GetSource() Coordinate {
	return r.source
}

// SetSource sets the source coordinate
func (r *Request) SetSource(source Coordinate) {
	r.source = source
}

// GetDestination returns the destination coordinate
func (r *Request) GetDestination() Coordinate {
	return r.destination
}

// SetDestination sets the destination coordinate
func (r *Request) SetDestination(destination Coordinate) {
	r.destination = destination
}

// GetEarliestDepartureTime returns the earliest departure time
func (r *Request) GetEarliestDepartureTime() time.Time {
	return r.earliestDepartureTime
}

// SetEarliestDepartureTime sets the earliest departure time
func (r *Request) SetEarliestDepartureTime(departureTime time.Time) {
	r.earliestDepartureTime = departureTime
}

// GetLatestArrivalTime returns the latest arrival time
func (r *Request) GetLatestArrivalTime() time.Time {
	return r.latestArrivalTime
}

// SetLatestArrivalTime sets the latest arrival time
func (r *Request) SetLatestArrivalTime(arrivalTime time.Time) {
	r.latestArrivalTime = arrivalTime
}

// GetMaxWalkingTime returns the maximum walking time
func (r *Request) GetMaxWalkingTime() time.Duration {
	return r.maxWalkingTime
}

// SetMaxWalkingTime sets the maximum walking time
func (r *Request) SetMaxWalkingTime(duration time.Duration) {
	r.maxWalkingTime = duration
}

// GetPreference returns the preference
func (r *Request) GetPreference() Preference {
	return r.preference
}

// SetPreference sets the preference
func (r *Request) SetPreference(preference Preference) {
	r.preference = preference
}

// GetNumberOfRiders returns the number of riders
func (r *Request) GetNumberOfRiders() int {
	return r.numberOfRiders
}

// SetNumberOfRiders sets the number of riders
func (r *Request) SetNumberOfRiders(count int) {
	r.numberOfRiders = count
}

// GetRoleType returns the role type of the request.
// This method identifies the request as a specific role type in the system.
func (r *Request) GetRoleType() enums.RoleType {
	return enums.Request
}
