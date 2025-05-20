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

			path, valid, err := matcher.matchEvaluator.Evaluate(offerNode, requestNode)

			if err != nil {
				return fmt.Errorf("error evaluating the match %v", err)
			}

			if !valid {
				requestSet.Remove(requestID)
				continue
			}

			hasNewEdge = true
			edge := model.NewEdge(requestNode, path)
			offerNode.AddEdge(edge)
			graph.AddOfferNode(offerNode)
			potentialRequests.Set(requestID, requestNode)

		}
		return nil
	})

	return hasNewEdge, err
}
