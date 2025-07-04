package mappers

import (
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"strings"
	"time"
)

type MatrixMapper struct{}

var _ routing.OperationMapper[
	*model.DistanceTimeMatrixParams,
	*model.DistanceTimeMatrix,
	model.OSRMTransport,
	map[string]any,
] = MatrixMapper{}

func (MatrixMapper) ToTransport(params *model.DistanceTimeMatrixParams) (model.OSRMTransport, error) {
	if params == nil {
		return model.OSRMTransport{}, fmt.Errorf("params cannot be nil")
	}

	sources := params.Sources()
	targets := params.Targets()

	allCoords := append(sources, targets...)

	coordStrs := make([]string, len(allCoords))
	for i, c := range allCoords {
		coordStrs[i] = fmt.Sprintf("%.6f,%.6f", c.Lng(), c.Lat())
	}

	sourceIndices := make([]string, len(sources))
	for i := range sources {
		sourceIndices[i] = fmt.Sprintf("%d", i)
	}

	targetIndices := make([]string, len(targets))
	for i := range targets {
		targetIndices[i] = fmt.Sprintf("%d", len(sources)+i)
	}

	return model.OSRMTransport{
		PathVars: map[string]string{
			"coordinates": strings.Join(coordStrs, ";"),
		},
		QueryParams: map[string]any{
			"sources":             strings.Join(sourceIndices, ";"),
			"destinations":        strings.Join(targetIndices, ";"),
			"annotations":         "duration,distance",
			"fallback_speed":      5,
			"fallback_coordinate": "snapped",
		},
	}, nil
}

func (MatrixMapper) FromTransport(response map[string]any) (*model.DistanceTimeMatrix, error) {
	if (response == nil) || (len(response) == 0) {
		return nil, fmt.Errorf("empty OSRM response")
	}

	rawDurations, ok := response["durations"].([]any)
	if !ok {
		return nil, fmt.Errorf("no durations in OSRM response")
	}

	rawDistances, ok := response["distances"].([]any)
	if !ok {
		return nil, fmt.Errorf("no distances in OSRM response")
	}

	if len(rawDurations) != len(rawDistances) {
		return nil, fmt.Errorf("mismatched matrix sizes: %d durations rows vs %d distances rows",
			len(rawDurations),
			len(rawDistances),
		)
	}

	n := len(rawDurations)
	distanceMatrix := make([][]model.Distance, n)
	timeMatrix := make([][]time.Duration, n)

	for i := 0; i < n; i++ {
		durRow, ok := rawDurations[i].([]any)
		if !ok {
			return nil, fmt.Errorf("invalid row %d in durations", i)
		}

		distRow, ok := rawDistances[i].([]any)
		if !ok {
			return nil, fmt.Errorf("invalid row %d in distances", i)
		}

		if len(durRow) != len(distRow) {
			return nil, fmt.Errorf("mismatched row lengths at row %d: %d durations vs %d distances", i,
				len(durRow),
				len(distRow),
			)
		}

		distanceMatrix[i] = make([]model.Distance, len(durRow))
		timeMatrix[i] = make([]time.Duration, len(durRow))

		for j := 0; j < len(durRow); j++ {
			durationInSec, ok := durRow[j].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid duration value at [%d][%d]", i, j)
			}

			timeMatrix[i][j] = time.Duration(durationInSec * float64(time.Second))

			distanceInMeters, ok := distRow[j].(float64)
			if !ok {
				return nil, fmt.Errorf("invalid distance value at [%d][%d]", i, j)
			}

			distanceInMeters = max(distanceInMeters, 0)
			distObj, err := model.NewDistance(float32(distanceInMeters), model.DistanceUnitMeter)
			if err != nil {
				return nil, fmt.Errorf("failed to create Distance at [%d][%d]: %w", i, j, err)
			}

			distanceMatrix[i][j] = *distObj
		}
	}

	return model.NewDistanceTimeMatrix(distanceMatrix, timeMatrix)
}
