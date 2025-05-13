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
	availableOffers        map[string]*model.OfferNode
	availableRequests      map[string]*model.RequestNode
	potentialOfferRequests map[string]*collections.Set[string]
	pathPlanner            path_generator.PathPlanner
	candidateGenerator     earlypruning.CandidateGenerator
	maximumMatching        maximummatching.MaximumMatching
}

// NewMatcher creates a new Matcher instance
func NewMatcher(planner path_generator.PathPlanner, generator earlypruning.CandidateGenerator, matching maximummatching.MaximumMatching) *Matcher {
	return &Matcher{
		availableOffers:        make(map[string]*model.OfferNode),
		availableRequests:      make(map[string]*model.RequestNode),
		potentialOfferRequests: make(map[string]*collections.Set[string]),
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

		if _, exists := matcher.potentialOfferRequests[offerID]; !exists {
			matcher.potentialOfferRequests[offerID] = collections.NewSet[string]()
		}
		matcher.potentialOfferRequests[offerID].Add(requestID)

		if _, exists := matcher.availableOffers[offerID]; !exists {
			matcher.availableOffers[offerID] = model.NewOfferNode(offer)
		}

		if _, exists := matcher.availableRequests[requestID]; !exists {
			matcher.availableRequests[requestID] = model.NewRequestNode(request)
		}
	}

	// Create a slice to store the matching results
	results := make([]model.MatchingResult, 0, len(matcher.potentialOfferRequests))
	graph := model.NewGraph()

	for len(matcher.availableOffers) > 0 && len(matcher.availableRequests) > 0 {
		hasNewEdge := false
		potentialRequests := make(map[string]*model.RequestNode)

		// Find feasible paths and build the graph
		for offerID, requests := range matcher.potentialOfferRequests {
			offerNode, exists := matcher.availableOffers[offerID]
			if !exists {
				continue
			}

			for _, requestID := range requests.ToSlice() {
				requestNode, exists := matcher.availableRequests[requestID]
				if !exists {
					continue
				}

				newPath, isFeasible, err := matcher.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
				if err != nil {
					return nil, fmt.Errorf("failed to find feasible path: %w", err)
				}

				if isFeasible {
					hasNewEdge = true
					edge := model.NewEdge(requestNode, newPath)
					offerNode.AddEdge(edge)
					potentialRequests[requestID] = requestNode
					graph.AddOfferNode(offerNode)
				} else {
					requests.Remove(requestID)
				}
			}
		}

		// If no new edges were found, break the loop
		if !hasNewEdge {
			break
		}

		// Process offers that are not in the graph
		potentialOffers := graph.OfferNodes()
		for offerID, offerNode := range matcher.availableOffers {
			if potentialOffers.Contains(offerNode) {
				continue
			}

			if offerNode.IsMatched() {
				results = append(results, *createMatchingResult(offerNode, offerNode.NewlyAssignedMatchedRequests()))
			}

			delete(matcher.potentialOfferRequests, offerID)
			delete(matcher.availableOffers, offerID)
		}

		// Update available requests
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
			delete(matcher.availableRequests, requestID)

			// Remove the request from the potential offer requests
			matcher.potentialOfferRequests[offerID].Remove(requestID)

			if len(matchedRequests) < limit {
				offerNode.SetMatched(true)
				offerNode.SetNewlyAssignedMatchedRequests(matchedRequests)
			} else {
				delete(matcher.potentialOfferRequests, offerID)
				delete(matcher.availableOffers, offerID)
				results = append(results, *createMatchingResult(offerNode, matchedRequests))
			}
		}

		// Clear the graph and edges for the next iteration
		graph.ClearOfferNodes()
		for _, offerNode := range matcher.availableOffers {
			offerNode.ClearEdges()
		}
	}

	// Add any remaining matched offers to the results
	for _, offerNode := range matcher.availableOffers {
		if offerNode.IsMatched() {
			results = append(results, *createMatchingResult(offerNode, offerNode.NewlyAssignedMatchedRequests()))
		}
	}

	return results, nil
}
