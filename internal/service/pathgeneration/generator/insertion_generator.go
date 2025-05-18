package generator

import (
	"fmt"
	"iter"
	"matching-engine/internal/model"
)

type InsertionPathGenerator struct {
}

func NewInsertionPathGenerator() *InsertionPathGenerator {
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

	return func(yield func([]model.PathPoint, error) bool) {

		// The following upper bounds guarantee that we are not losing any possible path though generating
		// paths is still expensive, especially that there are some combinations that are infeasible with high probability
		// Example:
		// if path times is [1, 5, 10, 15, 20]
		// and we want to insert a pickup at 8 and dropoff at 12
		// trying to insert the new points as follows [1, p, d, 5, 10, 15, 20]
		// might be theoretically feasible but with small probability in practical cases

		// Calculate upper bounds for insertion positions
		upperIndexForDropoff := pathLength - 1
		upperIndexForPickup := pathLength - 1
		for i := pathLength - 2; i > 0; i-- {
			if path[i].ExpectedArrivalTime().After(dropoff.ExpectedArrivalTime()) {
				upperIndexForDropoff = i - 1
			}
			if path[i].ExpectedArrivalTime().After(pickup.ExpectedArrivalTime()) {
				upperIndexForPickup = i - 1
				break
			}
		}

		// Buffer for building each new path (length = original + 2)
		newLen := pathLength + 2
		buf := make([]model.PathPoint, newLen)

		// Generate all valid paths
		for pickupPos := 1; pickupPos <= upperIndexForPickup; pickupPos++ {
			for dropoffPos := pickupPos; dropoffPos <= upperIndexForDropoff; dropoffPos++ {
				// Generate current path
				// 1) prefix [0:pickupPos)
				copy(buf[0:pickupPos], path[0:pickupPos])

				// 2) insert pickup
				buf[pickupPos] = *pickup

				// 3) segment [pickupPos:dropoffPos)
				copy(buf[pickupPos+1:pickupPos+1+(dropoffPos-pickupPos)], path[pickupPos:dropoffPos])

				// 4) insert dropoff
				buf[pickupPos+1+(dropoffPos-pickupPos)] = *dropoff

				// 5) suffix [dropoffPos:pathLength)
				copy(buf[pickupPos+2+(dropoffPos-pickupPos):], path[dropoffPos:])

				// Clone buffer into a fresh slice before yielding
				newPath := make([]model.PathPoint, newLen)
				copy(newPath, buf)

				if !yield(newPath, nil) {
					// Stop iterating if consumer is done
					return
				}
			}
		}
	}, nil
}
