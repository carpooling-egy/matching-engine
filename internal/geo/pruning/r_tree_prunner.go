package pruning

import (
	"github.com/rs/zerolog/log"
	"matching-engine/internal/collections"
	"sort"
	"time"

	"matching-engine/internal/geo"
	"matching-engine/internal/model"

	"github.com/dhconnelly/rtreego"
)

// RTreePruner implements the RoutePruner interface using R-tree spatial indexing
type RTreePruner struct {
	tree     *rtreego.Rtree
	indexMap map[string]int
}

// Prune filters the route to include only segments within the specified threshold
// from the origin point, implementing the RoutePruner interface
func (p *RTreePruner) Prune(origin *model.Coordinate, threshold time.Duration) (model.LineString, error) {
	log.Debug().Msg("RTreePruner.Prune called")

	// Convert query point to rtreego.Point
	query := rtreego.Point{origin.Lat(), origin.Lng()}

	// Convert time threshold to distance in degrees
	thresholdDistance := max(
		geo.MetersToDegrees(float64(threshold.Seconds())*geo.WalkingSpeedMPS),
		geo.MetersToDegrees(float64(MinThresholdDistanceInMeters)),
	)

	thresholdDistanceSquared := thresholdDistance * thresholdDistance

	// Create bounding box (square) around the circle for initial filtering
	searchRect, err := rtreego.NewRect(
		rtreego.Point{query[0] - thresholdDistance, query[1] - thresholdDistance},
		[]float64{2 * thresholdDistance, 2 * thresholdDistance},
	)
	if err != nil {
		return model.LineString{}, err
	}

	// Get candidate segments using R-tree spatial index
	results := p.tree.SearchIntersect(searchRect)

	seen := make(map[string]bool)
	var collected []collections.Tuple2[model.Coordinate, int]

	for _, result := range results {
		segment, ok := result.(*Segment)
		if !ok {
			continue // Skip if not a Segment
		}

		for _, point := range []model.Coordinate{segment.A, segment.B} {
			ptKey := point.Key()
			if seen[ptKey] {
				continue // Skip if already seen
			}

			if isWithinThreshold(origin, &point, thresholdDistanceSquared) {
				if idx, exists := p.indexMap[ptKey]; exists {
					seen[ptKey] = true
					collected = append(collected, collections.NewTuple2(point, idx))
				}
			}
		}
	}

	sort.Slice(collected, func(i, j int) bool {
		return collected[i].Second < collected[j].Second
	})

	subPath := make(model.LineString, 0, len(collected))
	for _, tuple := range collected {
		subPath = append(subPath, tuple.First)
	}

	return subPath, nil
}

// isWithinThreshold checks if a coordinate is within the squared threshold distance
func isWithinThreshold(origin *model.Coordinate, point *model.Coordinate, thresholdDistanceSquared float64) bool {
	return squaredDistance(origin, point) <= thresholdDistanceSquared
}

// squaredDistance calculates the squared Euclidean distance between two coordinates
func squaredDistance(a *model.Coordinate, b *model.Coordinate) float64 {
	latDiff := a.Lat() - b.Lat()
	lngDiff := a.Lng() - b.Lng()
	return latDiff*latDiff + lngDiff*lngDiff
}
