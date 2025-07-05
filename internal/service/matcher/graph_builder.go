package matcher

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
	"time"
	"matching-engine/cmd/appmetrics"
	"github.com/rs/zerolog/log"
)

// buildMatchingGraph constructs the graph by finding feasible paths and connecting offers with requests.
func (matcher *Matcher) buildMatchingGraph(graph *model.MaximumMatchingGraph) (bool, error) {
	hasNewEdge := false
	err := matcher.potentialOfferRequests.Range(func(offerID string, requestSet *collections.Set[string]) error {
		offerNode, exists := matcher.availableOffers.Get(offerID)
		if !exists || offerNode == nil {
			matcher.potentialOfferRequests.Delete(offerID)
			return nil
		}

		requestNodes := make([]*model.RequestNode, 0, requestSet.Size())
		requestSetSlice := requestSet.ToSlice()
		for _, requestID := range requestSetSlice {
			if requestNode, ok := matcher.availableRequests.Get(requestID); ok && requestNode != nil {
				requestNodes = append(requestNodes, requestNode)
			} else {
				requestSet.Remove(requestID)
			}
		}

		if len(requestNodes) == 0 {
			return nil
		}

		err := matcher.timeMatrixCachePopulator.Populate(offerNode, requestNodes)

		if err != nil {
			return err
		}

		for _, requestNode := range requestNodes {
			start := time.Now()
			path, valid, err := matcher.matchEvaluator.Evaluate(offerNode, requestNode)
			elapsed := time.Since(start)
			log.Info().Str("duration", elapsed.String())
			appmetrics.TrackTime("one_edge",elapsed)
			appmetrics.IncrementCount("one_edge", 1)
			if err != nil {
				return fmt.Errorf("error evaluating the match %v", err)
			}

			if !valid {
				requestSet.Remove(requestNode.Request().ID())
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
