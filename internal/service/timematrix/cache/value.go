package cache

import (
	"matching-engine/internal/model"
	"time"
)

type PathPointMappedTimeMatrix struct {
	timeMatrix     [][]time.Duration
	pointIdToIndex map[model.PathPointID]int
}

func NewPathPointMappedTimeMatrix(
	timeMatrix [][]time.Duration,
	pointIdToIndex map[model.PathPointID]int,
) *PathPointMappedTimeMatrix {
	return &PathPointMappedTimeMatrix{
		timeMatrix:     timeMatrix,
		pointIdToIndex: pointIdToIndex,
	}
}

func (m *PathPointMappedTimeMatrix) TimeMatrix() [][]time.Duration {
	return m.timeMatrix
}

func (m *PathPointMappedTimeMatrix) PointIdToIndex() map[model.PathPointID]int {
	return m.pointIdToIndex
}
