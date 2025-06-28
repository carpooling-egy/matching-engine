package timematrix

import (
	"fmt"
	"matching-engine/internal/model"

	"matching-engine/internal/service/timematrix/cache"
)

type DefaultSelector struct {
	cacheWithOfferId             *cache.TimeMatrixCacheWithOfferId
	cacheWithOfferIdAndRequestId *cache.TimeMatrixCacheWithOfferIdAndRequestId
}

func NewDefaultSelector(cacheWithOfferId *cache.TimeMatrixCacheWithOfferId, cacheWithOfferIdAndRequestId *cache.TimeMatrixCacheWithOfferIdAndRequestId) Selector {
	return &DefaultSelector{
		cacheWithOfferId:             cacheWithOfferId,
		cacheWithOfferIdAndRequestId: cacheWithOfferIdAndRequestId,
	}
}

func (selector *DefaultSelector) GetTimeMatrix(offer *model.OfferNode, requestNode *model.RequestNode) (*cache.PathPointMappedTimeMatrix, error) {

	// Check if the time matrix is already cached
	cachedValue, exists := selector.cacheWithOfferId.Get(offer.Offer().ID())
	if !exists {
		cachedValue, exists = selector.cacheWithOfferIdAndRequestId.Get(offer.Offer().ID(), requestNode.Request().ID())
		if !exists {
			return nil, fmt.Errorf("requestNode with request ID %s and offer ID %s not found in cache", requestNode.Request().ID(), offer.Offer().ID())
		}
	}
	return cachedValue, nil
}
