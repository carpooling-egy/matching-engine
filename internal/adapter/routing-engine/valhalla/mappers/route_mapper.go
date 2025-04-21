package mappers

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/model"
)

type RouteMapper struct{}

var _ re.OperationMapper[
	model.RouteParams,
	model.Route,
	*pb.Api,
	*pb.Api,
] = RouteMapper{}

func (RouteMapper) ToTransport(params model.RouteParams) *pb.Api {

	wps := params.Waypoints()
	locations := make([]*pb.Location, len(wps))
	for i, wp := range wps {
		locations[i] = valhalla.CreateLocation(wp.Lat(), wp.Lng(), valhalla.DefaultLocationType)
	}

	return &pb.Api{
		Options: &pb.Options{
			Action: pb.Options_route,
			Units:  valhalla.DefaultUnit,
			Format: valhalla.DefaultFormat,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): valhalla.DefaultAutoCosting,
			},
			HasShapeFormat: &pb.Options_ShapeFormat{
				ShapeFormat: valhalla.DefaultShapeFormat,
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
	}
}

func (RouteMapper) FromTransport(api *pb.Api) model.Route {
	polyline := api.GetDirections().GetRoutes()[0].GetLegs()[0].GetShape()

	var distance float32 = 0
	var time float64 = 0
	for _, leg := range api.GetDirections().GetRoutes()[0].GetLegs() {
		distance += leg.GetSummary().Length
		time += leg.GetSummary().Time
	}

	polylineObj, err := model.NewPolyline(polyline)
	if err != nil {
		// Handle error, maybe log it or return a default value
		// You could return an empty model.Route or handle accordingly
		return model.Route{}
	}

	distanceObj, err := model.NewDistance(float64(distance), valhalla.DefaultUnit)
	if err != nil {
		// Handle error, maybe log it or return a default value
		// You could return an empty model.Route or handle accordingly
		return model.Route{}
	}

	route, err := model.NewRoute(polylineObj, distanceObj, time)
	if err != nil {
		// Handle error if route creation fails
		// You could return an empty model.Route or handle accordingly
		return model.Route{}
	}

	return *route
}
