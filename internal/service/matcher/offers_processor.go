package matcher

import "matching-engine/internal/model"

// processUnmatchedOffers processes offers that are not in the graph and updates results
func (matcher *Matcher) processUnmatchedOffers(graph *model.Graph) {
	potentialOffers := graph.OfferNodes()
	matcher.availableOffers.ForEach(func(offerID string, offerNode *model.OfferNode) error {
		if potentialOffers.Contains(offerNode.Offer().ID()) {
			return nil // continue
		}

		if offerNode.IsMatched() {
			matcher.updateResults(offerNode)
		}

		matcher.potentialOfferRequests.Delete(offerID)
		return nil // continue
	})
	matcher.availableOffers = potentialOffers
}

// processRemainingOffers appends leftover matched offers to result.
func (matcher *Matcher) processRemainingOffers() error {
	return matcher.availableOffers.Range(func(offerID string, offerNode *model.OfferNode) error {
		if offerNode.IsMatched() {
			matcher.updateResults(offerNode)
		}
		return nil
	})
}
