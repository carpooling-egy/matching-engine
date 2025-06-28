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

type SnappedSourceDestinationGenerator struct {
	offerProcessorCache *collections.SyncMap[string, processor.GeospatialProcessor]
	routingEngine       routing.Engine
}

func NewSnappedSourceDestinationGenerator(engine routing.Engine) PickupDropoffGenerator {
	return &SnappedSourceDestinationGenerator{
		offerProcessorCache: collections.NewSyncMap[string, processor.GeospatialProcessor](),
		routingEngine:       engine,
	}
}

func (g *SnappedSourceDestinationGenerator) getPickupDropoffPoint(
	coord *model.Coordinate,
	pointType enums.PointType,
	timeValue time.Time,
	request *model.Request,
) (*model.PathPoint, error) {
	zeroWalkingDuration := 0 * time.Minute
	snappedCoord, snapErr := g.routingEngine.SnapPointToRoad(context.Background(), coord)
	if snapErr != nil {
		return nil, fmt.Errorf("failed to snap point to road: %w", snapErr)
	}
	return model.NewPathPoint(*snappedCoord, pointType, timeValue, request, zeroWalkingDuration), nil
	
}

func (g *SnappedSourceDestinationGenerator) GeneratePickupDropoffPoints(request *model.Request, offer *model.Offer) (pickup, dropoff *model.PathPoint, err error) {
	if request == nil || offer == nil {
		return nil, nil, fmt.Errorf("request or offer is nil")
	}

	pickup, err = g.getPickupDropoffPoint(
		request.Source(),
		enums.Pickup,
		request.EarliestDepartureTime(),
		request,
	)
	if err != nil {
		return nil, nil, err
	}
	dropoff, err = g.getPickupDropoffPoint(
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
