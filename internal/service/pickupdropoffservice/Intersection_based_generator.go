package pickupdropoffservice

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/enums"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/model"
	"time"
)

type IntersectionBasedGenerator struct {
	offerProcessor   *collections.SyncMap[string, processor.GeospatialProcessor]
	processorFactory processor.ProcessorFactory
}

func NewIntersectionBasedGenerator(factory processor.ProcessorFactory) *IntersectionBasedGenerator {
	return &IntersectionBasedGenerator{
		offerProcessor:   collections.NewSyncMap[string, processor.GeospatialProcessor](),
		processorFactory: factory,
	}
}

func (g *IntersectionBasedGenerator) GeneratePickupDropoffPoints(request *model.Request, offer *model.Offer) (pickup, dropoff *model.PathPoint, err error) {
	if request == nil || offer == nil {
		return nil, nil, fmt.Errorf("request or offer is nil")
	}
	geospatialProcessor, ok := g.offerProcessor.Get(offer.ID())
	if !ok {
		geospatialProcessor, err = g.processorFactory.CreateProcessor(offer)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create processor: %w", err)
		}
		// Store the geospatialProcessor in the map
		g.offerProcessor.Set(offer.ID(), geospatialProcessor)
	}
	noWalkingDuration := 0 * time.Minute
	pickupCoord, pickupDuration, err := geospatialProcessor.ComputeClosestRoutePoint(request.Source(), request.MaxWalkingDurationMinutes())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute closest route point: %w", err)
	}
	// Check intersection between the route and the source circle
	if pickupDuration > request.MaxWalkingDurationMinutes() {
		// If the pickup point is not within the max walking duration, use the original source
		pickup = model.NewPathPoint(*request.Source(), enums.Pickup, request.EarliestDepartureTime(), request, noWalkingDuration)
	} else {
		// If the pickup point is within the max walking duration, use the computed route point
		pickup = model.NewPathPoint(*pickupCoord, enums.Pickup, request.EarliestDepartureTime(), request, pickupDuration)
	}

	dropoffCoord, dropoffDuration, err := geospatialProcessor.ComputeClosestRoutePoint(request.Destination(), request.MaxWalkingDurationMinutes())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute closest route point: %w", err)
	}
	// Check the intersection between the route and the destination circle
	if dropoffDuration > request.MaxWalkingDurationMinutes() {
		// If the dropoff point is not within the max walking duration, use the original destination
		dropoff = model.NewPathPoint(*request.Destination(), enums.Dropoff, request.LatestArrivalTime(), request, noWalkingDuration)
	} else {
		// If the dropoff point is within the max walking duration, use the computed route point
		dropoff = model.NewPathPoint(*dropoffCoord, enums.Dropoff, request.LatestArrivalTime(), request, dropoffDuration)
	}
	return pickup, dropoff, nil
}
