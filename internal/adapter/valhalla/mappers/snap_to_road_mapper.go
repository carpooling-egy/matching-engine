package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
	"matching-engine/internal/model"
)

type SnapToRoadMapper struct{}

var _ re.OperationMapper[
	*model.Coordinate,
	*model.Coordinate,
	*pb.Api,
	*pb.Api,
] = &SnapToRoadMapper{}

func (SnapToRoadMapper) ToTransport(point *model.Coordinate) (*pb.Api, error) {
	if point == nil {
		return nil, fmt.Errorf("point cannot be nil")
	}

	return &pb.Api{
		Options: &pb.Options{
			Action: pb.Options_trace_attributes,
			Format: common.DefaultResponseFormat,
			Shape: []*pb.Location{ // same location twice to trick the engine to snap it to a road
				common.CreateLocation(point.Lat(), point.Lng(), pb.Location_kBreak),
				common.CreateLocation(point.Lat(), point.Lng(), pb.Location_kBreak),
			},
			CostingType: pb.Costing_auto_,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): common.DefaultAutoCosting,
			},
			ShapeMatch: pb.ShapeMatch_map_snap,
			HasSearchRadius: &pb.Options_SearchRadius{
				SearchRadius: float32(common.DefaultSearchRadiusInMeters),
			},
			PbfFieldSelector: &pb.PbfFieldSelector{
				Trip: true,
			},
		},
	}, nil
}

func (SnapToRoadMapper) FromTransport(response *pb.Api) (*model.Coordinate, error) {
	if response == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

	trip := response.GetTrip()
	if trip == nil {
		return nil, fmt.Errorf("trip cannot be nil")
	}

	routes := trip.GetRoutes()
	if len(routes) == 0 {
		return nil, fmt.Errorf("routes cannot be empty")
	}

	legs := routes[0].GetLegs()
	if len(legs) == 0 {
		return nil, fmt.Errorf("legs cannot be empty")
	}

	locations := legs[0].GetLocation()
	if len(locations) == 0 {
		return nil, fmt.Errorf("locations cannot be empty")
	}

	projectedLocation := locations[0].GetCorrelation().ProjectedLl
	if projectedLocation == nil {
		return nil, fmt.Errorf("projected location cannot be nil")
	}

	snappedPoint, err := model.NewCoordinate(
		projectedLocation.GetLat(),
		projectedLocation.GetLng(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create new coordinate: %w", err)
	}

	return snappedPoint, nil
}
