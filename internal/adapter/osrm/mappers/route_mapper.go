package mappers

import (
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"strings"
	"time"
)

type RouteMapper struct{}

var _ routing.OperationMapper[
	*model.RouteParams,
	*model.Route,
	model.OSRMTransport,
	map[string]any,
] = RouteMapper{}

func (RouteMapper) ToTransport(params *model.RouteParams) (model.OSRMTransport, error) {
	if params == nil {
		return model.OSRMTransport{}, fmt.Errorf("params cannot be nil")
	}
	coords := params.Waypoints()
	coordStrs := make([]string, len(coords))
	for i, c := range coords {
		coordStrs[i] = fmt.Sprintf("%.6f,%.6f", c.Lng(), c.Lat())
	}
	return model.OSRMTransport{
		PathVars: map[string]string{
			"coordinates": strings.Join(coordStrs, ";"),
		},
		QueryParams: map[string]any{
			"overview":   "full",
			"geometries": "polyline6",
		},
	}, nil
}

func (RouteMapper) FromTransport(response map[string]any) (*model.Route, error) {
	if response == nil || len(response) == 0 {
		return nil, fmt.Errorf("empty OSRM response")
	}

	routes, ok := response["routes"].([]any)
	if !ok || len(routes) == 0 {
		return nil, fmt.Errorf("no routes in OSRM response")
	}

	routeObj, ok := routes[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid route object in OSRM response")
	}

	geometry, _ := routeObj["geometry"].(string)
	distance, _ := routeObj["distance"].(float64)
	duration, _ := routeObj["duration"].(float64)

	polylineObj, err := model.NewPolyline(geometry)
	if err != nil {
		return nil, fmt.Errorf("invalid polyline in OSRM response: %w", err)
	}

	distanceObj, err := model.NewDistance(float32(distance), model.DistanceUnitMeter)
	if err != nil {
		return nil, fmt.Errorf("invalid distance in OSRM response: %w", err)
	}

	route, err := model.NewRoute(
		polylineObj, distanceObj,
		time.Duration(duration*float64(time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	return route, nil
}
