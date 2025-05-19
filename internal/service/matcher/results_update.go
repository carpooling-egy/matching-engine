package matcher

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
)

// validateOfferNode checks if the given offer node is valid and provides appropriate errors if not.
func validateOfferNode(offerNode *model.OfferNode) error {
	if offerNode == nil {
		return fmt.Errorf(errors.ErrNilOfferNode)
	}
	return offerNode.ValidateOffer()
}

// createMatchingResult generates a matching result for the given offer node.
func createMatchingResult(offerNode *model.OfferNode) (*model.MatchingResult, error) {
	if err := validateOfferNode(offerNode); err != nil {
		return nil, err
	}
	matchingResult, err := model.NewMatchingResultFromOfferNode(offerNode)
	if err != nil {
		return nil, fmt.Errorf("failed to create matching result: %w", err)
	}
	return matchingResult, nil
}

func (matcher *Matcher) updateResults(offerNode *model.OfferNode) {
	matchingResult, err := createMatchingResult(offerNode)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create matching result for offer %s", offerNode.Offer().ID())
		return // continue
	}
	matcher.results = append(matcher.results, *matchingResult)
}
