package generator

import (
	"iter"
	"matching-engine/internal/model"
)

// PathGenerator defines the interface for path generation strategies
type PathGenerator interface {
	// GeneratePaths returns an iterator that produces all possible paths
	GeneratePaths(
		path []model.PathPoint,
		pickup, dropoff *model.PathPoint,
	) (iter.Seq2[[]model.PathPoint, error], error)
}
