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

func (m *PathPointMappedTimeMatrix) GetTravelTime(from, to model.PathPointID) (time.Duration, bool) {
	fromIdx, fromOk := m.pointIdToIndex[from]
	toIdx, toOk := m.pointIdToIndex[to]

	if !fromOk || !toOk {
		return 0, false
	}

	if fromIdx >= len(m.timeMatrix) || toIdx >= len(m.timeMatrix[fromIdx]) {
		return 0, false
	}

	return m.timeMatrix[fromIdx][toIdx], true
}

func (m *PathPointMappedTimeMatrix) GetCumulativeTravelTime(pathPointIDs []model.PathPointID) (time.Duration, bool) {

	if len(pathPointIDs) == 0 {
		return 0, false
	}

	totalTime := time.Duration(0)
	for i := 0; i < len(pathPointIDs)-1; i++ {
		duration, ok := m.GetTravelTime(pathPointIDs[i], pathPointIDs[i+1])
		if !ok {
			return 0, false
		}
		totalTime += duration
	}

	return totalTime, true
}
