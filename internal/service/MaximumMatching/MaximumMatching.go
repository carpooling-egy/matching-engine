package MaximumMatching

import "matching-engine/internal/model"

type MaximumMatching interface {
	// FindMaximumMatching finds the maximum matching in a bipartite graph.
	findMaximumMatching(graph *model.Graph) (map[*model.OfferNode][]*model.Edge, error)
}
