package timematrix

import (
	"matching-engine/internal/model"

	"matching-engine/internal/service/timematrix/cache"
)

type DefaultSelector struct {
	generator Generator
	cache     *cache.TimeMatrixCache
}

func NewDefaultSelector(generator Generator, cache *cache.TimeMatrixCache) *DefaultSelector {
	return &DefaultSelector{
		generator: generator,
		cache:     cache,
	}
}

func (selector *DefaultSelector) GetTimeMatrix(offer *model.OfferNode) (*cache.PathPointMappedTimeMatrix, error) {

	// Check if the time matrix is already cached
	if cachedValue, ok := selector.cache.Get(offer.Offer().ID()); ok {
		return cachedValue, nil
	}

	// Call the underlying generator to get the time matrix and path point ID to index mapping
	pointMappedMatrix, err := selector.generator.Generate(offer)
	if err != nil {
		return nil, err
	}

	// Store the time matrix and path point ID to index mapping in the cache
	selector.cache.Set(offer.Offer().ID(), pointMappedMatrix)

	return pointMappedMatrix, nil
}
