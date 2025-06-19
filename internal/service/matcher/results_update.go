package matcher

import (
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
)

func (matcher *Matcher) updateResults(offerNode *model.OfferNode) {
	matchingResult, err := model.NewMatchingResultFromOfferNode(offerNode)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create matching result for offer %s", offerNode.Offer().ID())
		return // continue
	}
	matcher.results = append(matcher.results, matchingResult)
}
