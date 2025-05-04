package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
	"matching-engine/internal/model"
	"math"
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

	points := make([]*pb.Location, len(params.Points()))
	for i, point := range params.Points() {
		points[i] = common.CreateLocation(point.Lat(), point.Lng())
	}

	var (
		costingType pb.Costing_Type
		costing     *pb.Costing
	)

	if params.Profile() == model.Pedestrian {
		costingType = pb.Costing_pedestrian
		costing = common.DefaultPedestrianCosting
	} else { // assume auto
		costingType = pb.Costing_auto_
		costing = common.DefaultAutoCosting
	}

	return &pb.Api{
		Options: &pb.Options{
			Action:      pb.Options_sources_to_targets,
			Sources:     points,
			Targets:     points,
			Units:       common.DefaultUnit,
			Format:      common.DefaultFormat,
			CostingType: costingType,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_auto_): costing,
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

	// Assumes square matrix (sources == targets)
	length := int(math.Sqrt(float64(len(flattenedDistanceMatrix))))

	distanceMatrix := make([][]model.Distance, length)
	timeMatrix := make([][]time.Duration, length)

	distanceUnit, err := common.ToDomainDistanceUnit(common.DefaultUnit)
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
