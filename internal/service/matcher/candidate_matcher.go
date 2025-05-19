package matcher

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/collections"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
)

// buildCandidateMatches is responsible for matching offers and requests.
func (matcher *Matcher) buildCandidateMatches(offers []*model.Offer, requests []*model.Request) error {
	candidateIterator, err := matcher.candidateGenerator.GenerateCandidates(offers, requests)
	if err != nil {
		return err
	}

	for candidate, err := range candidateIterator.Candidates() {
		if err != nil {
			return fmt.Errorf("error during candidate iteration: %w", err)
		}
		if candidate == nil {
			log.Error().Msg("Candidate is nil, skipping")
			continue
		}

		if candidate.Offer() == nil || candidate.Request() == nil {
			log.Error().Msg("Candidate offer or request is nil, skipping")
			continue
		}

		offerID := candidate.Offer().ID()
		requestID := candidate.Request().ID()

		if offerID == "" || requestID == "" {
			log.Error().Msg(errors.ErrEmptyOfferIDOrRequestID)
			continue
		}

		requestSet, exists := matcher.potentialOfferRequests.Get(offerID)
		if !exists {
			requestSet = collections.NewSet[string]()
		}
		requestSet.Add(requestID)
		matcher.potentialOfferRequests.Set(offerID, requestSet)

		if _, exists := matcher.availableOffers.Get(offerID); !exists {
			matcher.availableOffers.Set(offerID, model.NewOfferNode(candidate.Offer()))
		}

		if _, exists := matcher.availableRequests.Get(requestID); !exists {
			matcher.availableRequests.Set(requestID, model.NewRequestNode(candidate.Request()))
		}
	}

	return nil
}
