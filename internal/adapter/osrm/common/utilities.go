package common

import (
	"fmt"
	"time"
)

func GetCumulativeDurations(timeMatrix [][]time.Duration, pathLength int) ([]time.Duration, error) {
	cumulativeDurations := make([]time.Duration, pathLength)
	cumulativeDurations[0] = 0
	for i := 0; i < pathLength-1; i++ {
		duration := timeMatrix[i][i+1]
		if duration < 0 {
			return nil, fmt.Errorf("negative duration found between points %d and %d", i, i+1)
		}
		cumulativeDurations[i+1] = cumulativeDurations[i] + duration
	}
	return cumulativeDurations, nil
}
