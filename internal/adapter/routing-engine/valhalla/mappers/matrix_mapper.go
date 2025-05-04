package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/adapter/routing-engine/valhalla/helpers"
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
	points := make([]*pb.Location, len(params.Points()))
	for i, point := range params.Points() {
		points[i] = helpers.CreateLocation(point.Lat(), point.Lng(), helpers.DefaultLocationType)
	}

	return &pb.Api{
		Options: &pb.Options{
			Action:      pb.Options_sources_to_targets,
			Sources:     points,
			Targets:     points,
			Units:       helpers.DefaultUnit,
			Format:      helpers.DefaultFormat,
			CostingType: pb.Costing_auto_,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): helpers.DefaultAutoCosting,
			},
			DateTimeType: pb.Options_depart_at,
			HasDateTime: &pb.Options_DateTime{
				// TODO check how valhalla handles timezones
				DateTime: params.DepartureTime().Format("2006-01-02T15:04"),
			},
			PbfFieldSelector: &pb.PbfFieldSelector{
				Matrix: true,
			},
		},
	}, nil
}

func (MatrixMapper) FromTransport(api *pb.Api) (*model.DistanceTimeMatrix, error) {
	matrix := api.GetMatrix()
	if matrix == nil {
		return nil, fmt.Errorf("matrix is nil")
	}

	flattenedDistanceMatrix := matrix.GetDistances()
	flattenedTimeMatrix := matrix.GetTimes()

	length := len(api.GetOptions().GetSources()) // Assumes square matrix (sources == targets)

	if len(flattenedDistanceMatrix) != length*length || len(flattenedTimeMatrix) != length*length {
		return nil, fmt.Errorf("flattened matrix size mismatch")
	}

	distanceMatrix := make([][]model.Distance, length)
	timeMatrix := make([][]time.Duration, length)

	distanceUnit, err := helpers.ToDomainDistanceUnit(helpers.DefaultUnit)
	if err != nil {
		return nil, err
	}

	for i := 0; i < length; i++ {
		distanceMatrix[i] = make([]model.Distance, length)
		timeMatrix[i] = make([]time.Duration, length)

		for j := 0; j < length; j++ {
			index := i*length + j

			distance, err := model.NewDistance(float32(flattenedDistanceMatrix[index]), distanceUnit)
			if err != nil {
				return nil, err
			}
			distanceMatrix[i][j] = *distance
			timeMatrix[i][j] = time.Duration(flattenedTimeMatrix[index]) * time.Second
		}
	}

	return model.NewDistanceTimeMatrix(distanceMatrix, timeMatrix)
}
