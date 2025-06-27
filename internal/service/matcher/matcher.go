package matcher

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/cmd/appmetrics"
	"matching-engine/internal/collections"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning"
	"matching-engine/internal/service/matchevaluator"
	"matching-engine/internal/service/maximummatching"
	"matching-engine/internal/service/timematrix"
	"time"
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
	timeMatrixCachePopulator *timematrix.CacheWithOfferIdPopulator
	limit                    int
}

// NewMatcher creates and initializes a new Matcher instance.
func NewMatcher(evaluator matchevaluator.Evaluator, generator earlypruning.CandidateGenerator, matching maximummatching.MaximumMatching, cachePopulator *timematrix.CacheWithOfferIdPopulator) *Matcher {
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
	startTime := time.Now()
	if err := matcher.buildCandidateMatches(offers, requests); err != nil {
		return nil, fmt.Errorf("failed to build candidate matches: %w", err)
	}
	candidateGenerationTime := time.Since(startTime)
	buildingGraphTime, maximummatchingTime := time.Duration(0), time.Duration(0)
	matcherStartTime := time.Now()

	graph := model.NewMaximumMatchingGraph()

	for matcher.availableOffers.Size() > 0 && matcher.availableRequests.Size() > 0 {
		// Build Matching Graph
		log.Debug().Msg("Building matching graph")
		// Build the matching graph with potential edges between offers and requests
		startTime = time.Now()
		hasNewEdge, err := matcher.buildMatchingGraph(graph)
		log.Debug().Msgf("Graph building took %s", time.Since(startTime))
		buildingGraphTime += time.Since(startTime)
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

		// Find Maximum Matching
		log.Debug().Msg("Processing maximum matching")
		startTime = time.Now()
		if err = matcher.processMaximumMatching(graph, matcher.limit); err != nil {
			return nil, fmt.Errorf("failed to process maximum matching: %w", err)
		}
		log.Debug().Msgf("Maximum matching processing took %s", time.Since(startTime))
		maximummatchingTime += time.Since(startTime)
		// Clear the graph and edges for the next iteration
		graph.Clear()

	}
	// Log the time taken for candidate generation
	log.Info().
		Str("duration", candidateGenerationTime.String()).
		Msg("Candidate generation completed")
	appmetrics.TrackTime("Average candidate generation duration", candidateGenerationTime)
	// Log the total time taken for the matcher loop
	log.Info().
		Str("duration", time.Since(matcherStartTime).String()).
		Msg("Matcher loop completed")
	appmetrics.TrackTime("Average matcher loop duration", time.Since(matcherStartTime))
	// Log the time taken for building the graph
	log.Info().
		Str("duration", buildingGraphTime.String()).
		Msg("Building matching graph completed")
	appmetrics.TrackTime("Average building graph duration", buildingGraphTime)
	// Log the time taken for maximum matching
	log.Info().
		Str("duration", maximummatchingTime.String()).
		Msg("Maximum matching processing completed")
	appmetrics.TrackTime("Average maximum matching duration", maximummatchingTime)
	// Handle remaining matched offers
	startTime = time.Now()
	if err := matcher.processRemainingOffers(); err != nil {
		return nil, fmt.Errorf("failed to process remaining offers: %w", err)
	}
	log.Info().
		Str("duration", time.Since(startTime).String()).
		Msg("Processing remaining offers completed")

	return matcher.results, nil
}
