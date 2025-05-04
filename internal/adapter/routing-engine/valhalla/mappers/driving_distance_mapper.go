package mappers

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/adapter/routing-engine/valhalla/helpers"
	"matching-engine/internal/model"
)

type DrivingDistanceMapper struct{}

var _ re.OperationMapper[
	*model.RouteParams,
	*model.Distance,
	*pb.Api,
	*pb.Api,
] = DrivingDistanceMapper{}

func (d DrivingDistanceMapper) ToTransport(params *model.RouteParams) (*pb.Api, error) {
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
				DateTime: params.DepartureTime().String(), // TODO check how valhalla handles timezones
			},
			PbfFieldSelector: &pb.PbfFieldSelector{
				Directions: true,
			},
		},
	}, nil
}

func (d DrivingDistanceMapper) FromTransport(api *pb.Api) (*model.Distance, error) {
	var distance float32 = 0
	for _, leg := range api.GetDirections().GetRoutes()[0].GetLegs() {
		distance += leg.GetSummary().Length
	}

	unit, err := helpers.ToDomainDistanceUnit(helpers.DefaultUnit)
	if err != nil {
		return nil, err
	}

	distanceObj, err := model.NewDistance(distance, unit)
	if err != nil {
		return nil, err
	}

	return distanceObj, nil
}
