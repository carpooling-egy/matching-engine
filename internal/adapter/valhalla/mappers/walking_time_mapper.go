package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
	"matching-engine/internal/model"
	"time"
)

type WalkingTimeMapper struct{}

var _ re.OperationMapper[
	*model.WalkParams,
	time.Duration,
	*pb.Api,
	*pb.Api,
] = WalkingTimeMapper{}

func (WalkingTimeMapper) ToTransport(params *model.WalkParams) (*pb.Api, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	wps := []*model.Coordinate{params.Origin(), params.Destination()}
	locations := make([]*pb.Location, 2)
	for i, wp := range wps {
		locations[i] = common.CreateLocation(wp.Lat(), wp.Lng())
	}

	return &pb.Api{
		Options: &pb.Options{
			Action:      pb.Options_route,
			Units:       common.DefaultUnit,
			Format:      common.DefaultResponseFormat,
			CostingType: pb.Costing_auto_,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): common.DefaultPedestrianCosting,
			},
			HasShapeFormat: &pb.Options_ShapeFormat{
				ShapeFormat: common.DefaultShapeFormat,
			},
			Locations: locations,
			PbfFieldSelector: &pb.PbfFieldSelector{
				Directions: true,
			},
		},
	}, nil
}

func (WalkingTimeMapper) FromTransport(response *pb.Api) (time.Duration, error) {
	if response == nil {
		return 0, fmt.Errorf("response cannot be nil")
	}

	var timeInSeconds float64 = 0

	for _, leg := range response.GetDirections().GetRoutes()[0].GetLegs() {
		timeInSeconds += leg.GetSummary().Time
	}

	return time.Duration(timeInSeconds * float64(time.Second)), nil
}
