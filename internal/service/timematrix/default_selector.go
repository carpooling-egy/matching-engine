package timematrix

import (
	"fmt"
	"matching-engine/internal/model"

	"matching-engine/internal/service/timematrix/cache"
)

type DefaultSelector struct {
	cache *cache.TimeMatrixCache
}

func NewDefaultSelector(cache *cache.TimeMatrixCache) *DefaultSelector {
	return &DefaultSelector{
		cache: cache,
	}
}

func (selector *DefaultSelector) GetTimeMatrix(offer *model.OfferNode) (*cache.PathPointMappedTimeMatrix, error) {

	// Check if the time matrix is already cached
	cachedValue, exists := selector.cache.Get(offer.Offer().ID())
	if !exists {
		return nil, fmt.Errorf("offer %s not found in cache, should have been populated", offer.Offer().ID())
	}
	return cachedValue, nil
}
