package service

import (
	"fmt"
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning"
	"matching-engine/internal/service/maximummatching"
	"matching-engine/internal/service/path-generator"
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
func createMatchingResult(offerNode *model.OfferNode, requests []*model.Request) *model.MatchingResult {
	offerID := offerNode.Offer().ID()
	return model.NewMatchingResult(offerID, offerID, requests, offerNode.Offer().Path(), len(offerNode.GetAllRequests()))
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
			continue
		}

		offer := candidate.Offer()
		request := candidate.Request()
		if offer == nil || request == nil {
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
		matcher.potentialOfferRequests.Range(func(offerID string, requests *collections.Set[string]) bool {
			offerNode, exists := matcher.availableOffers.Get(offerID)
			if !exists {
				return true // continue
			}

			for _, requestID := range requests.ToSlice() {
				requestNode, exists := matcher.availableRequests.Get(requestID)
				if !exists {
					continue
				}

				newPath, isFeasible, err := matcher.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
				if err != nil {
					// We can't return an error directly from the Range function, so we'll set it and break the loop
					err = fmt.Errorf("failed to find feasible path: %w", err)
					return false // break
				}

				if isFeasible {
					hasNewEdge = true
					edge := model.NewEdge(requestNode, newPath)
					offerNode.AddEdge(edge)
					potentialRequests.Set(requestID, requestNode)
					graph.AddOfferNode(offerNode)
				} else {
					requests.Remove(requestID)
				}
			}
			return true // continue
		})

		// If no new edges were found, break the loop
		if !hasNewEdge {
			break
		}

		// Process offers that are not in the graph
		potentialOffers := graph.OfferNodes()
		matcher.availableOffers.Range(func(offerID string, offerNode *model.OfferNode) bool {
			if potentialOffers.Contains(offerNode) {
				return true // continue
			}

			if offerNode.IsMatched() {
				results = append(results, *createMatchingResult(offerNode, offerNode.NewlyAssignedMatchedRequests()))
			}

			matcher.potentialOfferRequests.Delete(offerID)
			matcher.availableOffers.Delete(offerID)
			return true // continue
		})

		// Update available requests by replacing the current cache with the new one
		matcher.availableRequests = potentialRequests

		// Find maximum matching
		maximumMatchingEdges, err := matcher.maximumMatching.FindMaximumMatching(graph)
		if err != nil {
			return nil, fmt.Errorf("failed to find maximum matching: %w", err)
		}

		// Process the matching results
		for offerNode, edge := range maximumMatchingEdges {
			if edge == nil || edge.RequestNode() == nil || offerNode == nil {
				continue
			}

			requestNode := edge.RequestNode()
			request := requestNode.Request()
			offerID := offerNode.Offer().ID()
			requestID := request.ID()

			// Add the new request to the matched requests
			matchedRequests := append(offerNode.NewlyAssignedMatchedRequests(), request)

			// Set the path for the offer node
			offerNode.Offer().SetPath(edge.NewPath())

			// Delete the request node from the available requests
			matcher.availableRequests.Delete(requestID)

			// Remove the request from the potential offer requests
			requestSet, exists := matcher.potentialOfferRequests.Get(offerID)
			if exists {
				requestSet.Remove(requestID)
			}

			if len(matchedRequests) < limit {
				offerNode.SetMatched(true)
				offerNode.SetNewlyAssignedMatchedRequests(matchedRequests)
			} else {
				matcher.potentialOfferRequests.Delete(offerID)
				matcher.availableOffers.Delete(offerID)
				results = append(results, *createMatchingResult(offerNode, matchedRequests))
			}
		}

		// Clear the graph and edges for the next iteration
		graph.ClearOfferNodes()
		matcher.availableOffers.ForEach(func(offerID string, offerNode *model.OfferNode) {
			offerNode.ClearEdges()
		})
	}

	// Add any remaining matched offers to the results
	matcher.availableOffers.ForEach(func(offerID string, offerNode *model.OfferNode) {
		if offerNode.IsMatched() {
			results = append(results, *createMatchingResult(offerNode, offerNode.NewlyAssignedMatchedRequests()))
		}
	})

	return results, nil
}
