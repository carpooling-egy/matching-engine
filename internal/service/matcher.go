package service

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning"
	"matching-engine/internal/service/maximummatching"
	path_generator "matching-engine/internal/service/path-generator"

	"github.com/rs/zerolog/log"
)

type Matcher struct {
	availableOffers        *collections.SyncMap[string, *model.OfferNode]
	availableRequests      *collections.SyncMap[string, *model.RequestNode]
	potentialOfferRequests *collections.SyncMap[string, *collections.Set[string]]
	pathPlanner            path_generator.PathPlanner
	candidateGenerator     earlypruning.CandidateGenerator
	maximumMatching        maximummatching.MaximumMatching
}

// NewMatcher creates a new Matcher instance
func NewMatcher(planner path_generator.PathPlanner, generator earlypruning.CandidateGenerator, matching maximummatching.MaximumMatching) *Matcher {
	return &Matcher{
		availableOffers:        collections.NewSyncMap[string, *model.OfferNode](),
		availableRequests:      collections.NewSyncMap[string, *model.RequestNode](),
		potentialOfferRequests: collections.NewSyncMap[string, *collections.Set[string]](),
		pathPlanner:            planner,
		candidateGenerator:     generator,
		maximumMatching:        matching,
	}
}

// createMatchingResult is a helper function to create a matching result
func createMatchingResult(offerNode *model.OfferNode) (*model.MatchingResult, error) {
	if offerNode == nil {
		return nil, fmt.Errorf("offer node is nil")
	}

	offer := offerNode.Offer()
	if offer == nil {
		return nil, fmt.Errorf("offer node has nil offer")
	}

	offerID := offer.ID()
	if offerID == "" {
		return nil, fmt.Errorf("offer ID is empty")
	}
	path := offer.Path()
	if path == nil {
		return nil, fmt.Errorf("offer node has nil path")
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("offer node has empty path")
	}
	newlyMatchedRequests := offerNode.NewlyAssignedMatchedRequests()
	if newlyMatchedRequests == nil {
		return nil, fmt.Errorf("offer node has nil newly assigned matched requests")
	}
	if len(newlyMatchedRequests) == 0 {
		return nil, fmt.Errorf("offer node has empty newly assigned matched requests")
	}
	allRequests := offerNode.GetAllRequests()
	allRequestsCount := 0
	if allRequests != nil {
		allRequestsCount = len(allRequests)
	}
	return model.NewMatchingResult(offer.UserID(), offerID, newlyMatchedRequests, path, allRequestsCount), nil
}

