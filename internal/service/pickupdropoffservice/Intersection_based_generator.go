package pickupdropoffservice

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/collections"
	"matching-engine/internal/enums"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/model"
	"time"
)

var _ PickupDropoffGenerator = (*IntersectionBasedGenerator)(nil)

type IntersectionBasedGenerator struct {
	group               singleflight.Group
	offerProcessorCache *collections.SyncMap[string, processor.GeospatialProcessor]
	processorFactory    processor.ProcessorFactory
	routingEngine       routing.Engine
}

func NewIntersectionBasedGenerator(factory processor.ProcessorFactory, engine routing.Engine) PickupDropoffGenerator {
	return &IntersectionBasedGenerator{
		offerProcessorCache: collections.NewSyncMap[string, processor.GeospatialProcessor](),
		processorFactory:    factory,
		routingEngine:       engine,
	}
}

func (g *IntersectionBasedGenerator) computePickupDropoffPoint(
	geospatialProcessor processor.GeospatialProcessor,
	coord *model.Coordinate,
	pointType enums.PointType,
	timeValue time.Time,
	request *model.Request,
) (*model.PathPoint, error) {
	zeroWalkingDuration := 0 * time.Minute
	computedCoord, duration, err := geospatialProcessor.ComputeClosestRoutePoint(coord, request.MaxWalkingDurationMinutes())
	if err != nil {
		return nil, fmt.Errorf("failed to compute closest route point: %w", err)
	}
	if duration > request.MaxWalkingDurationMinutes() {
		snappedCoord, snapErr := g.routingEngine.SnapPointToRoad(context.Background(), coord)
		if snapErr != nil {
			return nil, fmt.Errorf("failed to snap point to road: %w", snapErr)
		}
		return model.NewPathPoint(*snappedCoord, pointType, timeValue, request, zeroWalkingDuration), nil
	}
	return model.NewPathPoint(*computedCoord, pointType, timeValue, request, duration), nil
}

func (g *IntersectionBasedGenerator) GeneratePickupDropoffPoints(
	request *model.Request,
	offer *model.Offer,
) (pickup, dropoff *model.PathPoint, err error) {
	if request == nil || offer == nil {
		return nil, nil, fmt.Errorf("request or offer is nil")
	}

	// Note: This may initialize the processor more than once if two goroutines race here.
	// The second initialization simply overwrites the first (harmless for correctness)
	// but wastes CPU and memory by doing duplicate work (API calls and geo operations).

	/*
		geospatialProcessor, exists := g.offerProcessorCache.Get(offer.ID())
		if !exists {
			geospatialProcessor, err = g.processorFactory.CreateProcessor(offer)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to create processor: %w", err)
			}
			g.offerProcessorCache.Set(offer.ID(), geospatialProcessor)
		}
	*/

	// Fast-path cache check: if a processor was already created, use it directly.
	if geospatialProcessor, exists := g.offerProcessorCache.Get(offer.ID()); exists {
		return g.getPickupAndDropoff(geospatialProcessor, request)
	}

	// Slow-path initialization via singleflight:
	// singleflight.Group.Do ensures that even if multiple goroutines
	// reach this point at the same time for the same offer.ID(),
	// the enclosed function will only be executed once.
	// All callers with the same key will wait for that one execution
	// and receive the same returned value, eliminating duplicate work.
	v, err, _ := g.group.Do(offer.ID(), func() (any, error) {

		// Secondary cache check inside the closure:
		// It’s possible another goroutine completed initialization
		// and stored the processor between our fast-path check and now.
		// This avoids a redundant CreateProcessor call in that window.
		if proc, exists := g.offerProcessorCache.Get(offer.ID()); exists {
			return proc, nil
		}

		// Actual creation and caching of the processor:
		newProc, err := g.processorFactory.CreateProcessor(offer)
		if err != nil {
			return nil, err
		}
		g.offerProcessorCache.Set(offer.ID(), newProc)
		return newProc, nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create processor: %w", err)
	}

	geospatialProcessor := v.(processor.GeospatialProcessor)
	return g.getPickupAndDropoff(geospatialProcessor, request)
}

func (g *IntersectionBasedGenerator) getPickupAndDropoff(
	geospatialProcessor processor.GeospatialProcessor,
	request *model.Request,
) (*model.PathPoint, *model.PathPoint, error) {
	// Compute the pickup point
	pickup, err := g.computePickupDropoffPoint(
		geospatialProcessor,
		request.Source(),
		enums.Pickup,
		request.EarliestDepartureTime(),
		request,
	)
	if err != nil {
		return nil, nil, err
	}

	// Compute the dropoff point
	dropoff, err := g.computePickupDropoffPoint(
		geospatialProcessor,
		request.Destination(),
		enums.Dropoff,
		request.LatestArrivalTime(),
		request,
	)
	if err != nil {
		return nil, nil, err
	}
	return pickup, dropoff, nil
}
