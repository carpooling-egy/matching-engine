package timematrix

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix/cache"
)

type CacheWithOfferIdPopulator struct {
	generator        Generator
	cacheWithOfferId *cache.TimeMatrixCacheWithOfferId
	cachingBound     int
}

func NewCacheWithOfferIdPopulator(generator Generator, cacheWithOfferId *cache.TimeMatrixCacheWithOfferId) *CacheWithOfferIdPopulator {
	return &CacheWithOfferIdPopulator{
		generator:        generator,
		cacheWithOfferId: cacheWithOfferId,
		cachingBound:     GetCachingBound(),
	}
}

func (p *CacheWithOfferIdPopulator) Populate(offer *model.OfferNode, requestNodes []*model.RequestNode) error {

	// early return if the number of request nodes exceeds the caching bound
	if len(requestNodes) > p.cachingBound {
		return nil
	}
	// Check if the time matrix is already cached
	_, exists := p.cacheWithOfferId.Get(offer.Offer().ID())
	if exists {
		return nil
	}

	// Create a new time matrix
	timeMatrix, err := p.generator.Generate(offer, requestNodes)
	if err != nil {
		return fmt.Errorf("could not generate time matrix for offer %s: %w", offer.Offer().ID(), err)
	}

	// Store the time matrix in the cacheWithOfferIdAndRequestId
	p.cacheWithOfferId.Set(offer.Offer().ID(), timeMatrix)
	return nil
}

func (p *CacheWithOfferIdPopulator) RemoveEntry(offer *model.OfferNode, requestNodes []*model.RequestNode) error {

	// Remove the time matrix from the cacheWithOfferId
	p.cacheWithOfferId.Delete(offer.Offer().ID())
	return nil
}
