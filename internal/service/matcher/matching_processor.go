package matcher

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
)

// processMaximumMatching finds maximum matches and updates results.
func (matcher *Matcher) processMaximumMatching(graph *model.MaximumMatchingGraph, limit int) error {
	maxPairs, err := matcher.maximumMatching.FindMaximumMatching(graph)
	if err != nil {
		return fmt.Errorf("failed to find maximum matching: %w", err)
	}

	if len(maxPairs) == 0 {
		log.Info().Msg("No maximum matching found")
		return nil
	}

	for _, pair := range maxPairs {
		offerNode := pair.First
		edge := pair.Second
		requestNode := edge.RequestNode()
		offerNode.SetMatched(true)
		offerNode.AddNewlyMatchedRequest(requestNode.Request())
		newPath := edge.NewPath()
		if newPath == nil {
			return fmt.Errorf("edge with nil path encountered for offer %s and request %s", offerNode.Offer().ID(), requestNode.Request().ID())
		}
		offerNode.Offer().SetPath(newPath)

		requestSet, exists := matcher.potentialOfferRequests.Get(offerNode.Offer().ID())
		if exists {
			requestSet.Remove(requestNode.Request().ID())
		}

		if len(offerNode.GetAllRequests()) >= limit {
			log.Info().Msgf("Exceeded maximum matching limit of %d requests", limit)
			matcher.updateResults(offerNode)
			matcher.availableOffers.Delete(offerNode.Offer().ID())
			matcher.potentialOfferRequests.Delete(offerNode.Offer().ID())
		}

		matcher.availableRequests.Delete(requestNode.Request().ID())
	}
	return nil
}
