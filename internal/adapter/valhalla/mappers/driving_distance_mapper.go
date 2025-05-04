package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
	"matching-engine/internal/model"
)

type DrivingDistanceMapper struct{}

var _ re.OperationMapper[
	*model.RouteParams,
	*model.Distance,
	*pb.Api,
	*pb.Api,
] = DrivingDistanceMapper{}

func (DrivingDistanceMapper) ToTransport(params *model.RouteParams) (*pb.Api, error) {
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

func (DrivingDistanceMapper) FromTransport(response *pb.Api) (*model.Distance, error) {
	if response == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

	var distance float32 = 0
	for _, leg := range response.GetDirections().GetRoutes()[0].GetLegs() {
		distance += leg.GetSummary().Length
	}

	unit, err := common.ToDomainDistanceUnit(common.DefaultUnit)
	if err != nil {
		return nil, err
	}

	distanceObj, err := model.NewDistance(distance, unit)
	if err != nil {
		return nil, err
	}

	return distanceObj, nil
}
