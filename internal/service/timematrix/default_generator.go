package timematrix

import (
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/timematrix/cache"
)

type DefaultGenerator struct {
	engine                routing.Engine
	pickupDropoffSelector pickupdropoffservice.PickupDropoffSelectorInterface
}

func NewDefaultGenerator(engine routing.Engine, pickupDropoffSelector pickupdropoffservice.PickupDropoffSelectorInterface) *DefaultGenerator {
	return &DefaultGenerator{
		engine:                engine,
		pickupDropoffSelector: pickupDropoffSelector,
	}
}
func (ds *DefaultGenerator) Generate(offerNode *model.OfferNode, requestNodes []*model.RequestNode) (*cache.PathPointMappedTimeMatrix, error) {
	if requestNodes == nil || len(requestNodes) == 0 {
		return nil, fmt.Errorf("requestNodes cannot be nil or empty")
	}

	pointToIdMap := make(map[model.PathPointID]int)

	var matrixPoints []model.Coordinate

	for _, point := range offerNode.Offer().PathPoints() {
		matrixPoints = append(matrixPoints, *point.Coordinate())
		pointToIdMap[point.ID()] = len(matrixPoints) - 1
	}

	// Add request pickup and dropoff points
	for _, requestNode := range requestNodes {
		pickupDropoff, err := ds.pickupDropoffSelector.GetPickupDropoffPointsAndDurations(requestNode.Request(), offerNode.Offer())
		if err != nil {
			// TODO: check if error is related to API calls before returning an error from the generator
			return nil, fmt.Errorf("failed to get pickup/dropoff points for requestNode %s: %w", requestNode.Request().ID(), err)
		}
		pickup := pickupDropoff.Pickup()
		dropoff := pickupDropoff.Dropoff()

		matrixPoints = append(matrixPoints, *pickup.Coordinate())
		pointToIdMap[pickup.ID()] = len(matrixPoints) - 1

		matrixPoints = append(matrixPoints, *dropoff.Coordinate())
		pointToIdMap[dropoff.ID()] = len(matrixPoints) - 1

	}

	if len(matrixPoints) <= 2 {
		return nil, fmt.Errorf("not enough points to generate a distance/time matrix")
	}

	// Call the routing engine to get the distance and time matrix
	params, err := model.NewDistanceTimeMatrixParams(
		matrixPoints,
		model.ProfileAuto,
		model.WithDepartureTime(offerNode.Offer().DepartureTime()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create distance time matrix params: %w", err)
	}

	distanceTimeMatrix, err := ds.engine.ComputeDistanceTimeMatrix(nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to compute distance time matrix for offer %s with %d matrixPoints: %w",
			offerNode.Offer().ID(), len(matrixPoints), err)
	}

	return cache.NewPathPointMappedTimeMatrix(distanceTimeMatrix.Times(), pointToIdMap), nil
}
