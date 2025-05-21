package pickupdropoffservice

import (
	"context"
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
)

type PickupDropoffSelector struct {
	// generator is the underlying generator that will be used to get the pickup and dropoff points
	generator             PickupDropoffGenerator
	walkingTimeCalculator *WalkingTimeCalculator
	// cache is a map that will store the cached pickup and dropoff points
	cache *pickupdropoffcache.PickupDropoffCache
}

func NewPickupDropoffSelector(generator PickupDropoffGenerator, walkingTimeCalculator *WalkingTimeCalculator, cache *pickupdropoffcache.PickupDropoffCache) PickupDropoffSelectorInterface {
	return &PickupDropoffSelector{
		generator:             generator,
		walkingTimeCalculator: walkingTimeCalculator,
		cache:                 cache,
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
	pickupWalkingDuration, dropoffWalkingDuration, err := selector.walkingTimeCalculator.ComputeWalkingDurations(context.Background(), request, pickup, dropoff)
	if err != nil {
		return nil, err
	}

	// Set the walking durations for the pickup and dropoff points
	pickup.SetWalkingDuration(pickupWalkingDuration)
	dropoff.SetWalkingDuration(dropoffWalkingDuration)

	// Set the expected arrival times:
	// - For the pickup point: earliest pickup time (departure time + walking duration)
	// - For the dropoff point: latest dropoff time (latest arrival time - walking duration)
	// NOTE: Be careful when changing these as some path generation logic depend on it
	pickup.SetExpectedArrivalTime(request.EarliestDepartureTime().Add(pickupWalkingDuration))
	dropoff.SetExpectedArrivalTime(request.LatestArrivalTime().Add(-dropoffWalkingDuration))

	// Store the pickup and dropoff points in the cache
	cacheValue := pickupdropoffcache.NewValue(pickup, dropoff)
	selector.cache.Set(cacheKey, cacheValue)

	// Return the pickup and dropoff points
	return cacheValue, nil
}
