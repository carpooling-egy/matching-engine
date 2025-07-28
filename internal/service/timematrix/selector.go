package timematrix

import (
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix/cache"
)

type Selector interface {
	// GetTimeMatrix retrieves the time matrix for a given offer.
	// It may return a cached result or generate a new one if not available.
	GetTimeMatrix(offer *model.OfferNode, request *model.RequestNode) (*cache.PathPointMappedTimeMatrix, error)
}
