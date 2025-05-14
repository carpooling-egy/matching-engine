package pickupdropoffservice

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"time"
)

type PickupDropoffSelector struct {
	// generator is the underlying generator that will be used to get the pickup and dropoff points
	generator PickupDropoffGenerator
	// engine is the routing engine that will be used to compute the walking time
	engine routing.Engine
	// cache is a map that will store the cached pickup and dropoff points
	cache *pickupdropoffcache.PickupDropoffCache
}

func NewPickupDropoffSelector(generator PickupDropoffGenerator, engine routing.Engine, cache *pickupdropoffcache.PickupDropoffCache) *PickupDropoffSelector {
	return &PickupDropoffSelector{
		generator: generator,
		engine:    engine,
		cache:     cache,
	}
}

func (selector *PickupDropoffSelector) computeWalkingDurations(request *model.Request, pickup, dropoff *model.PathPoint) (pickupWalkingDuration, dropoffWalkingDuration time.Duration, err error) {
	pickupWalkingParams, err := model.NewWalkParams(request.Source(), pickup.Coordinate())
	if err != nil {
		return 0, 0, fmt.Errorf("pickup walking params: %v", err)
	}
	pickupWalkingDuration, err = selector.engine.ComputeWalkingTime(context.Background(), pickupWalkingParams)
	if err != nil {
		return 0, 0, fmt.Errorf("pickup walking time: %v", err)
	}
	dropoffWalkingParams, err := model.NewWalkParams(dropoff.Coordinate(), request.Destination())
	if err != nil {
		return 0, 0, fmt.Errorf("dropoff walking params: %v", err)
	}
	dropoffWalkingDuration, err = selector.engine.ComputeWalkingTime(context.Background(), dropoffWalkingParams)
	if err != nil {
		return 0, 0, fmt.Errorf("dropoff walking time: %v", err)
	}
	return pickupWalkingDuration, dropoffWalkingDuration, nil
}

// GetPickupDropoffPointsAndDurations retrieves the pickup and dropoff points and durations for the given request and offer.
func (selector *PickupDropoffSelector) GetPickupDropoffPointsAndDurations(request *model.Request, offer *model.Offer) (value *pickupdropoffcache.Value, err error) {
	cacheKey := pickupdropoffcache.NewKey(
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
	pickupWalkingDuration, dropoffWalkingDuration, err := selector.computeWalkingDurations(request, pickup, dropoff)
	if err != nil {
		return nil, err
	}
	// Store the pickup and dropoff points in the cache
	cacheValue := pickupdropoffcache.NewValue(pickup, dropoff, pickupWalkingDuration, dropoffWalkingDuration)
	selector.cache.Set(cacheKey, cacheValue)
	// Return the pickup and dropoff points
	return cacheValue, nil
}
