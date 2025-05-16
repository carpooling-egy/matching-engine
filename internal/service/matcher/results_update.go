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
	offer := offerNode.Offer()
	if offer == nil {
		return fmt.Errorf(errors.ErrNilOffer)
	}
	if offer.UserID() == "" {
		return fmt.Errorf(errors.ErrEmptyUserID)
	}
	if offer.ID() == "" {
		return fmt.Errorf(errors.ErrEmptyOfferID)
	}
	if offer.Path() == nil {
		return fmt.Errorf(errors.ErrNilPath)
	}
	if len(offer.Path()) == 0 {
		return fmt.Errorf(errors.ErrEmptyPath)
	}
	return nil
}

// createMatchingResult generates a matching result for the given offer node.
func createMatchingResult(offerNode *model.OfferNode) (*model.MatchingResult, error) {
	if err := validateOfferNode(offerNode); err != nil {
		return nil, err
	}

	newlyMatchedRequests := offerNode.NewlyAssignedMatchedRequests()
	if newlyMatchedRequests == nil {
		return nil, fmt.Errorf(errors.ErrNilMatchedRequests)
	}
	if len(newlyMatchedRequests) == 0 {
		return nil, fmt.Errorf(errors.ErrEmptyMatchedRequests)
	}

	allRequestsCount := 0
	if allRequests := offerNode.GetAllRequests(); allRequests != nil {
		allRequestsCount = len(allRequests)
	}

	return model.NewMatchingResult(
		offerNode.Offer().UserID(),
		offerNode.Offer().ID(),
		newlyMatchedRequests,
		offerNode.Offer().Path(),
		allRequestsCount,
	), nil
}

func (matcher *Matcher) updateResults(offerNode *model.OfferNode) {
	matchingResult, err := createMatchingResult(offerNode)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create matching result for offer %s", offerNode.Offer().ID())
		return // continue
	}
	matcher.results = append(matcher.results, *matchingResult)
}