func (matcher *Matcher) Match(offers []*model.Offer, requests []*model.Request, limit int) ([]model.MatchingResult, error) {
	if len(offers) == 0 || len(requests) == 0 {
		return nil, fmt.Errorf("no offers or requests")
	}

	// Initialize the matcher with offers and requests
	candidateIterator, err := matcher.candidateGenerator.GenerateCandidates(offers, requests)
	if err != nil {
		return nil, fmt.Errorf("failed to generate candidates: %w", err)
	}

	// Iterate through the candidates and find matches
	for candidate, err := range candidateIterator.Candidates() {
		if err != nil {
			return nil, fmt.Errorf("failed to get candidates: %w", err)
		}

		if candidate == nil {
			log.Error().Msg("candidate is nil")
			continue
		}

		offer := candidate.Offer()
		request := candidate.Request()
		if offer == nil {
			log.Error().Msg("offer is nil")
			continue
		}
		if request == nil {
			log.Error().Msg("request is nil")
			continue
		}

		offerID := offer.ID()
		requestID := request.ID()

		// Get or create a set of request IDs for this offer
		requestSet, exists := matcher.potentialOfferRequests.Get(offerID)
		if !exists {
			requestSet = collections.NewSet[string]()
			matcher.potentialOfferRequests.Set(offerID, requestSet)
		}
		requestSet.Add(requestID)

		// Store the offer node if it doesn't exist
		_, exists = matcher.availableOffers.Get(offerID)
		if !exists {
			matcher.availableOffers.Set(offerID, model.NewOfferNode(offer))
		}

		// Store the request node if it doesn't exist
		_, exists = matcher.availableRequests.Get(requestID)
		if !exists {
			matcher.availableRequests.Set(requestID, model.NewRequestNode(request))
		}
	}

	// Create a slice to store the matching results
	results := make([]model.MatchingResult, 0, matcher.potentialOfferRequests.Size())
	graph := model.NewGraph()

	for matcher.availableOffers.Size() > 0 && matcher.availableRequests.Size() > 0 {
		hasNewEdge := false
		potentialRequests := collections.NewSyncMap[string, *model.RequestNode]()

		// Find feasible paths and build the graph
		err := matcher.potentialOfferRequests.Range(func(offerID string, requests *collections.Set[string]) error {
			if offerID == "" {
				return fmt.Errorf("empty offer ID encountered")
			}

			if requests == nil {
				return fmt.Errorf("nil requests set for offer ID: %s", offerID)
			}

			offerNode, exists := matcher.availableOffers.Get(offerID)
			if !exists {
				return nil // continue
			}

			if offerNode == nil {
				return fmt.Errorf("nil offer node found for ID: %s", offerID)
			}

			for _, requestID := range requests.ToSlice() {
				if requestID == "" {
					continue // Skip empty request IDs
				}

				requestNode, exists := matcher.availableRequests.Get(requestID)
				if !exists {
					continue
				}

				if requestNode == nil {
					continue // Skip nil request nodes
				}

				newPath, isFeasible, err := matcher.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
				if err != nil {
					return fmt.Errorf("failed to find feasible path for offer %s and request %s: %w", offerID, requestID, err)
				}

				if isFeasible {
					if newPath == nil {
						return fmt.Errorf("feasible path is nil for offer %s and request %s", offerID, requestID)
					}
					hasNewEdge = true
					edge := model.NewEdge(requestNode, newPath)
					offerNode.AddEdge(edge)
					potentialRequests.Set(requestID, requestNode)
					graph.AddOfferNode(offerNode)
				} else {
					log.Info().Msgf("no feasible path for offer %s and request %s", offerID, requestID)
					requests.Remove(requestID)
				}
			}
			return nil // continue
		})
		if err != nil {
			return nil, fmt.Errorf("failed to process potential offer requests: %w", err)
		}

		// If no new edges were found, break the loop
		if !hasNewEdge {
			break
		}

		// Process offers that are not in the graph
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
				results = append(results, *matchingResult)
			}

			matcher.potentialOfferRequests.Delete(offerID)
			matcher.availableOffers.Delete(offerID)
			return nil // continue
		})

		// Update available requests by replacing the current cache with the new one
		matcher.availableRequests = potentialRequests

		// Find maximum matching
		maximumMatchingEdges, err := matcher.maximumMatching.FindMaximumMatching(graph)
		if err != nil {
			return nil, fmt.Errorf("failed to find maximum matching: %w", err)
		}

		// Process the matching results
		err = maximumMatchingEdges.Range(func(offerNode *model.OfferNode, edge *model.Edge) error {
			if edge == nil {
				return fmt.Errorf("nil edge encountered in maximum matching")
			}

			if edge.RequestNode() == nil {
				return fmt.Errorf("edge with nil request node encountered")
			}

			if offerNode == nil {
				return fmt.Errorf("nil offer node encountered in maximum matching")
			}

			if offerNode.Offer() == nil {
				return fmt.Errorf("offer node with nil offer encountered")
			}

			requestNode := edge.RequestNode()
			request := requestNode.Request()
			if request == nil {
				return fmt.Errorf("request node with nil request encountered")
			}

			offerID := offerNode.Offer().ID()
			if offerID == "" {
				return fmt.Errorf("offer with empty ID encountered")
			}

			requestID := request.ID()
			if requestID == "" {
				return fmt.Errorf("request with empty ID encountered")
			}

			// Add the new request to the matched requests
			newlyMatchedRequests := offerNode.NewlyAssignedMatchedRequests()
			if newlyMatchedRequests == nil {
				newlyMatchedRequests = make([]*model.Request, 0)
			}
			newlyMatchedRequests = append(newlyMatchedRequests, request)

			// Set the path for the offer node
			newPath := edge.NewPath()
			if newPath == nil {
				return fmt.Errorf("edge with nil path encountered for offer %s and request %s", offerID, requestID)
			}
			offerNode.Offer().SetPath(newPath)

			// Delete the request node from the available requests
			matcher.availableRequests.Delete(requestID)

			// Remove the request from the potential offer requests
			requestSet, exists := matcher.potentialOfferRequests.Get(offerID)
			if exists {
				requestSet.Remove(requestID)
			}
			offerNode.SetMatched(true)
			offerNode.SetNewlyAssignedMatchedRequests(newlyMatchedRequests)

			if len(offerNode.GetAllRequests()) >= limit {
				matcher.potentialOfferRequests.Delete(offerID)
				matcher.availableOffers.Delete(offerID)
				matchingResult, err := createMatchingResult(offerNode)
				if err != nil {
					return fmt.Errorf("failed to create matching result for offer %s: %w", offerID, err)
				}
				results = append(results, *matchingResult)
			}
			return nil // continue
		})
		if err != nil {
			return nil, fmt.Errorf("failed to process maximum matching edges: %w", err)
		}

		// Clear the graph and edges for the next iteration
		graph.ClearOfferNodes()
		matcher.availableOffers.ForEach(func(offerID string, offerNode *model.OfferNode) error {
			if offerNode == nil {
				log.Warn().Msgf("nil offer node found with ID: %s", offerID)
				matcher.availableOffers.Delete(offerID)
				matcher.potentialOfferRequests.Delete(offerID)
				return nil
			}
			offerNode.ClearEdges()
			return nil
		})
	}

	// Add any remaining matched offers to the results
	matcher.availableOffers.ForEach(func(offerID string, offerNode *model.OfferNode) error {
		if offerNode == nil {
			log.Warn().Msgf("nil offer node found with ID: %s", offerID)
			return nil
		}

		if offerNode.IsMatched() {
			matchedRequests := offerNode.NewlyAssignedMatchedRequests()
			if matchedRequests == nil {
				return fmt.Errorf("matched offer node with ID: %s has nil matched requests", offerID)
			}
			matchingResult, err := createMatchingResult(offerNode)
			if err != nil {
				return fmt.Errorf("failed to create matching result for offer %s: %w", offerID, err)
			}
			results = append(results, *matchingResult)
		}
		return nil
	})

	return results, nil
}
