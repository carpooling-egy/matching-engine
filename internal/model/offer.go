package model

import (
	"time"
)

// Offer represents a service provider's offer
type Offer struct {
	id            string
	userID        string
	source        Coordinate
	destination   Coordinate
	detourDurMins time.Duration
	departureTime time.Time
	preference    Preference
	capacity      int

	currentNumberOfRequests int
	matchedRequests         []*MatchedRequest
	path                    []*PathPoint
}

// No need to validate parameters as they will be read from database
// This constructor should be only used from database entities
func NewOffer(
	id, userID string,
	source, destination Coordinate,
	departureTime time.Time,
	detourDurMins time.Duration,
	capacity int,
	preference Preference,
	currentNumberOfRequests int,
	path []*PathPoint,
	matchedRequests []*MatchedRequest,
) *Offer {
	if matchedRequests == nil {
		matchedRequests = make([]*MatchedRequest, 0)
	}
	if path == nil {
		path = make([]*PathPoint, 0)
	}
	return &Offer{
		id:                      id,
		userID:                  userID,
		source:                  source,
		destination:             destination,
		detourDurMins:           detourDurMins,
		departureTime:           departureTime,
		capacity:                capacity,
		matchedRequests:         matchedRequests,
		preference:              preference,
		currentNumberOfRequests: currentNumberOfRequests,
		path:                    path,
	}
}

// Getters for immutable fields
func (o *Offer) ID() string                           { return o.id }
func (o *Offer) UserID() string                       { return o.userID }
func (o *Offer) Source() Coordinate                   { return o.source }
func (o *Offer) Destination() Coordinate              { return o.destination }
func (o *Offer) DepartureTime() time.Time             { return o.departureTime }
func (o *Offer) DetourDurationMinutes() time.Duration { return o.detourDurMins }
func (o *Offer) Capacity() int                        { return o.capacity }
func (o *Offer) Preferences() Preference              { return o.preference }
func (o *Offer) CurrentNumberOfRequests() int         { return o.currentNumberOfRequests }
func (o *Offer) PathPoints() []*PathPoint             { return o.path }
func (o *Offer) MatchedRequests() []*MatchedRequest {
	return o.matchedRequests
}

// SetMatchedRequests sets the matched requests
func (o *Offer) SetMatchedRequests(matchedRequests []*MatchedRequest) {
	o.matchedRequests = matchedRequests
}

// Path returns the path
func (o *Offer) Path() []*PathPoint {
	return o.path
}

// SetPath sets the path
func (o *Offer) SetPath(path []*PathPoint) {
	o.path = path
}

func (o *Offer) AsOffer() (*Offer, bool) {
	return o, true
}

func (o *Offer) AsRequest() (*Request, bool) {
	return nil, false
}

func (o *Offer) SetCurrentNumberOfRequests(count int) {
	o.currentNumberOfRequests = count
}
