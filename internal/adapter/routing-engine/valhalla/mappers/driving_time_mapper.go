package mappers

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/adapter/routing-engine/valhalla/helpers"
	"matching-engine/internal/model"
	"time"
)

type DrivingTimeMapper struct{}

var _ re.OperationMapper[
	*model.RouteParams,
	time.Duration,
	*pb.Api,
	*pb.Api,
] = DrivingTimeMapper{}

func (d DrivingTimeMapper) ToTransport(params *model.RouteParams) (*pb.Api, error) {
	wps := params.Waypoints()
	locations := make([]*pb.Location, len(wps))
	for i, wp := range wps {
		locations[i] = helpers.CreateLocation(wp.Lat(), wp.Lng(), helpers.DefaultLocationType)
	}

	return &pb.Api{
		Options: &pb.Options{
			Action: pb.Options_route,
			Units:  helpers.DefaultUnit,
			Format: helpers.DefaultFormat,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): helpers.DefaultAutoCosting,
			},
			HasShapeFormat: &pb.Options_ShapeFormat{
				ShapeFormat: helpers.DefaultShapeFormat,
			},
			Locations:    locations,
			DateTimeType: pb.Options_depart_at,
			HasDateTime: &pb.Options_DateTime{
				// TODO check how valhalla handles timezones
				DateTime: params.DepartureTime().String(),
			},
			PbfFieldSelector: &pb.PbfFieldSelector{
				Directions: true,
			},
		},
	}, nil
}

func (d DrivingTimeMapper) FromTransport(api *pb.Api) (time.Duration, error) {
	var timeInSeconds float64 = 0

	for _, leg := range api.GetDirections().GetRoutes()[0].GetLegs() {
		timeInSeconds += leg.GetSummary().Time
	}

	return time.Duration(timeInSeconds * float64(time.Second)), nil
}
