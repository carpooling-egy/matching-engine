package timematrix

import (
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix/cache"
)

// Generator generates travel time matrices between points
type Generator interface {
	// Generate creates a time matrix for an offer and its potential requests in the system
	Generate(offer *model.OfferNode) (*cache.PathPointMappedTimeMatrix, error)
}
