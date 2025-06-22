package matcher

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/collections"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning"
	"matching-engine/internal/service/matchevaluator"
	"matching-engine/internal/service/maximummatching"
	"matching-engine/internal/service/timematrix"
)

const (
	// DefaultLimit is the default limit for the number of requests per offer.
	DefaultLimit = 5
)

type Matcher struct {
	availableOffers          *collections.SyncMap[string, *model.OfferNode]
	availableRequests        *collections.SyncMap[string, *model.RequestNode]
	potentialOfferRequests   *collections.SyncMap[string, *collections.Set[string]]
	results                  []*model.MatchingResult
	matchEvaluator           matchevaluator.Evaluator
	candidateGenerator       earlypruning.CandidateGenerator
	maximumMatching          maximummatching.MaximumMatching
	timeMatrixCachePopulator timematrix.Populator
	limit                    int
}

// NewMatcher creates and initializes a new Matcher instance.
func NewMatcher(evaluator matchevaluator.Evaluator, generator earlypruning.CandidateGenerator, matching maximummatching.MaximumMatching, cachePopulator timematrix.Populator) *Matcher {
	if evaluator == nil {
		log.Error().Msg("Matcher: Evaluator is nil")
		panic("Matcher: Evaluator is nil")
		return nil
	}
	return &Matcher{
		availableOffers:          collections.NewSyncMap[string, *model.OfferNode](),
		availableRequests:        collections.NewSyncMap[string, *model.RequestNode](),
		potentialOfferRequests:   collections.NewSyncMap[string, *collections.Set[string]](),
		results:                  make([]*model.MatchingResult, 0),
		matchEvaluator:           evaluator,
		candidateGenerator:       generator,
		maximumMatching:          matching,
		limit:                    DefaultLimit,
		timeMatrixCachePopulator: cachePopulator,
	}
}

// Match performs the matching process for the input offers and requests with the given limit.
func (matcher *Matcher) Match(offers []*model.Offer, requests []*model.Request) ([]*model.MatchingResult, error) {
	if offers == nil || requests == nil || len(offers) == 0 || len(requests) == 0 {
		return nil, fmt.Errorf(errors.ErrNoOffersOrRequests)
	}

	// Generate Candidates
	if err := matcher.buildCandidateMatches(offers, requests); err != nil {
		return nil, fmt.Errorf("failed to build candidate matches: %w", err)
	}

	//fmt.Println(matcher.potentialOfferRequests, "potential offer requests found")
	//v, ok := matcher.potentialOfferRequests.Get("1")
	//if !ok {
	//	return nil, nil
	//}
	//fmt.Println("Potential requests for offer 1:", v)
	//for _, requestID := range v.ToSlice() {
	//	fmt.Println("Potential requests for offer 1:", requestID)
	//}
	graph := model.NewGraph()

	for matcher.availableOffers.Size() > 0 && matcher.availableRequests.Size() > 0 {
		// Build Matching Graph
		log.Info().Msg("Building matching graph")
		fmt.Printf("matcher.availableOffers.Size(): %d, matcher.availableRequests.Size(): %d\n", matcher.availableOffers.Size(), matcher.availableRequests.Size())
		matcher.availableRequests.Range(func(requestId string, requestNode *model.RequestNode) error {
			fmt.Println("Processing request:", requestId, requestNode.Request().ID())
			return nil
		})
		matcher.availableOffers.Range(func(offerID string, offerNode *model.OfferNode) error {
			fmt.Println("Processing offer:", offerID)
			return nil
		})
		// Build the matching graph with potential edges between offers and requests
		fmt.Println(matcher.availableOffers.Size())
		fmt.Println(matcher.availableRequests.Size())
		hasNewEdge, err := matcher.buildMatchingGraph(graph)
		if err != nil {
			return nil, fmt.Errorf("failed to build matching graph: %w", err)
		}

		if !hasNewEdge {
			log.Info().Msg("No new edges found, stopping matching process")
			break
		}

		// Process unmatched offers
		matcher.processUnmatchedOffers(graph)

		// Update the graph with potential offers
		matcher.availableOffers = graph.OfferNodes()

		// Update the graph with potential requests
		matcher.availableRequests = graph.RequestNodes()

		fmt.Println(matcher.availableRequests)
		fmt.Println("Available requests:", matcher.availableRequests.Size())
		fmt.Println("Available offers:", matcher.availableOffers.Size())

		// Find Maximum Matching
		if err = matcher.processMaximumMatching(graph, matcher.limit); err != nil {
			return nil, fmt.Errorf("failed to process maximum matching: %w", err)
		}
		// Clear the graph and edges for the next iteration
		graph.Clear()
		fmt.Println("Graph cleared for next iteration")
		matcher.availableRequests.Range(func(requestId string, requestNode *model.RequestNode) error {
			fmt.Println("Processing request:", requestId, requestNode.Request().ID())
			return nil
		})
		matcher.availableOffers.Range(func(offerID string, offerNode *model.OfferNode) error {
			fmt.Println("Processing offer:", offerID)
			return nil
		})

	}

	// Handle remaining matched offers
	if err := matcher.processRemainingOffers(); err != nil {
		return nil, fmt.Errorf("failed to process remaining offers: %w", err)
	}

	return matcher.results, nil
}
