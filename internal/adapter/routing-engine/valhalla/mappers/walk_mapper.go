package mappers

import (
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/model"
)

type WalkMapper struct{}

var _ re.OperationMapper[
	model.WalkParams,
	model.Distance,
	*pb.Api,
	*pb.Api] = WalkMapper{}

func (WalkMapper) ToTransport(params model.WalkParams) *pb.Api {
	// Convert model.WalkParams to pb.Api
	// This is a placeholder implementation.
	return &pb.Api{}
}

func (WalkMapper) FromTransport(api *pb.Api) model.Distance {
	// Convert pb.Api to model.Distance
	// This is a placeholder implementation.
	return model.Distance{}
}
