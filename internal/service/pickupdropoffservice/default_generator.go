package pickupdropoffservice

import (
	"fmt"
	"matching-engine/internal/enums"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
)

type DefaultGenerator struct {
	processor geo.GeospatialProcessor
}

func NewDefaultSelector(processor geo.GeospatialProcessor) *DefaultGenerator {
	return &DefaultGenerator{
		processor: processor,
	}
}
func (ds *DefaultGenerator) Generate(request *model.Request, offer *model.Offer) (pickup, dropoff *model.PathPoint, err error) {
	if request == nil || offer == nil {
		return nil, nil, fmt.Errorf("request or offer is nil")
	}
	// TODO: Implement the logic to generate pickup and dropoff points based on the walking time of the request
	// Initially, we set the pickup time to the earliest departure time of the request
	pickupPoint := model.NewPathPoint(request.Source(), enums.Pickup, request.EarliestDepartureTime(), request)
	// The dropoff time is set to the latest arrival time of the request
	dropoffPoint := model.NewPathPoint(offer.Destination(), enums.Dropoff, request.LatestArrivalTime(), offer)
	return pickupPoint, dropoffPoint, nil
}
