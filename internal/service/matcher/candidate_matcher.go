package matcher

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"matching-engine/internal/collections"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
	"runtime"
)

var (
	workerCount = runtime.GOMAXPROCS(0) * 4
)

// buildCandidateMatches is responsible for matching offers and requests.
func (matcher *Matcher) buildCandidateMatches(offers []*model.Offer, requests []*model.Request) error {
	candidateIterator, err := matcher.candidateGenerator.GenerateCandidates(offers, requests)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	candidatesChannel := candidateIterator.Candidates(ctx, g)

	matcher.runCandidateProcessors(ctx, g, candidatesChannel)

	if err := g.Wait(); err != nil {
		return fmt.Errorf("matching failed: %w", err)
	}
	return nil
}

func (matcher *Matcher) runCandidateProcessors(
	ctx context.Context,
	g *errgroup.Group,
	candidatesChannel <-chan *model.MatchCandidate,
) {
	for i := 0; i < workerCount; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case matchCandidate, ok := <-candidatesChannel:
					if !ok {
						return nil
					}
					if err := matcher.processCandidate(matchCandidate); err != nil {
						return err
					}
				}
			}
		})
	}
}

func (matcher *Matcher) processCandidate(candidate *model.MatchCandidate) error {

	if candidate == nil {
		log.Error().Msg("Candidate is nil, skipping")
		return nil
	}

	if candidate.Offer() == nil || candidate.Request() == nil {
		log.Error().Msg("Candidate offer or request is nil, skipping")
		return nil
	}

	offerID := candidate.Offer().ID()
	requestID := candidate.Request().ID()

	if offerID == "" || requestID == "" {
		log.Error().Msg(errors.ErrEmptyOfferIDOrRequestID)
		return nil
	}

	// NOTE: The following block is not thread-safe if run concurrently for the same offerID.
	// Two goroutines could both see that the request set doesn't exist, create separate sets,
	// and overwrite each other when calling Set(), leading to lost request IDs.

	/*
		requestSet, exists := matcher.potentialOfferRequests.Get(offerID)
		if !exists {
			requestSet = collections.NewSet[string]()
		}
		requestSet.Add(requestID)
		matcher.potentialOfferRequests.Set(offerID, requestSet)
	*/

	// Thread-safe alternative using LoadOrStore to avoid lost updates.
	// This ensures only one set is created per offerID, even if accessed concurrently.
	requestSet, _ := matcher.potentialOfferRequests.GetOrStore(offerID, collections.NewSet[string]())
	requestSet.Add(requestID)

	// Note: It's okay if two goroutines race to write the same key here,
	// they will both write an identical OfferNode, so the second write
	// simply overwrites the first without any adverse effect.
	if _, exists := matcher.availableOffers.Get(offerID); !exists {
		matcher.availableOffers.Set(offerID, model.NewOfferNode(candidate.Offer()))
	}

	// Same goes here
	if _, exists := matcher.availableRequests.Get(requestID); !exists {
		matcher.availableRequests.Set(requestID, model.NewRequestNode(candidate.Request()))
	}
	return nil
}
