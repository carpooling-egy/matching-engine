package service

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning"
	"matching-engine/internal/service/maximummatching"
	pathgenerator "matching-engine/internal/service/path-generator"

	"github.com/rs/zerolog/log"
)

type Matcher struct {
	availableOffers        *collections.SyncMap[string, *model.OfferNode]
	availableRequests      *collections.SyncMap[string, *model.RequestNode]
	potentialOfferRequests *collections.SyncMap[string, *collections.Set[string]]
	pathPlanner            pathgenerator.PathPlanner
	candidateGenerator     earlypruning.CandidateGenerator
	maximumMatching        maximummatching.MaximumMatching
}

// NewMatcher creates and initializes a new Matcher instance.
func NewMatcher(planner pathgenerator.PathPlanner, generator earlypruning.CandidateGenerator, matching maximummatching.MaximumMatching) *Matcher {
	return &Matcher{
		availableOffers:        collections.NewSyncMap[string, *model.OfferNode](),
		availableRequests:      collections.NewSyncMap[string, *model.RequestNode](),
		potentialOfferRequests: collections.NewSyncMap[string, *collections.Set[string]](),
		pathPlanner:            planner,
		candidateGenerator:     generator,
		maximumMatching:        matching,
	}
}

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

// Match performs the matching process for the input offers and requests with the given limit.
func (matcher *Matcher) Match(offers []*model.Offer, requests []*model.Request, limit int) ([]model.MatchingResult, error) {
	if len(offers) == 0 || len(requests) == 0 {
		return nil, fmt.Errorf(errors.ErrNoOffersOrRequests)
	}

	// Generate Candidates
	if err := matcher.generateCandidates(offers, requests); err != nil {
		return nil, fmt.Errorf("failed to generate candidates: %w", err)
	}

	results := make([]model.MatchingResult, 0)
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
		matcher.processUnmatchedOffers(graph, &results)

		// Find Maximum Matching
		if err = matcher.processMaximumMatching(graph, &results, limit); err != nil {
			return nil, fmt.Errorf("failed to process maximum matching: %w", err)
		}
	}

	// Handle remaining matched offers
	if err := matcher.processRemainingOffers(&results); err != nil {
		return nil, fmt.Errorf("failed to process remaining offers: %w", err)
	}

	return results, nil
}

// generateCandidates populates potential offer-request pairs.
func (matcher *Matcher) generateCandidates(offers []*model.Offer, requests []*model.Request) error {
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

		offerID := candidate.Offer().ID()
		requestID := candidate.Request().ID()

		if offerID == "" || requestID == "" {
			log.Error().Msg(errors.ErrEmptyOfferIDORRequestID)
			continue
		}

		if _, exists := matcher.potentialOfferRequests.Get(offerID); !exists {
			matcher.potentialOfferRequests.Set(offerID, collections.NewSet[string]())
		}
		requestSet, _ := matcher.potentialOfferRequests.Get(offerID)
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

// buildMatchingGraph constructs the graph by finding feasible paths and connecting offers with requests.
func (matcher *Matcher) buildMatchingGraph(graph *model.Graph, potentialRequests *collections.SyncMap[string, *model.RequestNode]) (bool, error) {
	hasNewEdge := false
	err := matcher.potentialOfferRequests.Range(func(offerID string, requestSet *collections.Set[string]) error {
		offerNode, exists := matcher.availableOffers.Get(offerID)
		if !exists || offerNode == nil {
			matcher.potentialOfferRequests.Delete(offerID)
			return nil
		}

		for _, requestID := range requestSet.ToSlice() {
			requestNode, exists := matcher.availableRequests.Get(requestID)
			if !exists || requestNode == nil {
				requestSet.Remove(requestID)
				continue
			}

			newPath, isFeasible, err := matcher.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
			if err != nil {
				return fmt.Errorf("failed to find feasible path for offer %s and request %s: %w", offerID, requestID, err)
			}

			if isFeasible && newPath != nil {
				hasNewEdge = true
				edge := model.NewEdge(requestNode, newPath)
				offerNode.AddEdge(edge)
				graph.AddOfferNode(offerNode)
				potentialRequests.Set(requestID, requestNode)
			} else {
				requestSet.Remove(requestID)
			}
		}
		return nil
	})

	return hasNewEdge, err
}

// processUnmatchedOffers processes offers that are not in the graph and updates results
func (matcher *Matcher) processUnmatchedOffers(graph *model.Graph, results *[]model.MatchingResult) {
	potentialOffers := graph.OfferNodes()
	matcher.availableOffers.ForEach(func(offerID string, offerNode *model.OfferNode) error {
		if potentialOffers.Contains(offerNode) {
			return nil // continue
		}

		if offerNode.IsMatched() {
			matchingResult, err := createMatchingResult(offerNode)
			if err != nil {
				log.Error().Err(err).Msgf("failed to create matching result for offer %s", offerID)
				return nil // continue
			}
			*results = append(*results, *matchingResult)
		}

		matcher.potentialOfferRequests.Delete(offerID)
		matcher.availableOffers.Delete(offerID)
		return nil // continue
	})
}

// processMaximumMatching finds maximum matches and updates results.
func (matcher *Matcher) processMaximumMatching(graph *model.Graph, results *[]model.MatchingResult, limit int) error {
	maxEdges, err := matcher.maximumMatching.FindMaximumMatching(graph)
	if err != nil {
		return fmt.Errorf("failed to find maximum matching: %w", err)
	}

	return maxEdges.Range(func(offerNode *model.OfferNode, edge *model.Edge) error {
		requestNode := edge.RequestNode()
		if requestNode == nil {
			log.Warn().Msg("Nil request node in maximum matching")
			return nil
		}

		offerNode.SetMatched(true)
		offerNode.AddMatchedRequest(requestNode.Request())

		if len(offerNode.GetAllRequests()) >= limit {
			matchingResult, err := createMatchingResult(offerNode)
			if err != nil {
				log.Error().Err(err).Msg("Failed to create matching result")
				return nil
			}
			*results = append(*results, *matchingResult)
			matcher.availableOffers.Delete(offerNode.Offer().ID())
			matcher.potentialOfferRequests.Delete(offerNode.Offer().ID())
		}

		matcher.availableRequests.Delete(requestNode.Request().ID())
		return nil
	})
}

// processRemainingOffers appends leftover matched offers to results.
func (matcher *Matcher) processRemainingOffers(results *[]model.MatchingResult) error {
	return matcher.availableOffers.Range(func(offerID string, offerNode *model.OfferNode) error {
		if offerNode.IsMatched() {
			matchingResult, err := createMatchingResult(offerNode)
			if err != nil {
				log.Warn().Err(err).Msgf("Failed to create matching result for offer: %s", offerID)
				return nil
			}
			*results = append(*results, *matchingResult)
		}
		return nil
	})
}
