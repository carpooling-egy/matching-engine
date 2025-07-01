package mappers

import (
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
)

type SnapToRoadMapper struct{}

var _ routing.OperationMapper[
	*model.Coordinate,
	*model.Coordinate,
	model.OSRMTransport,
	map[string]any,
] = SnapToRoadMapper{}

func (SnapToRoadMapper) ToTransport(point *model.Coordinate) (model.OSRMTransport, error) {
	if point == nil {
		return model.OSRMTransport{}, fmt.Errorf("point cannot be nil")
	}
	return model.OSRMTransport{
		PathVars: map[string]string{
			"coordinates": fmt.Sprintf("%.6f,%.6f", point.Lng(), point.Lat()),
		},
		QueryParams: map[string]any{},
	}, nil
}

func (SnapToRoadMapper) FromTransport(response map[string]any) (*model.Coordinate, error) {
	if response == nil || len(response) == 0 {
		return nil, fmt.Errorf("empty OSRM response")
	}

	waypoints, ok := response["waypoints"].([]any)
	if !ok || len(waypoints) == 0 {
		return nil, fmt.Errorf("no waypoints in OSRM response")
	}

	waypoint, ok := waypoints[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid waypoint in OSRM response")
	}

	location, ok := waypoint["location"].([]any)
	if !ok || len(location) != 2 {
		return nil, fmt.Errorf("invalid location in OSRM response")
	}

	lng, _ := location[0].(float64)
	lat, _ := location[1].(float64)
	return model.NewCoordinate(lat, lng)
}
