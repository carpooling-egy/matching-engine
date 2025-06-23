package cache

import (
	"fmt"
	"matching-engine/internal/model"
	"time"
)

type PathPointMappedTimeMatrix struct {
	timeMatrix     [][]time.Duration
	pointIdToIndex map[model.PathPointID]int
	pointIdToPoint map[model.PathPointID]model.PathPoint
}

func NewPathPointMappedTimeMatrix(
	timeMatrix [][]time.Duration,
	pointIdToIndex map[model.PathPointID]int,
	pointIdToPoint map[model.PathPointID]model.PathPoint,
) *PathPointMappedTimeMatrix {
	return &PathPointMappedTimeMatrix{
		timeMatrix:     timeMatrix,
		pointIdToIndex: pointIdToIndex,
		pointIdToPoint: pointIdToPoint,
	}
}

func (m *PathPointMappedTimeMatrix) TimeMatrix() [][]time.Duration {
	return m.timeMatrix
}

func (m *PathPointMappedTimeMatrix) PointIdToIndex() map[model.PathPointID]int {
	return m.pointIdToIndex
}

func (m *PathPointMappedTimeMatrix) PrintMatrix() {
	fmt.Println("pointIdToIndex map:")
	for k, v := range m.pointIdToIndex {
		fmt.Printf("  %v: %d\n", k, v)
	}
	fmt.Println("pointIdToPoint map:")
	for k, v := range m.pointIdToPoint {
		fmt.Printf("  %v: %v\n", k, v.Coordinate())
	}
	fmt.Println("timeMatrix:")
	for i, row := range m.timeMatrix {
		fmt.Printf("Row %d: ", i)
		for _, val := range row {
			fmt.Printf("%v ", val)
		}
		fmt.Println()
	}
}
