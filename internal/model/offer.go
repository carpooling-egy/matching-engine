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
func (d *Offer) ID() string                           { return d.id }
func (d *Offer) UserID() string                       { return d.userID }
func (d *Offer) Source() Coordinate                   { return d.source }
func (d *Offer) Destination() Coordinate              { return d.destination }
func (d *Offer) DepartureTime() time.Time             { return d.departureTime }
func (d *Offer) DetourDurationMinutes() time.Duration { return d.detourDurMins }
func (d *Offer) Capacity() int                        { return d.capacity }
func (d *Offer) Preferences() Preference              { return d.preference }
func (d *Offer) CurrentNumberOfRequests() int         { return d.currentNumberOfRequests }
func (d *Offer) PathPoints() []*PathPoint {return d.path}
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

func (d *Offer) SetCurrentNumberOfRequests(count int) {
	d.currentNumberOfRequests = count
}
