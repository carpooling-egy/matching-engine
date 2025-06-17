package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
	"matching-engine/internal/model"
	"time"
)

type MatrixMapper struct{}

var _ re.OperationMapper[
	*model.DistanceTimeMatrixParams,
	*model.DistanceTimeMatrix,
	*pb.Api,
	*pb.Api,
] = MatrixMapper{}

func (MatrixMapper) ToTransport(params *model.DistanceTimeMatrixParams) (*pb.Api, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	sources := common.WayPointsToLocations(params.Sources())
	targets := common.WayPointsToLocations(params.Targets())

	costingType, costing, err := common.MapProfileToCosting(params.Profile())
	if err != nil {
		return nil, err
	}

	return &pb.Api{
		Options: &pb.Options{
			Action:      pb.Options_sources_to_targets,
			Sources:     sources,
			Targets:     targets,
			Units:       common.DefaultUnit,
			Format:      common.DefaultResponseFormat,
			CostingType: costingType,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): costing,
			},
			DateTimeType: pb.Options_depart_at,
			HasDateTime: &pb.Options_DateTime{
				// TODO check how valhalla handles timezones
				DateTime: params.DepartureTime().Format(common.DefaultTimeFormat),
			},
			PbfFieldSelector: &pb.PbfFieldSelector{
				Matrix: true,
			},
		},
	}, nil
}

func (MatrixMapper) FromTransport(response *pb.Api) (*model.DistanceTimeMatrix, error) {
	if response == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

	matrix := response.GetMatrix()
	if matrix == nil {
		return nil, fmt.Errorf("matrix is nil")
	}

	flattenedDistanceMatrix := matrix.GetDistances()
	flattenedTimeMatrix := matrix.GetTimes()

	fromIndexes := matrix.FromIndices
	toIndexes := matrix.ToIndices

	if len(fromIndexes) == 0 || len(toIndexes) == 0 {
		return nil, fmt.Errorf("from or to indices are empty")
	}

	numOfRows := int(matrix.FromIndices[len(fromIndexes)-1] + 1)
	numOfCols := int(matrix.ToIndices[len(toIndexes)-1] + 1)

	distanceMatrix := make([][]model.Distance, numOfRows)
	timeMatrix := make([][]time.Duration, numOfRows)

	distanceUnit, err := common.ToDomainDistanceUnit(common.DefaultUnit)
	if err != nil {
		return nil, err
	}

	for i := 0; i < numOfRows; i++ {
		distanceMatrix[i] = make([]model.Distance, numOfCols)
		timeMatrix[i] = make([]time.Duration, numOfCols)

		for j := 0; j < numOfCols; j++ {
			index := i*numOfCols + j

			distance, err := model.NewDistance(float32(flattenedDistanceMatrix[index]), distanceUnit)
			if err != nil {
				return nil, err
			}
			distanceMatrix[i][j] = *distance
			timeMatrix[i][j] = time.Duration(float64(flattenedTimeMatrix[index]) * float64(time.Minute))
		}
	}

	return model.NewDistanceTimeMatrix(distanceMatrix, timeMatrix)
}
