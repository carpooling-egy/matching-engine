package matcher

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
)

// processMaximumMatching finds maximum matches and updates results.
func (matcher *Matcher) processMaximumMatching(graph *model.Graph, limit int) error {
	maxEdges, err := matcher.maximumMatching.FindMaximumMatching(graph)
	if err != nil {
		return fmt.Errorf("failed to find maximum matching: %w", err)
	}

	return maxEdges.Range(func(offerNode *model.OfferNode, edge *model.Edge) error {
		requestNode := edge.RequestNode()
		if requestNode == nil {
			log.Warn().Msg("Nil request node in maximum matching")
			return nil
		}

		offerNode.SetMatched(true)
		offerNode.AddNewlyMatchedRequest(requestNode.Request())

		if len(offerNode.GetAllRequests()) >= limit {
			matcher.updateResults(offerNode)
			matcher.availableOffers.Delete(offerNode.Offer().ID())
			matcher.potentialOfferRequests.Delete(offerNode.Offer().ID())
		}

		matcher.availableRequests.Delete(requestNode.Request().ID())
		return nil
	})
}
