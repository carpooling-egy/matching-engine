package timematrix

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix/cache"
)

type DefaultPopulator struct {
	generator Generator
	cache     *cache.TimeMatrixCache
}

func NewDefaultPopulator(generator Generator, cache *cache.TimeMatrixCache) Populator {
	return &DefaultPopulator{
		generator: generator,
		cache:     cache,
	}
}

func (p *DefaultPopulator) Populate(offer *model.OfferNode, requestNodes []*model.RequestNode) error {

	// Check if the time matrix is already cached
	_, exists := p.cache.Get(offer.Offer().ID())
	if exists {
		return nil
	}

	// Create a new time matrix
	timeMatrix, err := p.generator.Generate(offer, requestNodes)
	if err != nil {
		return fmt.Errorf("could not generate time matrix for offer %s: %w", offer.Offer().ID(), err)
	}

	// Store the time matrix in the cache
	p.cache.Set(offer.Offer().ID(), timeMatrix)

	return nil
}
