package model

import (
	"errors"
	"time"
)

type DistanceTimeMatrix struct {
	distances [][]Distance
	times     [][]time.Duration
}

func NewDistanceTimeMatrix(distances [][]Distance, times [][]time.Duration) (
	*DistanceTimeMatrix, error,
) {
	if len(distances) == 0 || len(times) == 0 {
		return nil, errors.New("distances or times matrix is empty")
	}
	if len(distances) != len(times) {
		return nil, errors.New("distances and times matrix must have the same shape")
	}
	for i := range distances {
		if len(distances[i]) != len(times[i]) {
			return nil, errors.New("distances and times matrix must have the same shape")
		}
		for j := range distances[i] {
			if distances[i][j].value < 0 {
				return nil, errors.New("distance cannot be negative")
			}
			if distances[i][j].unit != Meter && distances[i][j].unit != Kilometer &&
				distances[i][j].unit != Mile {
				return nil, errors.New("invalid distance unit")
			}
			if times[i][j] < 0 {
				return nil, errors.New("time duration cannot be negative")
			}
		}
	}

	return &DistanceTimeMatrix{
		distances: distances,
		times:     times,
	}, nil
}

func (dtm *DistanceTimeMatrix) Distances() [][]Distance {
	return dtm.distances
}

func (dtm *DistanceTimeMatrix) Times() [][]time.Duration {
	return dtm.times
}
