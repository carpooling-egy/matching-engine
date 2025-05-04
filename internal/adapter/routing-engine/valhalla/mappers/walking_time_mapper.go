package mappers

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/adapter/routing-engine/valhalla/helpers"
	"matching-engine/internal/model"
	"time"
)

type WalkingTimeMapper struct{}

var _ re.OperationMapper[
	*model.WalkParams,
	time.Duration,
	*pb.Api,
	*pb.Api,
] = WalkingTimeMapper{}

func (w WalkingTimeMapper) ToTransport(params *model.WalkParams) (*pb.Api, error) {

	wps := []*model.Coordinate{params.Origin(), params.Destination()}
	locations := make([]*pb.Location, 2)
	for i, wp := range wps {
		locations[i] = helpers.CreateLocation(wp.Lat(), wp.Lng(), helpers.DefaultLocationType)
	}

	return &pb.Api{
		Options: &pb.Options{
			Action: pb.Options_route,
			Units:  helpers.DefaultUnit,
			Format: helpers.DefaultFormat,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): helpers.DefaultPedestrianCosting,
			},
			HasShapeFormat: &pb.Options_ShapeFormat{
				ShapeFormat: helpers.DefaultShapeFormat,
			},
			Locations: locations,
			PbfFieldSelector: &pb.PbfFieldSelector{
				Directions: true,
			},
		},
	}, nil
}

func (w WalkingTimeMapper) FromTransport(api *pb.Api) (time.Duration, error) {
	var timeInSeconds float64 = 0

	for _, leg := range api.GetDirections().GetRoutes()[0].GetLegs() {
		timeInSeconds += leg.GetSummary().Time
	}

	return time.Duration(timeInSeconds * float64(time.Second)), nil
}
