package pickupdropoffservice

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/collections"
	"matching-engine/internal/enums"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/model"
	"time"
)

var _ PickupDropoffGenerator = (*IntersectionBasedGenerator)(nil)

type IntersectionBasedGenerator struct {
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

func (g *IntersectionBasedGenerator) getPickupDropoffPoint(
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

func (g *IntersectionBasedGenerator) GeneratePickupDropoffPoints(request *model.Request, offer *model.Offer) (pickup, dropoff *model.PathPoint, err error) {
	if request == nil || offer == nil {
		return nil, nil, fmt.Errorf("request or offer is nil")
	}
	geospatialProcessor, exists := g.offerProcessorCache.Get(offer.ID())
	if !exists {
		geospatialProcessor, err = g.processorFactory.CreateProcessor(offer)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create processor: %w", err)
		}
		g.offerProcessorCache.Set(offer.ID(), geospatialProcessor)
	}

	pickup, err = g.getPickupDropoffPoint(
		geospatialProcessor,
		request.Source(),
		enums.Pickup,
		request.EarliestDepartureTime(),
		request,
	)
	if err != nil {
		return nil, nil, err
	}
	dropoff, err = g.getPickupDropoffPoint(
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
