package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
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

func (DrivingTimeMapper) ToTransport(params *model.RouteParams) (*pb.Api, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	wps := params.Waypoints()
	locations := make([]*pb.Location, len(wps))
	for i, wp := range wps {
		locations[i] = common.CreateLocation(wp.Lat(), wp.Lng())
	}

	return &pb.Api{
		Options: &pb.Options{
			Action:      pb.Options_route,
			Units:       common.DefaultUnit,
			Format:      common.DefaultFormat,
			CostingType: pb.Costing_auto_,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): common.DefaultAutoCosting,
			},
			HasShapeFormat: &pb.Options_ShapeFormat{
				ShapeFormat: common.DefaultShapeFormat,
			},
			Locations:    locations,
			DateTimeType: pb.Options_depart_at,
			HasDateTime: &pb.Options_DateTime{
				// TODO check how valhalla handles timezones
				DateTime: params.DepartureTime().Format("2006-01-02T15:04"),
			},
			PbfFieldSelector: &pb.PbfFieldSelector{
				Directions: true,
			},
		},
	}, nil
}

func (DrivingTimeMapper) FromTransport(response *pb.Api) (time.Duration, error) {
	if response == nil {
		return 0, fmt.Errorf("response cannot be nil")
	}

	var timeInSeconds float64 = 0

	for _, leg := range response.GetDirections().GetRoutes()[0].GetLegs() {
		timeInSeconds += leg.GetSummary().Time
	}

	return time.Duration(timeInSeconds * float64(time.Second)), nil
}
