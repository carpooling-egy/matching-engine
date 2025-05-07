package valhalla

import (
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/mappers"
	"matching-engine/internal/model"
	"time"
)

type Mapper struct {
	RouteMapper       re.OperationMapper[*model.RouteParams, *model.Route, *pb.Api, *pb.Api]
	DrivingTimeMapper re.OperationMapper[*model.RouteParams, []time.Duration, *pb.Api, *pb.Api]
	WalkingTimeMapper re.OperationMapper[*model.WalkParams, time.Duration, *pb.Api, *pb.Api]
	IsochroneMapper   re.OperationMapper[*model.IsochroneParams, *model.Isochrone, *pb.Api, *pb.Api]
	MatrixMapper      re.OperationMapper[*model.DistanceTimeMatrixParams, *model.DistanceTimeMatrix, *pb.Api, *pb.Api]
	SnapToRoadMapper  re.OperationMapper[*model.Coordinate, *model.Coordinate, *pb.Api, *pb.Api]
}

func NewMapper() *Mapper {
	return &Mapper{
		RouteMapper:       mappers.RouteMapper{},
		DrivingTimeMapper: mappers.DrivingTimeMapper{},
		WalkingTimeMapper: mappers.WalkingTimeMapper{},
		IsochroneMapper:   mappers.IsochroneMapper{},
		MatrixMapper:      mappers.MatrixMapper{},
		SnapToRoadMapper:  mappers.SnapToRoadMapper{},
	}
}
