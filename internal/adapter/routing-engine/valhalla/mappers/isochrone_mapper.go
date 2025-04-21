package mappers

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/model"
)

type IsochroneMapper struct{}

var _ re.OperationMapper[
	model.IsochroneParams,
	*model.Isochrone,
	*pb.Api,
	*pb.Api,
] = IsochroneMapper{}

func (IsochroneMapper) ToTransport(params model.IsochroneParams) *pb.Api {
	// Convert model.IsochroneParams to pb.Api
	// This is a placeholder implementation.
	return &pb.Api{}
}

func (IsochroneMapper) FromTransport(api *pb.Api) *model.Isochrone {
	// Convert pb.Api to model.Isochrone
	// This is a placeholder implementation.
	return &model.Isochrone{}
}
