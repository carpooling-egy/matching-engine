package osrm

import (
	"matching-engine/internal/adapter/osrm/mappers"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
)

type Mapper struct {
	RouteMapper      re.OperationMapper[*model.RouteParams, *model.Route, model.OSRMTransport, map[string]any]
	MatrixMapper     re.OperationMapper[*model.DistanceTimeMatrixParams, *model.DistanceTimeMatrix, model.OSRMTransport, map[string]any]
	SnapToRoadMapper re.OperationMapper[*model.Coordinate, *model.Coordinate, model.OSRMTransport, map[string]any]
}

func NewMapper() *Mapper {
	return &Mapper{
		RouteMapper:      mappers.RouteMapper{},
		MatrixMapper:     mappers.MatrixMapper{},
		SnapToRoadMapper: mappers.SnapToRoadMapper{},
	}
}
