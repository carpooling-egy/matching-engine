package matcher

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
)

// buildMatchingGraph constructs the graph by finding feasible paths and connecting offers with requests.
func (matcher *Matcher) buildMatchingGraph(graph *model.Graph) (bool, error) {
	hasNewEdge := false
	err := matcher.potentialOfferRequests.Range(func(offerID string, requestSet *collections.Set[string]) error {
		offerNode, exists := matcher.availableOffers.Get(offerID)
		if !exists || offerNode == nil {
			matcher.potentialOfferRequests.Delete(offerID)
			return nil
		}

		var requestNodes []*model.RequestNode
		requestSetSlice := requestSet.ToSlice()
		for _, requestID := range requestSetSlice {
			if requestNode, ok := matcher.availableRequests.Get(requestID); ok && requestNode != nil {
				requestNodes = append(requestNodes, requestNode)
			} else {
				requestSet.Remove(requestID)
			}
		}

		err := matcher.timeMatrixCachePopulator.Populate(offerNode, requestNodes, false)
		if err != nil {
			return err
		}

		for _, requestID := range requestSetSlice {
			requestNode, _ := matcher.availableRequests.Get(requestID)

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
			graph.AddOfferNode(offerNode)
			graph.AddRequestNode(requestNode)
			graph.AddEdge(offerNode, requestNode, edge)

		}
		return nil
	})

	return hasNewEdge, err
}
