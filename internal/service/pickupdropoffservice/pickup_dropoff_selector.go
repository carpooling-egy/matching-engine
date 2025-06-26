package pickupdropoffservice

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
)

type PickupDropoffSelector struct {
	// generator is the underlying generator that will be used to get the pickup and dropoff points
	generator PickupDropoffGenerator
	// cache is a map that will store the cached pickup and dropoff points
	cache *pickupdropoffcache.PickupDropoffCache
}

func NewPickupDropoffSelector(generator PickupDropoffGenerator, cache *pickupdropoffcache.PickupDropoffCache) PickupDropoffSelectorInterface {
	return &PickupDropoffSelector{
		generator: generator,
		cache:     cache,
	}
}

// GetPickupDropoffPointsAndDurations retrieves the pickup and dropoff points and durations for the given request and offer.
func (selector *PickupDropoffSelector) GetPickupDropoffPointsAndDurations(request *model.Request, offer *model.Offer) (value *pickupdropoffcache.Value, err error) {
	cacheKey := model.NewOfferRequestKey(
		offer.ID(),
		request.ID())
	if cachedValue, ok := selector.cache.Get(cacheKey); ok {
		return cachedValue, nil
	}

	// Call the underlying generator to get the pickup and dropoff points
	pickup, dropoff, err := selector.generator.GeneratePickupDropoffPoints(request, offer)
	if err != nil {
		return nil, fmt.Errorf("pickup dropoff generator error: %v", err)
	}

	// Set the expected arrival times:
	// - For the pickup point: earliest pickup time (departure time + walking duration)
	// - For the dropoff point: latest dropoff time (latest arrival time - walking duration)
	// NOTE: Be careful when changing these as some path generation logic depends on it
	pickup.SetExpectedArrivalTime(request.EarliestDepartureTime().Add(pickup.WalkingDuration()))
	dropoff.SetExpectedArrivalTime(request.LatestArrivalTime().Add(-dropoff.WalkingDuration()))

	// Store the pickup and dropoff points in the cache
	cacheValue := pickupdropoffcache.NewValue(pickup, dropoff)
	selector.cache.Set(cacheKey, cacheValue)

	// Return the pickup and dropoff points
	return cacheValue, nil
}
