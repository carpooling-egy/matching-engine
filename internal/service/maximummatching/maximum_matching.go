package maximummatching

import (
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
)

type MaximumMatching interface {
	// FindMaximumMatching finds the maximum matching in a bipartite graph.
	FindMaximumMatching(graph *model.Graph) (collections.SyncMap[*model.OfferNode, *model.Edge], error)
}
