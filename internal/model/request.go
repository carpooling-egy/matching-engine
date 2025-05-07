package model

import (
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

// ID returns the request ID
func (r *Request) ID() string {
	return r.id
}

// SetID sets the request ID
func (r *Request) SetID(id string) {
	r.id = id
}

// UserID returns the user ID
func (r *Request) UserID() string {
	return r.userID
}

// SetUserID sets the user ID
func (r *Request) SetUserID(userID string) {
	r.userID = userID
}

// Source returns the source coordinate
func (r *Request) Source() *Coordinate {
	return &r.source
}

// SetSource sets the source coordinate
func (r *Request) SetSource(source Coordinate) {
	r.source = source
}

// Destination returns the destination coordinate
func (r *Request) Destination() *Coordinate {
	return &r.destination
}

// SetDestination sets the destination coordinate
func (r *Request) SetDestination(destination Coordinate) {
	r.destination = destination
}

// EarliestDepartureTime returns the earliest departure time
func (r *Request) EarliestDepartureTime() time.Time {
	return r.earliestDepartureTime
}

// SetEarliestDepartureTime sets the earliest departure time
func (r *Request) SetEarliestDepartureTime(departureTime time.Time) {
	r.earliestDepartureTime = departureTime
}

// LatestArrivalTime returns the latest arrival time
func (r *Request) LatestArrivalTime() time.Time {
	return r.latestArrivalTime
}

// SetLatestArrivalTime sets the latest arrival time
func (r *Request) SetLatestArrivalTime(arrivalTime time.Time) {
	r.latestArrivalTime = arrivalTime
}

// MaxWalkingTime returns the maximum walking time
func (r *Request) MaxWalkingTime() time.Duration {
	return r.maxWalkingTime
}

// SetMaxWalkingTime sets the maximum walking time
func (r *Request) SetMaxWalkingTime(duration time.Duration) {
	r.maxWalkingTime = duration
}

// Preference returns the preference
func (r *Request) Preference() Preference {
	return r.preference
}

// SetPreference sets the preference
func (r *Request) SetPreference(preference Preference) {
	r.preference = preference
}

// NumberOfRiders returns the number of riders
func (r *Request) NumberOfRiders() int {
	return r.numberOfRiders
}

// SetNumberOfRiders sets the number of riders
func (r *Request) SetNumberOfRiders(count int) {
	r.numberOfRiders = count
}

func (r *Request) AsOffer() (*Offer, bool) {
	return nil, false
}

func (r *Request) AsRequest() (*Request, bool) {
	return r, true
}
