package model

import (
	"matching-engine/internal/enums"
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

// GetID returns the offer ID
func (o *Offer) GetID() string {
	return o.id
}

// SetID sets the offer ID
func (o *Offer) SetID(id string) {
	o.id = id
}

// GetUserID returns the user ID
func (o *Offer) GetUserID() string {
	return o.userID
}

// SetUserID sets the user ID
func (o *Offer) SetUserID(userID string) {
	o.userID = userID
}

// GetSource returns the source coordinate
func (o *Offer) GetSource() Coordinate {
	return o.source
}

// SetSource sets the source coordinate
func (o *Offer) SetSource(source Coordinate) {
	o.source = source
}

// GetDestination returns the destination coordinate
func (o *Offer) GetDestination() Coordinate {
	return o.destination
}

// SetDestination sets the destination coordinate
func (o *Offer) SetDestination(destination Coordinate) {
	o.destination = destination
}

// GetDetourTime returns the detour time
func (o *Offer) GetDetourTime() time.Duration {
	return o.detourTime
}

// SetDetourTime sets the detour time
func (o *Offer) SetDetourTime(detourTime time.Duration) {
	o.detourTime = detourTime
}

// GetDepartureTime returns the departure time
func (o *Offer) GetDepartureTime() time.Time {
	return o.departureTime
}

// SetDepartureTime sets the departure time
func (o *Offer) SetDepartureTime(departureTime time.Time) {
	o.departureTime = departureTime
}

// GetMatchedRequests returns the matched requests
func (o *Offer) GetMatchedRequests() []*MatchedRequest {
	return o.matchedRequests
}

// SetMatchedRequests sets the matched requests
func (o *Offer) SetMatchedRequests(matchedRequests []*MatchedRequest) {
	o.matchedRequests = matchedRequests
}

// GetPreference returns the preference
func (o *Offer) GetPreference() Preference {
	return o.preference
}

// SetPreference sets the preference
func (o *Offer) SetPreference(preference Preference) {
	o.preference = preference
}

// GetPath returns the path
func (o *Offer) GetPath() []*Point {
	return o.path
}

// SetPath sets the path
func (o *Offer) SetPath(path []*Point) {
	o.path = path
}

// GetRoleType returns the role type of the offer
func (o *Offer) GetRoleType() enums.RoleType {
	return enums.Offer
}
