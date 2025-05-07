package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
	"matching-engine/internal/model"
	"time"
)

type DrivingTimeMapper struct{}

var _ re.OperationMapper[
	*model.RouteParams,
	[]time.Duration,
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

func (DrivingTimeMapper) FromTransport(response *pb.Api) ([]time.Duration, error) {
	if response == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

	directions := response.GetDirections()
	if directions == nil {
		return nil, fmt.Errorf("directions cannot be nil")
	}

	routes := directions.GetRoutes()
	if len(routes) == 0 {
		return nil, fmt.Errorf("no routes found in the response")
	}

	legs := routes[0].GetLegs()
	durations := make([]time.Duration, len(legs)+1)

	durations[0] = 0
	for i, leg := range legs {
		timeInSeconds := leg.GetSummary().Time
		durations[i+1] = time.Duration(timeInSeconds * float64(time.Second))
	}

	return durations, nil
}
