package model

import (
	"time"
)

// Offer represents a service provider's offer
type Offer struct {
	id              string
	userID          string
	source          Coordinate
	destination     Coordinate
	detourTime      time.Duration
	departureTime   time.Time
	matchedRequests []*MatchedRequest
	preference      Preference
	path            []*Point
}

// NewOffer creates a new Offer
func NewOffer(id, userID string, source, destination Coordinate, detourTime time.Duration, departureTime time.Time, matchedRequests []*MatchedRequest, preference Preference, path []*Point) *Offer {
	if matchedRequests == nil {
		matchedRequests = make([]*MatchedRequest, 0)
	}
	if path == nil {
		path = make([]*Point, 0)
	}
	return &Offer{
		id:              id,
		userID:          userID,
		source:          source,
		destination:     destination,
		detourTime:      detourTime,
		departureTime:   departureTime,
		matchedRequests: matchedRequests,
		preference:      preference,
		path:            path,
	}
}

// ID returns the offer ID
func (o *Offer) ID() string {
	return o.id
}

// SetID sets the offer ID
func (o *Offer) SetID(id string) {
	o.id = id
}

// UserID returns the user ID
func (o *Offer) UserID() string {
	return o.userID
}

// SetUserID sets the user ID
func (o *Offer) SetUserID(userID string) {
	o.userID = userID
}

// Source returns the source coordinate
func (o *Offer) Source() Coordinate {
	return o.source
}

// SetSource sets the source coordinate
func (o *Offer) SetSource(source Coordinate) {
	o.source = source
}

// Destination returns the destination coordinate
func (o *Offer) Destination() *Coordinate {
	return &o.destination
}

// SetDestination sets the destination coordinate
func (o *Offer) SetDestination(destination Coordinate) {
	o.destination = destination
}

// DetourTime returns the detour time
func (o *Offer) DetourTime() time.Duration {
	return o.detourTime
}

// SetDetourTime sets the detour time
func (o *Offer) SetDetourTime(detourTime time.Duration) {
	o.detourTime = detourTime
}

// DepartureTime returns the departure time
func (o *Offer) DepartureTime() time.Time {
	return o.departureTime
}

// SetDepartureTime sets the departure time
func (o *Offer) SetDepartureTime(departureTime time.Time) {
	o.departureTime = departureTime
}

// MatchedRequests returns the matched requests
func (o *Offer) MatchedRequests() []*MatchedRequest {
	return o.matchedRequests
}

// SetMatchedRequests sets the matched requests
func (o *Offer) SetMatchedRequests(matchedRequests []*MatchedRequest) {
	o.matchedRequests = matchedRequests
}

// Preference returns the preference
func (o *Offer) Preference() Preference {
	return o.preference
}

// SetPreference sets the preference
func (o *Offer) SetPreference(preference Preference) {
	o.preference = preference
}

// Path returns the path
func (o *Offer) Path() []*Point {
	return o.path
}

// SetPath sets the path
func (o *Offer) SetPath(path []*Point) {
	o.path = path
}

func (o *Offer) AsOffer() (*Offer, bool) {
	return o, true
}

func (o *Offer) AsRequest() (*Request, bool) {
	return nil, false
}
