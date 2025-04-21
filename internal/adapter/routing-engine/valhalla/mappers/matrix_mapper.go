package mappers

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/model"
)

type MatrixMapper struct{}

var _ re.OperationMapper[
	model.DistanceTimeMatrixParams,
	*model.DistanceTimeMatrix,
	*pb.Api,
	*pb.Api,
] = MatrixMapper{}

func (MatrixMapper) ToTransport(params model.DistanceTimeMatrixParams) *pb.Api {
	// Convert model.IsochroneParams to pb.Api
	// This is a placeholder implementation.
	return &pb.Api{}
}

func (MatrixMapper) FromTransport(api *pb.Api) *model.DistanceTimeMatrix {
	// Convert pb.Api to model.Isochrone
	// This is a placeholder implementation.
	return &model.DistanceTimeMatrix{}
}
