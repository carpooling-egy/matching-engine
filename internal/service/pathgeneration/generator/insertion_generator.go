package generator

import (
	"fmt"
	"iter"
	"matching-engine/internal/model"
)

type InsertionPathGenerator struct {
}

func NewInsertionPathGenerator() PathGenerator {
	return &InsertionPathGenerator{}
}

// GeneratePaths returns an iterator that generates all possible paths by inserting
// pickup and dropoff points into the existing path.
func (ip *InsertionPathGenerator) GeneratePaths(
	path []model.PathPoint,
	pickup, dropoff *model.PathPoint,
) (iter.Seq2[[]model.PathPoint, error], error) {
	pathLength := len(path)
	if pathLength < 2 {
		return nil, fmt.Errorf("path must contain at least two points")
	}

	// The following upper bound guarantee that we are not losing any possible path though generating
	// paths is still expensive, especially that there are some combinations that are infeasible with high probability
	// Example:
	// if path times is [1, 5, 10, 15, 20]
	// and we want to insert a pickup at 8 and dropoff at 12
	// trying to insert the new points as follows [1, p, d, 5, 10, 15, 20]
	// might be theoretically feasible but with small probability in practical cases

	// We can also add an upper bound to the pickup position pick up can't be added after a point whose expected
	// arrival time is greater than the latest pickup time (latest dropoff time - duration from pickup to

	// Calculate upper bounds for insertion positions
	upperIndex := pathLength - 1
	for i := pathLength - 2; i > 0; i-- {
		if path[i].ExpectedArrivalTime().After(dropoff.ExpectedArrivalTime()) {
			upperIndex = i
		} else {
			break
		}
	}

	// Buffer for building each new path (length = original + 2)
	newLen := pathLength + 2

	return func(yield func([]model.PathPoint, error) bool) {

		// Generate all valid paths
		for pickupPos := 1; pickupPos <= upperIndex; pickupPos++ {
			for dropoffPos := pickupPos; dropoffPos <= upperIndex; dropoffPos++ {

				//// Generate current path
				prefix := path[:pickupPos]
				middle := path[pickupPos:dropoffPos]
				suffix := path[dropoffPos:]

				newPath := make([]model.PathPoint, 0, newLen)
				newPath = append(newPath, prefix...)
				newPath = append(newPath, *pickup)
				newPath = append(newPath, middle...)
				newPath = append(newPath, *dropoff)
				newPath = append(newPath, suffix...)

				if !yield(newPath, nil) {
					// Stop iterating if consumer is done
					return
				}
			}
		}
	}, nil
}
