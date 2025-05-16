package matcher

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
)

// buildMatchingGraph constructs the graph by finding feasible paths and connecting offers with requests.
func (matcher *Matcher) buildMatchingGraph(graph *model.Graph, potentialRequests *collections.SyncMap[string, *model.RequestNode]) (bool, error) {
	hasNewEdge := false
	err := matcher.potentialOfferRequests.Range(func(offerID string, requestSet *collections.Set[string]) error {
		offerNode, exists := matcher.availableOffers.Get(offerID)
		if !exists || offerNode == nil {
			matcher.potentialOfferRequests.Delete(offerID)
			return nil
		}

		for _, requestID := range requestSet.ToSlice() {
			requestNode, exists := matcher.availableRequests.Get(requestID)
			if !exists || requestNode == nil {
				requestSet.Remove(requestID)
				continue
			}

			newPath, isFeasible, err := matcher.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
			if err != nil {
				return fmt.Errorf("failed to find feasible path for offer %s and request %s: %w", offerID, requestID, err)
			}

			if isFeasible && newPath != nil {
				hasNewEdge = true
				edge := model.NewEdge(requestNode, newPath)
				offerNode.AddEdge(edge)
				graph.AddOfferNode(offerNode)
				potentialRequests.Set(requestID, requestNode)
			} else {
				requestSet.Remove(requestID)
			}
		}
		return nil
	})

	return hasNewEdge, err
}
