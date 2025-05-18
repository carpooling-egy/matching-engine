package matcher

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning"
	"matching-engine/internal/service/maximummatching"
	"matching-engine/internal/service/pathgeneration/planner"

	"github.com/rs/zerolog/log"
)

type Matcher struct {
	availableOffers        *collections.SyncMap[string, *model.OfferNode]
	availableRequests      *collections.SyncMap[string, *model.RequestNode]
	potentialOfferRequests *collections.SyncMap[string, *collections.Set[string]]
	results                []model.MatchingResult
	pathPlanner            planner.PathPlanner
	candidateGenerator     earlypruning.CandidateGenerator
	maximumMatching        maximummatching.MaximumMatching
}

// NewMatcher creates and initializes a new Matcher instance.
func NewMatcher(planner planner.PathPlanner, generator earlypruning.CandidateGenerator, matching maximummatching.MaximumMatching) *Matcher {
	return &Matcher{
		availableOffers:        collections.NewSyncMap[string, *model.OfferNode](),
		availableRequests:      collections.NewSyncMap[string, *model.RequestNode](),
		potentialOfferRequests: collections.NewSyncMap[string, *collections.Set[string]](),
		results:                make([]model.MatchingResult, 0),
		pathPlanner:            planner,
		candidateGenerator:     generator,
		maximumMatching:        matching,
	}
}

// Match performs the matching process for the input offers and requests with the given limit.
func (matcher *Matcher) Match(offers []*model.Offer, requests []*model.Request, limit int) ([]model.MatchingResult, error) {
	if offers == nil || requests == nil || len(offers) == 0 || len(requests) == 0 {
		return nil, fmt.Errorf(errors.ErrNoOffersOrRequests)
	}

	// Generate Candidates
	if err := matcher.buildCandidateMatches(offers, requests); err != nil {
		return nil, fmt.Errorf("failed to build candidate matches: %w", err)
	}

	graph := model.NewGraph()

	for matcher.availableOffers.Size() > 0 && matcher.availableRequests.Size() > 0 {
		// Build Matching Graph
		potentialRequests := collections.NewSyncMap[string, *model.RequestNode]()
		hasNewEdge, err := matcher.buildMatchingGraph(graph, potentialRequests)
		if err != nil {
			return nil, fmt.Errorf("failed to build matching graph: %w", err)
		}

		if !hasNewEdge {
			log.Info().Msg("No new edges found, stopping matching process")
			break
		}
		// Update the graph with potential requests
		matcher.availableRequests = potentialRequests

		// Process unmatched offers
		matcher.processUnmatchedOffers(graph)

		// Find Maximum Matching
		if err = matcher.processMaximumMatching(graph, limit); err != nil {
			return nil, fmt.Errorf("failed to process maximum matching: %w", err)
		}
	}

	// Handle remaining matched offers
	if err := matcher.processRemainingOffers(); err != nil {
		return nil, fmt.Errorf("failed to process remaining offers: %w", err)
	}

	return matcher.results, nil
}
