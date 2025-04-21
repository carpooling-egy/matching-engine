package valhalla

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/model"
)

type Mapper struct {
	RouteMapper     re.OperationMapper[model.RouteParams, model.Route, *pb.Api, *pb.Api]
	WalkMapper      re.OperationMapper[model.WalkParams, model.Distance, *pb.Api, *pb.Api]
	IsochroneMapper re.OperationMapper[model.IsochroneParams, *model.Isochrone, *pb.Api, *pb.Api]
	MatrixMapper    re.OperationMapper[model.DistanceTimeMatrixParams, *model.DistanceTimeMatrix, *pb.Api, *pb.Api]
}

func NewMapper() *Mapper {
	return &Mapper{
		//Route:     NewRouteMapper(),
		//Walk:      NewWalkMapper(),
		//Isochrone: NewIsochroneMapper(),
		//Matrix:    NewDistanceTimeMatrixMapper(),
	}
}
