package timematrix

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix/cache"
)

type CacheWithOfferIdRequestIdPopulator struct {
	generator                    Generator
	cacheWithOfferIdAndRequestId *cache.TimeMatrixCacheWithOfferIdAndRequestId
	CacheWithOfferId             *cache.TimeMatrixCacheWithOfferId
}

func NewCacheWithOfferIdRequestIdPopulator(generator Generator, cacheWithOfferIdAndRequestId *cache.TimeMatrixCacheWithOfferIdAndRequestId, CacheWithOfferId *cache.TimeMatrixCacheWithOfferId) *CacheWithOfferIdRequestIdPopulator {
	return &CacheWithOfferIdRequestIdPopulator{
		generator:                    generator,
		cacheWithOfferIdAndRequestId: cacheWithOfferIdAndRequestId,
		CacheWithOfferId:             CacheWithOfferId,
	}
}

func (p *CacheWithOfferIdRequestIdPopulator) Populate(offer *model.OfferNode, requestNodes []*model.RequestNode) error {

	_, exists := p.CacheWithOfferId.Get(offer.Offer().ID())
	if exists {
		return nil
	}

	if len(requestNodes) != 1 {
		return fmt.Errorf("requestNodes should be contain 1 node, got %d request nodes", len(requestNodes))
	}
	// Check if the time matrix is already cached
	_, exists = p.cacheWithOfferIdAndRequestId.Get(offer.Offer().ID(), requestNodes[0].Request().ID())
	if exists {
		return nil
	}

	// Create a new time matrix
	timeMatrix, err := p.generator.Generate(offer, requestNodes)

	if err != nil {
		return fmt.Errorf("could not generate time matrix for offer %s: %w", offer.Offer().ID(), err)
	}

	// Store the time matrix in the cacheWithOfferIdAndRequestId
	p.cacheWithOfferIdAndRequestId.Set(offer.Offer().ID(), requestNodes[0].Request().ID(), timeMatrix)
	return nil
}

func (p *CacheWithOfferIdRequestIdPopulator) RemoveEntry(offer *model.OfferNode, requestNodes []*model.RequestNode) error {
	if len(requestNodes) != 1 {
		return fmt.Errorf("requestNodes should contain exactly 1 node, got %d request nodes", len(requestNodes))
	}

	// Remove the time matrix from the cacheWithOfferIdAndRequestId
	p.cacheWithOfferIdAndRequestId.Delete(offer.Offer().ID(), requestNodes[0].Request().ID())
	return nil
}
